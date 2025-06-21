#!/bin/bash

echo "🔨 Task Manager Build Test"
echo "=========================="

# Change to the project directory
cd /Users/kyiku/GolandProjects/hackz-giganoto/microservices/task-manager

echo "📁 Current directory: $(pwd)"
echo ""

echo "🔍 Checking Go version..."
go version
echo ""

echo "📦 Running go mod tidy..."
go mod tidy
echo ""

echo "🔨 Attempting to build..."
if go build -o task-manager-bin cmd/task_manager/main.go; then
    echo "✅ Build successful!"
    echo "🎯 Binary created: task-manager-bin"
    echo ""
    echo "🚀 To run the server:"
    echo "   ./task-manager-bin"
    echo "   or"
    echo "   go run cmd/task_manager/main.go"
else
    echo "❌ Build failed. Checking for specific errors..."
    echo ""
    echo "🔍 Running go build with verbose output:"
    go build -v cmd/task_manager/main.go 2>&1 | head -20
fi

echo ""
echo "📋 Project structure check:"
ls -la cmd/task_manager/
ls -la task_service/
ls -la comment_service/
ls -la label_service/