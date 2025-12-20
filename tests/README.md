# HTTP Server Tests

This folder contains tests for the Milho HTTP server implementation.

## Test Files

- **test_http_server.sh** - Automated integration tests for HTTP endpoints
- **test_http.milho** - Simple HTTP server example for manual testing
- **test_http_debug.milho** - HTTP server with debug output for troubleshooting

## Running Tests

### Automated Tests

Run the complete test suite:

```bash
./tests/test_http_server.sh
```

This will:
- Start the HTTP server in the background
- Test multiple endpoints (GET, POST, JSON, query params)
- Display results and server logs
- Clean up after completion

### Manual Testing

Start a simple test server:

```bash
./milho tests/test_http.milho
```

Then in another terminal:

```bash
curl http://localhost:8080/hello
```

### Debug Mode

For detailed output:

```bash
./milho tests/test_http_debug.milho
```

## Test Coverage

The automated tests verify:

✅ JSON responses (`/hello`)
✅ Query parameter handling (`/greet?name=Milho`)
✅ POST request body handling (`/echo`)
✅ Request metadata (`/info` - method, path, headers, query)

## Expected Results

All endpoints should return appropriate responses:

```json
# GET /hello
{"message":"Hello, World!","timestamp":null}

# GET /greet?name=Milho
{"greeting":"Hello, Milho!","query-params":{"name":"Milho"}}

# POST /echo -d "Test message"
Test message

# GET /info
{"headers":{...},"method":"GET","path":"/info","query":{}}
```
