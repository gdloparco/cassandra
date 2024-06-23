package controllers

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"main.go/errors"
	"main.go/models"
	"main.go/services"
)

// Local storage for interpretations and UUIDs
var LocalStorage map[string]string = make(map[string]string)

// Function to select 3 random tarot cards from the deck
func GetRandomCard(deck []models.Card, currentCards []models.Card) models.Card {
	// Initialize a random number generator with the current time as the seed
	randomiser := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// Select a random index
		randomIndex := randomiser.Intn(len(deck))
		// Get the card at the random index
		randomCard := deck[randomIndex]

		// Check if the card is a duplicate
		isDuplicate := false
		for _, card := range currentCards {
			if card.CardName == randomCard.CardName {
				isDuplicate = true
				break
			}
		}

		// Return the card if it's not a duplicate
		if !isDuplicate {
			return randomCard
		}
	}
}

// Function to get and interpret 3 tarot cards
func GetandInterpretThreeCards(ctx *gin.Context) {
	// call below returns a type of []Card containing the whole deck from the API
	deck, err := services.FetchTarotCards() 
	if err != nil {
		errors.SendInternalError(ctx, err)
		return
	}

	// Generate a new UUID for the request
	requestID := uuid.New()
	var threeCards []models.Card
	threeCards = append(threeCards, GetRandomCard(deck, threeCards))
	threeCards = append(threeCards, GetRandomCard(deck, threeCards))
	threeCards = append(threeCards, GetRandomCard(deck, threeCards))

	// from here below we convert the three Card objects into three JSONCard objects

	var jsonCards []models.JSONCard
	var cardNames []string

	for _, card := range threeCards {
		//decide if card is reversed or not
		reversed := ReverseRandomiser()

		//edit the title with (Reversed) if applicable
		var FinalCardName string
		if reversed {
			FinalCardName = card.CardName + " (Reversed)"
		} else {
			FinalCardName = card.CardName
		}

		jsonCards = append(jsonCards, models.JSONCard{
			CardName:       FinalCardName,
			Type:           card.Type,
			MeaningUp:      card.MeaningUp,
			MeaningReverse: card.MeaningReverse,
			Description:    card.Description,
			ImageName:      card.ShortName + ".jpg",
			Reversed:       reversed,
		})

		var reversedValue string
		card.Reversed = reversed
		if card.Reversed {
			reversedValue = "(Reversed)"
		} else {
			reversedValue = ""
		}

		cardNames = append(cardNames, card.CardName, reversedValue)
	}

	//here we send our three Cards and the requestID in JSON form to the client, to be rendered in the UI.

	ctx.JSON(http.StatusOK, gin.H{"cards": jsonCards, "requestID": requestID})
	userStory := ctx.Query("userstory")
	userName := ctx.Query("name")

	// here we use Open AI's API to generate a reading of our three cards, we store this reading locally to return it to the user later.
	go func() {
		testing := os.Getenv("TESTING")
		if testing == "True" {
			interpretation := "This is a test interpretation"
			LocalStorage[requestID.String()] = interpretation
			return
		}
		// Get the OpenAI API key from environment variables
		apiKey := os.Getenv("API_KEY")
		interpretation, err := services.InterpretTarotCards(apiKey, cardNames, requestID, userStory, userName)
		if err != nil {
			errors.SendInternalError(ctx, err)
			return
		}
		// Store the interpretation in local storage
		LocalStorage[requestID.String()] = interpretation
		GetInterpretation(ctx)
	}()
}

// function to send the interpretation from the internal storage to the frontend
func GetInterpretation(ctx *gin.Context) {
	// Get the UUID from the request parameters
	requestID := ctx.Param("uuid")

	// Retrieve the interpretation from local storage
	interpretation, ok := LocalStorage[requestID]

	// Check if the interpretation was found
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No interpretation found for this UUID"})
		return
	}

	// Return the interpretation
	ctx.JSON(http.StatusOK, gin.H{"interpretation": interpretation})
}

// function to generate reversed cards at random
func ReverseRandomiser() bool {
	// Initialize a random number generator with the current time as the seed
	randomiser := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate a random number, either 0 or 1
	randomBool := randomiser.Intn(2)
	return randomBool == 0
}
