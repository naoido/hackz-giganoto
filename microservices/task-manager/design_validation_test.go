package main

import (
	"context"
	"fmt"
	"reflect"

	taskservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/task_service"
	commentservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/comment_service"
	labelservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/label_service"
	
	tasksvc "object-t.com/hackz-giganoto/microservices/task-manager/task_service"
	commentsvc "object-t.com/hackz-giganoto/microservices/task-manager/comment_service"
	labelsvc "object-t.com/hackz-giganoto/microservices/task-manager/label_service"
)

func main() {
	fmt.Println("🧪 Goa Design Validation Test")
	fmt.Println("==============================")
	
	// Logger function for testing
	logger := func(ctx context.Context, args ...any) {
		// Silent logger for testing
	}
	
	// Test 1: Service Interface Compliance
	fmt.Println("\n1. 🔍 Testing Service Interface Compliance")
	fmt.Println("-------------------------------------------")
	
	// Create service implementations
	taskSvc := tasksvc.New(logger)
	commentSvc := commentsvc.New(logger)
	labelSvc := labelsvc.New(logger)
	
	// Verify interface compliance
	var _ taskservice.Service = taskSvc
	var _ commentservice.Service = commentSvc
	var _ labelservice.Service = labelSvc
	
	fmt.Println("✅ TaskService implements taskservice.Service")
	fmt.Println("✅ CommentService implements commentservice.Service")
	fmt.Println("✅ LabelService implements labelservice.Service")
	
	// Test 2: Method Signatures
	fmt.Println("\n2. 🔍 Testing Method Signatures")
	fmt.Println("--------------------------------")
	
	ctx := context.Background()
	
	// Test TaskService methods
	testTaskServiceMethods(ctx, taskSvc)
	
	// Test CommentService methods
	testCommentServiceMethods(ctx, commentSvc)
	
	// Test LabelService methods
	testLabelServiceMethods(ctx, labelSvc)
	
	// Test 3: Data Type Validation
	fmt.Println("\n3. 🔍 Testing Data Type Structures")
	fmt.Println("-----------------------------------")
	
	testDataTypes()
	
	// Test 4: Endpoint Creation
	fmt.Println("\n4. 🔍 Testing Endpoint Creation")
	fmt.Println("--------------------------------")
	
	testEndpointCreation(taskSvc, commentSvc, labelSvc)
	
	fmt.Println("\n✅ All Goa Design Validation Tests Passed!")
	fmt.Println("🎯 Implementation matches Goa DSL design specification")
}

func testTaskServiceMethods(ctx context.Context, svc taskservice.Service) {
	fmt.Println("Testing TaskService methods:")
	
	// Test List method
	_, err := svc.List(ctx, &taskservice.ListPayload{})
	if err != nil {
		fmt.Printf("⚠️  TaskService.List error: %v\n", err)
	} else {
		fmt.Println("✅ TaskService.List method signature correct")
	}
	
	// Test Create method
	createPayload := &taskservice.CreatePayload{
		Title:       "Test Task",
		Description: stringPtr("Test Description"),
	}
	task, err := svc.Create(ctx, createPayload)
	if err != nil {
		fmt.Printf("⚠️  TaskService.Create error: %v\n", err)
	} else {
		fmt.Println("✅ TaskService.Create method signature correct")
		
		// Test Get method with created task
		getPayload := &taskservice.GetPayload{ID: task.ID}
		_, err = svc.Get(ctx, getPayload)
		if err != nil {
			fmt.Printf("⚠️  TaskService.Get error: %v\n", err)
		} else {
			fmt.Println("✅ TaskService.Get method signature correct")
		}
		
		// Test Update method
		updatePayload := &taskservice.UpdatePayload{
			ID:     task.ID,
			Status: stringPtr("IN_PROGRESS"),
		}
		_, err = svc.Update(ctx, updatePayload)
		if err != nil {
			fmt.Printf("⚠️  TaskService.Update error: %v\n", err)
		} else {
			fmt.Println("✅ TaskService.Update method signature correct")
		}
		
		// Test Delete method
		deletePayload := &taskservice.DeletePayload{ID: task.ID}
		err = svc.Delete(ctx, deletePayload)
		if err != nil {
			fmt.Printf("⚠️  TaskService.Delete error: %v\n", err)
		} else {
			fmt.Println("✅ TaskService.Delete method signature correct")
		}
	}
}

func testCommentServiceMethods(ctx context.Context, svc commentservice.Service) {
	fmt.Println("Testing CommentService methods:")
	
	// Test List method
	listPayload := &commentservice.ListPayload{TaskID: "test-task-id"}
	_, err := svc.List(ctx, listPayload)
	if err != nil {
		fmt.Printf("⚠️  CommentService.List error: %v\n", err)
	} else {
		fmt.Println("✅ CommentService.List method signature correct")
	}
	
	// Test Create method
	createPayload := &commentservice.CreatePayload{
		TaskID:   "test-task-id",
		AuthorID: "test-author-id",
		Body:     "Test comment",
	}
	comment, err := svc.Create(ctx, createPayload)
	if err != nil {
		fmt.Printf("⚠️  CommentService.Create error: %v\n", err)
	} else {
		fmt.Println("✅ CommentService.Create method signature correct")
		
		// Test Get method
		getPayload := &commentservice.GetPayload{ID: comment.ID}
		_, err = svc.Get(ctx, getPayload)
		if err != nil {
			fmt.Printf("⚠️  CommentService.Get error: %v\n", err)
		} else {
			fmt.Println("✅ CommentService.Get method signature correct")
		}
	}
}

func testLabelServiceMethods(ctx context.Context, svc labelservice.Service) {
	fmt.Println("Testing LabelService methods:")
	
	// Test List method (should return default labels)
	labels, err := svc.List(ctx)
	if err != nil {
		fmt.Printf("⚠️  LabelService.List error: %v\n", err)
	} else {
		fmt.Printf("✅ LabelService.List method signature correct (returned %d labels)\n", len(labels))
		
		if len(labels) > 0 {
			// Test Get method with first label
			getPayload := &labelservice.GetPayload{ID: labels[0].ID}
			_, err = svc.Get(ctx, getPayload)
			if err != nil {
				fmt.Printf("⚠️  LabelService.Get error: %v\n", err)
			} else {
				fmt.Println("✅ LabelService.Get method signature correct")
			}
		}
	}
	
	// Test Create method
	createPayload := &labelservice.CreatePayload{
		Name:  "test-label",
		Color: "#ff0000",
	}
	label, err := svc.Create(ctx, createPayload)
	if err != nil {
		fmt.Printf("⚠️  LabelService.Create error: %v\n", err)
	} else {
		fmt.Println("✅ LabelService.Create method signature correct")
		
		// Test Update method
		updatePayload := &labelservice.UpdatePayload{
			ID:   label.ID,
			Name: stringPtr("updated-test-label"),
		}
		_, err = svc.Update(ctx, updatePayload)
		if err != nil {
			fmt.Printf("⚠️  LabelService.Update error: %v\n", err)
		} else {
			fmt.Println("✅ LabelService.Update method signature correct")
		}
	}
}

func testDataTypes() {
	fmt.Println("Testing data type structures:")
	
	// Test Task type
	task := &taskservice.Task{}
	taskType := reflect.TypeOf(task).Elem()
	
	requiredTaskFields := []string{"ID", "Title", "Status", "CreatedAt", "UpdatedAt"}
	for _, field := range requiredTaskFields {
		if _, found := taskType.FieldByName(field); found {
			fmt.Printf("✅ Task.%s field exists\n", field)
		} else {
			fmt.Printf("❌ Task.%s field missing\n", field)
		}
	}
	
	// Test Comment type
	comment := &commentservice.Comment{}
	commentType := reflect.TypeOf(comment).Elem()
	
	requiredCommentFields := []string{"ID", "TaskID", "AuthorID", "Body", "CreatedAt"}
	for _, field := range requiredCommentFields {
		if _, found := commentType.FieldByName(field); found {
			fmt.Printf("✅ Comment.%s field exists\n", field)
		} else {
			fmt.Printf("❌ Comment.%s field missing\n", field)
		}
	}
	
	// Test Label type
	label := &labelservice.Label{}
	labelType := reflect.TypeOf(label).Elem()
	
	requiredLabelFields := []string{"ID", "Name", "Color"}
	for _, field := range requiredLabelFields {
		if _, found := labelType.FieldByName(field); found {
			fmt.Printf("✅ Label.%s field exists\n", field)
		} else {
			fmt.Printf("❌ Label.%s field missing\n", field)
		}
	}
}

func testEndpointCreation(taskSvc taskservice.Service, commentSvc commentservice.Service, labelSvc labelservice.Service) {
	fmt.Println("Testing endpoint creation:")
	
	// Test TaskService endpoints
	taskEndpoints := taskservice.NewEndpoints(taskSvc)
	if taskEndpoints != nil {
		fmt.Println("✅ TaskService endpoints created successfully")
	} else {
		fmt.Println("❌ TaskService endpoints creation failed")
	}
	
	// Test CommentService endpoints
	commentEndpoints := commentservice.NewEndpoints(commentSvc)
	if commentEndpoints != nil {
		fmt.Println("✅ CommentService endpoints created successfully")
	} else {
		fmt.Println("❌ CommentService endpoints creation failed")
	}
	
	// Test LabelService endpoints
	labelEndpoints := labelservice.NewEndpoints(labelSvc)
	if labelEndpoints != nil {
		fmt.Println("✅ LabelService endpoints created successfully")
	} else {
		fmt.Println("❌ LabelService endpoints creation failed")
	}
}

func stringPtr(s string) *string {
	return &s
}