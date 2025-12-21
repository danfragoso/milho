#!/bin/bash

cd "$(dirname "$0")/.."

echo "=== Testing Guestbook Application ==="
echo ""

echo "Creating first entry..."
curl -s -X POST http://localhost:8080/create_entry \
  -d 'name=Alice&message=Hello%20from%20the%20Milho%20guestbook!'
echo ""
echo ""

echo "Creating second entry..."
curl -s -X POST http://localhost:8080/create_entry \
  -d 'name=Bob&message=This%20is%20amazing!%20TCP%20+%20HTTP%20+%20KV%20=%20Win!'
echo ""
echo ""

echo "Creating third entry..."
curl -s -X POST http://localhost:8080/create_entry \
  -d 'name=Charlie&message=Milho%20rocks!'
echo ""
echo ""

echo "Fetching all entries..."
curl -s http://localhost:8080/entries
echo ""
