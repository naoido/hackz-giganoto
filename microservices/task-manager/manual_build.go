// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("🔨 Manual Build Process")
	fmt.Println("========================")
	
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}
	fmt.Printf("Working directory: %s\n", wd)
	
	// Change to project directory
	projectDir := "/Users/kyiku/GolandProjects/hackz-giganoto/microservices/task-manager"
	if err := os.Chdir(projectDir); err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
		return
	}
	fmt.Printf("Changed to: %s\n", projectDir)
	
	// Check if key files exist
	files := []string{
		"go.mod",
		"simple_main.go", 
		"gen/task_service/service.go",
		"task_service/task_service.go",
	}
	
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			fmt.Printf("✅ %s exists\n", file)
		} else {
			fmt.Printf("❌ %s missing\n", file)
		}
	}
	
	fmt.Println("\n🔧 Attempting to build simple_main.go...")
	
	// Try to build
	cmd := exec.Command("go", "build", "-o", "task-manager-binary", "simple_main.go")
	cmd.Dir = projectDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Build failed: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
	} else {
		fmt.Printf("✅ Build successful!\n")
		
		// Check if binary was created
		if _, err := os.Stat(filepath.Join(projectDir, "task-manager-binary")); err == nil {
			fmt.Println("🎯 Binary 'task-manager-binary' created successfully!")
			fmt.Println("\n🚀 To run the server:")
			fmt.Println("   ./task-manager-binary")
		}
	}
}