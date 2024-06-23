package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main.go/env"
	"main.go/routes"
)

func main() {
	env.LoadEnv()
	// Load environment variables using the custom env package
	mode := os.Getenv("GIN_MODE")

	if mode == "" { // If GIN_MODE is not set to debug
		// Default to release mode
		mode = gin.ReleaseMode
	}

	gin.SetMode(mode)

	// Call the setupApp function to initialize the Gin engine
	app := setupApp()

	// Start the Gin server on port 8082
	app.Run(":8082")
}

func setupApp() *gin.Engine {
	// Create a new Gin engine with default middleware
	app := gin.Default()
	// Call setupCORS to configure CORS settings
	setupCORS(app)
	// Call SetupRoutes from the routes package to define application routes
	routes.SetupRoutes(app)
	// Return the configured Gin engine
	return app
}

func setupCORS(app *gin.Engine) {
	// CORS is a mechanism that allows servers to specify which origins are permitted to access their resources.
	// It's essential for security,as it prevents unauthorized access to sensitive data and APIs across different
	// domains while still enabling legitimate cross-origin interactions when properly configured.

	// Get the default CORS configuration
	config := cors.DefaultConfig()
	// Allow all origins to access the server
	config.AllowAllOrigins = true
	// Allow credentials (cookies, authorization headers, etc.) in CORS requests
	config.AllowCredentials = true
	// Specify allowed headers in CORS requests
	config.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type", "Accept"}

	// Apply the CORS middleware with the configured settings to the Gin engine
	app.Use(cors.New(config))
}
