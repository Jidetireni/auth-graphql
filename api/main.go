package main

func main() {

	// Initialize the server
	server := NewServer()

	// graphql server
	server.SetGqlServer()

	// Mount the routes
	server.MountRoutes()

	// Start the server
	server.Start()
}
