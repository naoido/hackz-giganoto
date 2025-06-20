package design

import (
	. "goa.design/goa/v3/dsl"
)

// Service describes a service
var _ = Service("profile", func() {
	Description("The profile service allows management of user profiles.")

	// Method describes a service method (endpoint)
	Method("create", func() {
		// Payload describes the method payload
		// Here, we define a ProfilePayload type
		Payload(ProfilePayload)
		// Result describes the method result
		// Here, we define a Profile type
		Result(Profile)
		// HTTP describes the HTTP transport mapping
		HTTP(func() {
			// Requests to the service consist of POST requests
			// The payload is encoded in the request body
			POST("/profile")
			// Responses use a 201 Created status code
			Response(StatusCreated)
		})
	})

	Method("get", func() {
		Payload(func() {
			Attribute("id", String, "Profile ID")
			Required("id")
		})
		Result(Profile)
		HTTP(func() {
			GET("/profile/{id}")
			Response(StatusOK)
		})
	})
})

// ProfilePayload defines the data structure for creating a profile.
var ProfilePayload = Type("ProfilePayload", func() {
	Description("Payload for creating a user profile.")
	Attribute("username", String, "The user's name")
	Required("username")
})

// Profile defines the data structure for a profile.
var Profile = ResultType("application/vnd.goa.example.profile", func() {
	Description("A user profile")
	Attribute("id", String, "Unique profile ID")
	Attribute("username", String, "The user's name")
	Required("id", "username")
})
