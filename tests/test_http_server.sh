#!/bin/bash

# Run from project root
cd "$(dirname "$0")/.."

# Start server in background
./milho examples/http.milho > /tmp/http_server_test.log 2>&1 &
SERVER_PID=$!
sleep 3

echo "=== Testing /hello endpoint ==="
curl -s http://localhost:8080/hello
echo ""

echo ""
echo "=== Testing /greet with query param ==="
curl -s "http://localhost:8080/greet?name=Milho"
echo ""

echo ""
echo "=== Testing POST /echo ==="
curl -s -X POST http://localhost:8080/echo -d "Test message"
echo ""

echo ""
echo "=== Testing /info ==="
curl -s http://localhost:8080/info
echo ""

# Stop server
kill $SERVER_PID 2>/dev/null

echo ""
echo "=== Server log ==="
cat /tmp/http_server_test.log
