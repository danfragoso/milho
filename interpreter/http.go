package interpreter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/danfragoso/milho/mir"
)

// HTTP server state management
var (
	httpServers     = make(map[int64]*httpServerInstance)
	httpServerMutex sync.RWMutex
	nextServerID    int64 = 1
)

type httpServerInstance struct {
	ID            int64
	Port          int64
	Server        *http.Server
	Mux           *http.ServeMux
	Routes        map[string]*routeHandler // method:path -> handler
	PathHandlers  map[string]map[string]mir.Expression // path -> method -> handler
	RegisteredPaths map[string]bool // tracks which paths are registered with mux
	Session       *mir.Session
	stopChan      chan struct{}
}

type routeHandler struct {
	Method  string
	Path    string
	Handler mir.Expression
}

// __httpServerCreate creates a new HTTP server instance
func __httpServerCreate(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-server-create: expected 1 parameter (port), got %d", len(params))
	}

	portExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if portExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("http-server-create: port must be a number, got %s", portExpr.Type())
	}

	port := portExpr.(*mir.NumberExpression).Numerator

	httpServerMutex.Lock()
	serverID := nextServerID
	nextServerID++

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	instance := &httpServerInstance{
		ID:              serverID,
		Port:            port,
		Server:          server,
		Mux:             mux,
		Routes:          make(map[string]*routeHandler),
		PathHandlers:    make(map[string]map[string]mir.Expression),
		RegisteredPaths: make(map[string]bool),
		Session:         session,
		stopChan:        make(chan struct{}),
	}

	httpServers[serverID] = instance
	httpServerMutex.Unlock()

	return mir.CreateNumberExpression(serverID, 1)
}

// __httpServerRoute registers a route handler
func __httpServerRoute(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 4 {
		return nil, fmt.Errorf("http-server-route: expected 4 parameters (serverID, method, path, handler), got %d", len(params))
	}

	serverIDExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if serverIDExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("http-server-route: serverID must be a number, got %s", serverIDExpr.Type())
	}

	serverID := serverIDExpr.(*mir.NumberExpression).Numerator

	methodExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if methodExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("http-server-route: method must be a string, got %s", methodExpr.Type())
	}

	method := methodExpr.(*mir.StringExpression).Val

	pathExpr, err := evaluate(params[2], session)
	if err != nil {
		return nil, err
	}

	if pathExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("http-server-route: path must be a string, got %s", pathExpr.Type())
	}

	path := pathExpr.(*mir.StringExpression).Val

	handlerExpr, err := evaluate(params[3], session)
	if err != nil {
		return nil, err
	}

	if handlerExpr.Type() != mir.FunctionExpr {
		return nil, fmt.Errorf("http-server-route: handler must be a function, got %s", handlerExpr.Type())
	}

	httpServerMutex.Lock()
	instance, exists := httpServers[serverID]
	if !exists {
		httpServerMutex.Unlock()
		return nil, fmt.Errorf("http-server-route: server %d not found", serverID)
	}

	// Store handler in path -> method map
	if instance.PathHandlers[path] == nil {
		instance.PathHandlers[path] = make(map[string]mir.Expression)
	}
	instance.PathHandlers[path][method] = handlerExpr

	// Also store in legacy Routes map for compatibility
	routeKey := method + ":" + path
	route := &routeHandler{
		Method:  method,
		Path:    path,
		Handler: handlerExpr,
	}
	instance.Routes[routeKey] = route

	// Register the path with mux only if not already registered
	if !instance.RegisteredPaths[path] {
		instance.RegisteredPaths[path] = true

		// Register a dispatcher that handles all methods for this path
		instance.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			// Look up the handler for this method
			httpServerMutex.RLock()
			methodHandlers, pathExists := instance.PathHandlers[path]
			if !pathExists {
				httpServerMutex.RUnlock()
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			handlerExpr, methodExists := methodHandlers[r.Method]
			if !methodExists {
				httpServerMutex.RUnlock()
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			httpServerMutex.RUnlock()

			// Create request map
			reqMap, err := createRequestMap(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating request: %s", err), http.StatusInternalServerError)
				return
			}

			// Create response map
			resMap, err := createResponseMap()
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating response: %s", err), http.StatusInternalServerError)
				return
			}

			// Clone the session for this request to avoid state contamination
			requestSession := instance.Session.Clone()

			// Call the handler function
			handlerParams := []mir.Expression{reqMap, resMap}
			resultExpr, err := callFunction(handlerExpr.(*mir.FunctionExpression), handlerParams, requestSession)
			if err != nil {
				http.Error(w, fmt.Sprintf("Handler error: %s", err), http.StatusInternalServerError)
				return
			}

			// Extract response from result (should be the modified response map)
			if resultExpr.Type() == mir.MapExpr {
				resMap = resultExpr.(*mir.MapExpression)
			}

			// Write response
			writeResponse(w, resMap)
		})
	}

	httpServerMutex.Unlock()

	return mir.CreateNilExpression()
}

// __httpServerStart starts the HTTP server (blocking)
func __httpServerStart(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-server-start: expected 1 parameter (serverID), got %d", len(params))
	}

	serverIDExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if serverIDExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("http-server-start: serverID must be a number, got %s", serverIDExpr.Type())
	}

	serverID := serverIDExpr.(*mir.NumberExpression).Numerator

	httpServerMutex.RLock()
	instance, exists := httpServers[serverID]
	httpServerMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("http-server-start: server %d not found", serverID)
	}

	// Start server (blocking)
	if err := instance.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return nil, fmt.Errorf("http-server-start: %s", err)
	}

	return mir.CreateNilExpression()
}

// __httpServerStartAsync starts the HTTP server in the background
func __httpServerStartAsync(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-server-start-async: expected 1 parameter (serverID), got %d", len(params))
	}

	serverIDExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if serverIDExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("http-server-start-async: serverID must be a number, got %s", serverIDExpr.Type())
	}

	serverID := serverIDExpr.(*mir.NumberExpression).Numerator

	httpServerMutex.RLock()
	instance, exists := httpServers[serverID]
	httpServerMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("http-server-start-async: server %d not found", serverID)
	}

	// Start server in background
	go func() {
		if err := instance.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %s\n", err)
		}
	}()

	return mir.CreateNilExpression()
}

// __httpServerStop stops a running HTTP server
func __httpServerStop(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-server-stop: expected 1 parameter (serverID), got %d", len(params))
	}

	serverIDExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if serverIDExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("http-server-stop: serverID must be a number, got %s", serverIDExpr.Type())
	}

	serverID := serverIDExpr.(*mir.NumberExpression).Numerator

	httpServerMutex.Lock()
	instance, exists := httpServers[serverID]
	if !exists {
		httpServerMutex.Unlock()
		return nil, fmt.Errorf("http-server-stop: server %d not found", serverID)
	}

	// Shutdown server gracefully
	ctx := context.Background()
	if err := instance.Server.Shutdown(ctx); err != nil {
		httpServerMutex.Unlock()
		return nil, fmt.Errorf("http-server-stop: %s", err)
	}

	close(instance.stopChan)
	delete(httpServers, serverID)
	httpServerMutex.Unlock()

	return mir.CreateNilExpression()
}

// Request helper functions

// __httpRequestMethod gets the HTTP method from request
func __httpRequestMethod(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-request-method: expected 1 parameter (request), got %d", len(params))
	}

	reqExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if reqExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-request-method: request must be a map, got %s", reqExpr.Type())
	}

	reqMap := reqExpr.(*mir.MapExpression)
	if method, ok := reqMap.Values["method"]; ok {
		return method, nil
	}

	return mir.CreateStringExpression("")
}

// __httpRequestPath gets the path from request
func __httpRequestPath(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-request-path: expected 1 parameter (request), got %d", len(params))
	}

	reqExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if reqExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-request-path: request must be a map, got %s", reqExpr.Type())
	}

	reqMap := reqExpr.(*mir.MapExpression)
	if path, ok := reqMap.Values["path"]; ok {
		return path, nil
	}

	return mir.CreateStringExpression("")
}

// __httpRequestHeaders gets headers from request
func __httpRequestHeaders(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-request-headers: expected 1 parameter (request), got %d", len(params))
	}

	reqExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if reqExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-request-headers: request must be a map, got %s", reqExpr.Type())
	}

	reqMap := reqExpr.(*mir.MapExpression)
	if headers, ok := reqMap.Values["headers"]; ok {
		return headers, nil
	}

	return mir.CreateMapExpression(make(map[string]mir.Expression))
}

// __httpRequestBody gets body from request
func __httpRequestBody(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-request-body: expected 1 parameter (request), got %d", len(params))
	}

	reqExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if reqExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-request-body: request must be a map, got %s", reqExpr.Type())
	}

	reqMap := reqExpr.(*mir.MapExpression)
	if body, ok := reqMap.Values["body"]; ok {
		return body, nil
	}

	return mir.CreateStringExpression("")
}

// __httpRequestQuery gets query parameters from request
func __httpRequestQuery(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("http-request-query: expected 1 parameter (request), got %d", len(params))
	}

	reqExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if reqExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-request-query: request must be a map, got %s", reqExpr.Type())
	}

	reqMap := reqExpr.(*mir.MapExpression)
	if query, ok := reqMap.Values["query"]; ok {
		return query, nil
	}

	return mir.CreateMapExpression(make(map[string]mir.Expression))
}

// Response helper functions

// __httpResponseStatus sets response status code
func __httpResponseStatus(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("http-response-status: expected 2 parameters (response, status), got %d", len(params))
	}

	resExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if resExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-response-status: response must be a map, got %s", resExpr.Type())
	}

	statusExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if statusExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("http-response-status: status must be a number, got %s", statusExpr.Type())
	}

	resMap := resExpr.(*mir.MapExpression)
	resMap.Values["status"] = statusExpr

	return resMap, nil
}

// __httpResponseHeader sets a response header
func __httpResponseHeader(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 3 {
		return nil, fmt.Errorf("http-response-header: expected 3 parameters (response, key, value), got %d", len(params))
	}

	resExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if resExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-response-header: response must be a map, got %s", resExpr.Type())
	}

	keyExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if keyExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("http-response-header: key must be a string, got %s", keyExpr.Type())
	}

	valueExpr, err := evaluate(params[2], session)
	if err != nil {
		return nil, err
	}

	if valueExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("http-response-header: value must be a string, got %s", valueExpr.Type())
	}

	resMap := resExpr.(*mir.MapExpression)

	// Get or create headers map
	var headersMap *mir.MapExpression
	if headers, ok := resMap.Values["headers"]; ok && headers.Type() == mir.MapExpr {
		headersMap = headers.(*mir.MapExpression)
	} else {
		headersMap, _ = mir.CreateMapExpression(make(map[string]mir.Expression))
		resMap.Values["headers"] = headersMap
	}

	headersMap.Values[keyExpr.(*mir.StringExpression).Val] = valueExpr

	return resMap, nil
}

// __httpResponseBody sets response body
func __httpResponseBody(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("http-response-body: expected 2 parameters (response, body), got %d", len(params))
	}

	resExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if resExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-response-body: response must be a map, got %s", resExpr.Type())
	}

	bodyExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if bodyExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("http-response-body: body must be a string, got %s", bodyExpr.Type())
	}

	resMap := resExpr.(*mir.MapExpression)
	resMap.Values["body"] = bodyExpr

	return resMap, nil
}

// __httpResponseJson sends a JSON response
func __httpResponseJson(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("http-response-json: expected 2 parameters (response, jsonData), got %d", len(params))
	}

	resExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if resExpr.Type() != mir.MapExpr {
		return nil, fmt.Errorf("http-response-json: response must be a map, got %s", resExpr.Type())
	}

	jsonExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	var jsonStr string
	if jsonExpr.Type() == mir.StringExpr {
		jsonStr = jsonExpr.(*mir.StringExpression).Val
	} else if jsonExpr.Type() == mir.MapExpr {
		// Convert map to JSON
		jsonMap := mapExpressionToGoMap(jsonExpr.(*mir.MapExpression))
		jsonBytes, err := json.Marshal(jsonMap)
		if err != nil {
			return nil, fmt.Errorf("http-response-json: error marshaling JSON: %s", err)
		}
		jsonStr = string(jsonBytes)
	} else {
		return nil, fmt.Errorf("http-response-json: jsonData must be a string or map, got %s", jsonExpr.Type())
	}

	resMap := resExpr.(*mir.MapExpression)
	resMap.Values["body"], _ = mir.CreateStringExpression(jsonStr)

	// Set Content-Type header
	var headersMap *mir.MapExpression
	if headers, ok := resMap.Values["headers"]; ok && headers.Type() == mir.MapExpr {
		headersMap = headers.(*mir.MapExpression)
	} else {
		headersMap, _ = mir.CreateMapExpression(make(map[string]mir.Expression))
		resMap.Values["headers"] = headersMap
	}

	contentTypeExpr, _ := mir.CreateStringExpression("application/json")
	headersMap.Values["Content-Type"] = contentTypeExpr

	return resMap, nil
}

// Helper functions

func createRequestMap(r *http.Request) (*mir.MapExpression, error) {
	reqMap := make(map[string]mir.Expression)

	// Method
	method, _ := mir.CreateStringExpression(r.Method)
	reqMap["method"] = method

	// Path
	path, _ := mir.CreateStringExpression(r.URL.Path)
	reqMap["path"] = path

	// Headers
	headersMap := make(map[string]mir.Expression)
	for key, values := range r.Header {
		if len(values) > 0 {
			headerVal, _ := mir.CreateStringExpression(values[0])
			headersMap[key] = headerVal
		}
	}
	headers, _ := mir.CreateMapExpression(headersMap)
	reqMap["headers"] = headers

	// Query parameters
	queryMap := make(map[string]mir.Expression)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			queryVal, _ := mir.CreateStringExpression(values[0])
			queryMap[key] = queryVal
		}
	}
	query, _ := mir.CreateMapExpression(queryMap)
	reqMap["query"] = query

	// Body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	body, _ := mir.CreateStringExpression(string(bodyBytes))
	reqMap["body"] = body

	return mir.CreateMapExpression(reqMap)
}

func createResponseMap() (*mir.MapExpression, error) {
	resMap := make(map[string]mir.Expression)

	// Default status 200
	status, _ := mir.CreateNumberExpression(200, 1)
	resMap["status"] = status

	// Empty headers
	headers, _ := mir.CreateMapExpression(make(map[string]mir.Expression))
	resMap["headers"] = headers

	// Empty body
	body, _ := mir.CreateStringExpression("")
	resMap["body"] = body

	return mir.CreateMapExpression(resMap)
}

func writeResponse(w http.ResponseWriter, resMap *mir.MapExpression) {
	// Set status code
	status := 200
	if statusExpr, ok := resMap.Values["status"]; ok && statusExpr.Type() == mir.NumberExpr {
		status = int(statusExpr.(*mir.NumberExpression).Numerator)
	}

	// Set headers
	if headersExpr, ok := resMap.Values["headers"]; ok && headersExpr.Type() == mir.MapExpr {
		headersMap := headersExpr.(*mir.MapExpression)
		for key, valueExpr := range headersMap.Values {
			if valueExpr.Type() == mir.StringExpr {
				w.Header().Set(key, valueExpr.(*mir.StringExpression).Val)
			}
		}
	}

	w.WriteHeader(status)

	// Write body
	if bodyExpr, ok := resMap.Values["body"]; ok && bodyExpr.Type() == mir.StringExpr {
		w.Write([]byte(bodyExpr.(*mir.StringExpression).Val))
	}
}

func callFunction(fn *mir.FunctionExpression, params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	arity := len(params)
	fnArity, ok := fn.Arities[arity]
	if !ok {
		return nil, fmt.Errorf("function %s does not support arity %d", fn.Identifier, arity)
	}

	// Create new scope for function
	session.CallStack.Push(params)
	defer session.CallStack.Pop()

	// Bind parameters
	for i, paramName := range fnArity.Parameters() {
		session.AddObject(paramName, params[i])
	}

	// Evaluate function body
	return evaluate(fnArity.Body(), session)
}

func mapExpressionToGoMap(mapExpr *mir.MapExpression) map[string]interface{} {
	if mapExpr == nil {
		return make(map[string]interface{})
	}

	result := make(map[string]interface{})
	for key, expr := range mapExpr.Values {
		if expr == nil {
			result[key] = nil
			continue
		}

		switch expr.Type() {
		case mir.NilExpr:
			result[key] = nil
		case mir.StringExpr:
			result[key] = expr.(*mir.StringExpression).Val
		case mir.NumberExpr:
			numExpr := expr.(*mir.NumberExpression)
			if numExpr.Denominator == 1 {
				result[key] = numExpr.Numerator
			} else {
				result[key] = float64(numExpr.Numerator) / float64(numExpr.Denominator)
			}
		case mir.BooleanExpr:
			result[key] = expr.(*mir.BooleanExpression).Val
		case mir.MapExpr:
			result[key] = mapExpressionToGoMap(expr.(*mir.MapExpression))
		case mir.ListExpr:
			listExpr := expr.(*mir.ListExpression)
			listResult := make([]interface{}, 0)
			for _, item := range listExpr.Expressions {
				if item == nil {
					listResult = append(listResult, nil)
					continue
				}
				switch item.Type() {
				case mir.NilExpr:
					listResult = append(listResult, nil)
				case mir.StringExpr:
					listResult = append(listResult, item.(*mir.StringExpression).Val)
				case mir.NumberExpr:
					numExpr := item.(*mir.NumberExpression)
					if numExpr.Denominator == 1 {
						listResult = append(listResult, numExpr.Numerator)
					} else {
						listResult = append(listResult, float64(numExpr.Numerator)/float64(numExpr.Denominator))
					}
				case mir.BooleanExpr:
					listResult = append(listResult, item.(*mir.BooleanExpression).Val)
				case mir.MapExpr:
					listResult = append(listResult, mapExpressionToGoMap(item.(*mir.MapExpression)))
				default:
					listResult = append(listResult, nil)
				}
			}
			result[key] = listResult
		default:
			result[key] = expr.Value()
		}
	}
	return result
}
