package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"main.go/controllers"
	"main.go/env"
	"main.go/models"
)

// TestSuiteEnv is a struct that embeds suite.Suite and includes additional fields
// for testing HTTP requests with Gin
type TestSuiteEnv struct {
	suite.Suite
	app *gin.Engine
	// ResponseRecorder captures HTTP responses
	res *httptest.ResponseRecorder
}

// RequestSetup is a helper function to set up and execute HTTP requests for testing
func RequestSetup(app *gin.Engine, suite *TestSuiteEnv, reqType string, path string) []byte {
	req, _ := http.NewRequest(reqType, path, nil) // Create a new HTTP request
	app.ServeHTTP(suite.res, req) // Serve the request and capture the response
	responseData, _ := io.ReadAll(suite.res.Body) // Read the response body (type of []byte)
	return responseData 
}

// SetupSuite is run once before all tests in the suite
func (suite *TestSuiteEnv) SetupSuite() {
	// Load environment variables from .test.env file
	env.LoadEnv(".test.env") 
	// Set up the Gin application
	suite.app = setupApp()
}

func (suite *TestSuiteEnv) SetupTest() {
	// Create a new ResponseRecorder for each test
	suite.res = httptest.NewRecorder()
}

// This gets run automatically by `go test` (Entry Point) so we call `suite.Run` inside it
func TestSuite(t *testing.T) {
	// This is what actually runs our suite
	suite.Run(t, new(TestSuiteEnv))
}

// Checks that the response code is 200 for GET /cards
func (suite *TestSuiteEnv) Test_GetThreeCards_ResponseCode() {
	app := suite.app
	// Make a GET request to /cards
	responseData := RequestSetup(app, suite, "GET", "/cards")

	var jsonCards struct {
		Cards []models.JSONCard
	}

	 // Unmarshal the response into jsonCards
	_ = json.Unmarshal(responseData, &jsonCards)

	// Assert that the response code is 200
	assert.Equal(suite.T(), 200, suite.res.Code)
}

// Checks that the response contains exactly 3 cards
func (suite *TestSuiteEnv) Test_GetThreeCards_ExpectedFormat() {
	app := suite.app
	responseData := RequestSetup(app, suite, "GET", "/cards")

	var jsonCards struct {
		Cards []models.JSONCard
	}

	_ = json.Unmarshal(responseData, &jsonCards)

	// Assert that there are exactly 3 cards
	assert.Len(suite.T(), jsonCards.Cards, 3)
}

// Checks that two different requests return two different sets of cards
func (suite *TestSuiteEnv) Test_GetThreeCardsIsRandom() {
	app := suite.app

	//Response 1
	responseData := RequestSetup(app, suite, "GET", "/cards")
	var jsonCards struct {
		Cards []models.JSONCard
	}

	_ = json.Unmarshal(responseData, &jsonCards)

	//Response 2
	responseData2 := RequestSetup(app, suite, "GET", "/cards")
	var jsonCards2 struct {
		Cards []models.JSONCard
	}

	_ = json.Unmarshal(responseData2, &jsonCards2)

	assert.NotEqual(suite.T(), jsonCards.Cards[0].CardName, jsonCards2.Cards[0].CardName) // 0.0041% of failure chances
}

// Checks that the response code is 200 for GET /cards/interpret/:uuid
func (suite *TestSuiteEnv) Test_GetAndInterpretCards_ResponseCode() {
	app := suite.app

	// Get the cards
	responseData := RequestSetup(app, suite, "GET", "/cards")
	var jsonCards struct {
		Cards     []models.JSONCard
		RequestID uuid.UUID
	}

	_ = json.Unmarshal(responseData, &jsonCards)

	// Request the interpretation
	_ = RequestSetup(app, suite, "GET", "/cards/interpret/"+jsonCards.RequestID.String())

	// Assert that response code is 200
	assert.Equal(suite.T(), 200, suite.res.Code)
}

// Checks that the interpretation response is as expected
func (suite *TestSuiteEnv) Test_GetAndInterpretCards_ExpectedFormat() {
	app := suite.app

	// Get the cards
	responseData := RequestSetup(app, suite, "GET", "/cards")
	var jsonCards struct {
		Cards     []models.JSONCard
		RequestID uuid.UUID
	}

	_ = json.Unmarshal(responseData, &jsonCards)

	// Request the interpretation
	responseData2 := RequestSetup(app, suite, "GET", "/cards/interpret/"+jsonCards.RequestID.String())
	var interpretationResponse struct {
		Interpretation string
	}

	_ = json.Unmarshal(responseData2, &interpretationResponse)

	// Assert that the interpretation in LocalStorage matches the expected string
	assert.Equal(suite.T(), controllers.LocalStorage[jsonCards.RequestID.String()], "This is a test interpretation")
}
