#!/bin/bash

# Test script for meme generation with different modes

echo "Testing meme generation with different modes..."

# Set data directory
export DATA_DIR=./data

# Test 1: Default mode (synchronous)
echo "Test 1: Default mode (synchronous)"
unset GENERATE_MEME
go run cmd/web/main.go &
PID=$!
sleep 2

# Create a test template
curl -X POST http://localhost:8080/api/templates -H "Content-Type: application/json" -d '{"name":"test"}'

# Create a test meme
curl -X POST http://localhost:8080/api/memes -H "Content-Type: application/json" -d '{"template":"test","text_top":"Hello","text_bottom":"World"}'

kill $PID
echo "Test 1 completed"

# Test 2: Background mode
echo "Test 2: Background mode"
export GENERATE_MEME=background
go run cmd/web/main.go &
PID=$!
sleep 2

# Create a test meme
curl -X POST http://localhost:8080/api/memes -H "Content-Type: application/json" -d '{"template":"test","text_top":"Hello","text_bottom":"Background"}'

kill $PID
echo "Test 2 completed"

# Test 3: Container App Job mode
echo "Test 3: Container App Job mode"
export GENERATE_MEME=containerappjob
export PROJECT_ID=your-project-id
export CLOUDRU_KEY_ID=your-key-id
export CLOUDRU_KEY_SECRET=your-key-secret
go run cmd/web/main.go &
PID=$!
sleep 2

# Create a test meme
curl -X POST http://localhost:8080/api/memes -H "Content-Type: application/json" -d '{"template":"test","text_top":"Hello","text_bottom":"ContainerApp"}'

kill $PID
echo "Test 3 completed"

echo "All tests completed!"