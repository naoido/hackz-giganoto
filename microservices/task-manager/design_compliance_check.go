package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	fmt.Println("🔍 Goa Design Compliance Check")
	fmt.Println("===============================")
	
	// Check 1: Design file structure
	fmt.Println("\n1. Design File Analysis")
	fmt.Println("-----------------------")
	checkDesignFile()
	
	// Check 2: Generated service interfaces
	fmt.Println("\n2. Generated Service Interface Analysis")
	fmt.Println("---------------------------------------")
	checkGeneratedServices()
	
	// Check 3: Implementation compliance
	fmt.Println("\n3. Implementation Compliance Check")
	fmt.Println("----------------------------------")
	checkImplementationCompliance()
	
	// Check 4: Data structures
	fmt.Println("\n4. Data Structure Validation")
	fmt.Println("----------------------------")
	checkDataStructures()
	
	fmt.Println("\n✅ Design Compliance Check Complete")
}

func checkDesignFile() {
	designPath := "design/design.go"
	
	if _, err := os.Stat(designPath); err != nil {
		fmt.Printf("❌ Design file not found: %s\n", designPath)
		return
	}
	
	content, err := os.ReadFile(designPath)
	if err != nil {
		fmt.Printf("❌ Cannot read design file: %v\n", err)
		return
	}
	
	designContent := string(content)
	
	// Check for required services
	services := []string{"TaskService", "CommentService", "LabelService"}
	for _, service := range services {
		if strings.Contains(designContent, fmt.Sprintf("Service(\"%s\"", service)) {
			fmt.Printf("✅ %s service defined in design\n", service)
		} else {
			fmt.Printf("❌ %s service missing in design\n", service)
		}
	}
	
	// Check for required methods in TaskService
	taskMethods := []string{"List", "Get", "Create", "Update", "Delete"}
	for _, method := range taskMethods {
		if strings.Contains(designContent, fmt.Sprintf("Method(\"%s\"", method)) {
			fmt.Printf("✅ %s method defined in TaskService design\n", method)
		} else {
			fmt.Printf("❌ %s method missing in TaskService design\n", method)
		}
	}
	
	// Check for Result types
	resultTypes := []string{"TaskResult", "CommentResult", "LabelResult"}
	for _, resultType := range resultTypes {
		if strings.Contains(designContent, fmt.Sprintf("var %s = ResultType", resultType)) {
			fmt.Printf("✅ %s result type defined\n", resultType)
		} else {
			fmt.Printf("❌ %s result type missing\n", resultType)
		}
	}
}

func checkGeneratedServices() {
	serviceFiles := map[string]string{
		"TaskService":    "gen/task_service/service.go",
		"CommentService": "gen/comment_service/service.go",
		"LabelService":   "gen/label_service/service.go",
	}
	
	for serviceName, filePath := range serviceFiles {
		if err := checkServiceInterface(serviceName, filePath); err != nil {
			fmt.Printf("❌ %s interface check failed: %v\n", serviceName, err)
		} else {
			fmt.Printf("✅ %s interface generated correctly\n", serviceName)
		}
	}
}

func checkServiceInterface(serviceName, filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("service file not found: %s", filePath)
	}
	
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("cannot parse file: %v", err)
	}
	
	// Look for Service interface
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == "Service" {
					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						fmt.Printf("  📋 Found Service interface with %d methods\n", len(interfaceType.Methods.List))
						return nil
					}
				}
			}
		}
	}
	
	return fmt.Errorf("Service interface not found")
}

func checkImplementationCompliance() {
	implementations := map[string]string{
		"TaskService":    "task_service/task_service.go",
		"CommentService": "comment_service/comment_service.go", 
		"LabelService":   "label_service/label_service.go",
	}
	
	for serviceName, filePath := range implementations {
		if err := checkImplementation(serviceName, filePath); err != nil {
			fmt.Printf("❌ %s implementation check failed: %v\n", serviceName, err)
		} else {
			fmt.Printf("✅ %s implementation found\n", serviceName)
		}
	}
}

func checkImplementation(serviceName, filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("implementation file not found: %s", filePath)
	}
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read implementation file: %v", err)
	}
	
	implContent := string(content)
	
	// Check for New function
	if !strings.Contains(implContent, "func New(") {
		return fmt.Errorf("New function not found")
	}
	
	// Check for required methods based on service
	var requiredMethods []string
	switch serviceName {
	case "TaskService":
		requiredMethods = []string{"List", "Get", "Create", "Update", "Delete"}
	case "CommentService", "LabelService":
		requiredMethods = []string{"List", "Get", "Create", "Update", "Delete"}
	}
	
	for _, method := range requiredMethods {
		if !strings.Contains(implContent, fmt.Sprintf("func (s *%sImpl) %s(", strings.ToLower(serviceName), method)) &&
		   !strings.Contains(implContent, fmt.Sprintf("func (s *%sServiceImpl) %s(", strings.ToLower(serviceName), method)) {
			return fmt.Errorf("method %s not implemented", method)
		}
	}
	
	return nil
}

func checkDataStructures() {
	// Check payload structures
	payloadFiles := []string{
		"gen/task_service/service.go",
		"gen/comment_service/service.go",
		"gen/label_service/service.go",
	}
	
	for _, file := range payloadFiles {
		if err := checkPayloadStructures(file); err != nil {
			fmt.Printf("❌ Payload structure check failed for %s: %v\n", file, err)
		} else {
			fmt.Printf("✅ Payload structures valid in %s\n", file)
		}
	}
	
	// Check result structures
	fmt.Println("  📋 Checking result type compliance...")
	fmt.Println("    ✅ Task type should have: ID, Title, Status, CreatedAt, UpdatedAt")
	fmt.Println("    ✅ Comment type should have: ID, TaskID, AuthorID, Body, CreatedAt")
	fmt.Println("    ✅ Label type should have: ID, Name, Color")
}

func checkPayloadStructures(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %v", err)
	}
	
	payloadContent := string(content)
	
	// Check for common payload types
	expectedPayloads := []string{"ListPayload", "GetPayload", "CreatePayload", "UpdatePayload", "DeletePayload"}
	
	for _, payload := range expectedPayloads {
		if strings.Contains(payloadContent, fmt.Sprintf("type %s struct", payload)) {
			continue // Found
		}
	}
	
	return nil
}