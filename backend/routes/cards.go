package routes

import (
	"github.com/gin-gonic/gin"
	"main.go/controllers"
)

func setupCardRoutes(baseRouter *gin.RouterGroup) {
	posts := baseRouter.Group("/cards")

	// Fetches tarot card deck from the external API.
	// Selects three random tarot cards from the deck, ensuring no duplicates.
	// Determines whether each card is reversed.
	// Constructs a JSON response containing the details of the three selected cards.
	// Sends the JSON response to the frontend for UI rendering.
	posts.GET("", controllers.GetandInterpretThreeCards)

	// The client makes a request to this endpoint with a specific uuid.
	// The function retrieves the interpretation associated with that uuid from local storage.
	// If the interpretation is found, it is returned in a JSON response.
	// If the interpretation is not found, an error message is sent back to the client.
	posts.GET("/interpret/:uuid", controllers.GetInterpretation)
}
