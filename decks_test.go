package ankiconnect

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestDecksManager_GetAll(t *testing.T) {
	getAllRequest := []byte(`{
    "action": "deckNames",
    "version": 6
}`)
	getAllResult := []byte(`{
    "result": [
        "Default",
        "Deck01",
        "Deck02"
    ],
    "error": null
}`)
	t.Run("success", func(t *testing.T) {
		defer httpmock.Reset()

		registerVerifiedPayload(t, getAllRequest, getAllResult)

		decks, restErr := client.Decks.GetAll()
		assert.NotNil(t, decks)
		assert.Nil(t, restErr)
		assert.Equal(t, 3, len(*decks))
	})

	t.Run("error", func(t *testing.T) {
		defer httpmock.Reset()

		registerErrorResponse(t)

		decks, restErr := client.Decks.GetAll()
		assert.Nil(t, decks)
		assert.NotNil(t, restErr)
		assert.Equal(t, http.StatusBadRequest, restErr.StatusCode)
		assert.Equal(t, "some error message", restErr.Message)
	})

	t.Run("http request error", func(t *testing.T) {
		defer httpmock.Reset()

		decks, restErr := client.Decks.GetAll()
		assert.Nil(t, decks)
		assert.NotNil(t, restErr)
		assert.Equal(t, http.StatusInternalServerError, restErr.StatusCode)
		assert.Equal(t, "Internal Server Error", restErr.Message)
	})
}

func TestDecksManager_Create(t *testing.T) {
	createRequest := []byte(`{
    "action": "createDeck",
    "version": 6,
    "params": {
        "deck": "Japanese::Tokyo"
    }
}`)
	createResponse := []byte(`{
    "result": 1659294179522,
    "error": null
}`)

	t.Run("success", func(t *testing.T) {
		defer httpmock.Reset()

		registerVerifiedPayload(t, createRequest, createResponse)

		restErr := client.Decks.Create("Japanese::Tokyo")
		assert.Nil(t, restErr)
	})

	t.Run("error", func(t *testing.T) {
		defer httpmock.Reset()

		registerErrorResponse(t)

		restErr := client.Decks.Create("test")
		assert.NotNil(t, restErr)
		assert.Equal(t, http.StatusBadRequest, restErr.StatusCode)
		assert.Equal(t, "some error message", restErr.Message)
	})
}

func TestDecksManagerDelete(t *testing.T) {
	deleteDeckRequest := []byte(`{
    "action": "deleteDecks",
    "version": 6,
    "params": {
        "decks": ["test"],
        "cardsToo": true
    }
}`)

	t.Run("success", func(t *testing.T) {
		defer httpmock.Reset()

		registerVerifiedPayload(t, deleteDeckRequest, genericSuccessJson)

		restErr := client.Decks.Delete("test")
		assert.Nil(t, restErr)
	})

	t.Run("error", func(t *testing.T) {
		defer httpmock.Reset()

		registerErrorResponse(t)

		restErr := client.Decks.Delete("test")
		assert.NotNil(t, restErr)
		assert.Equal(t, http.StatusBadRequest, restErr.StatusCode)
		assert.Equal(t, "some error message", restErr.Message)
	})
}

func TestDecksManagerGetReviews(t *testing.T) {
	getReviewsRequest := []byte(`{
    "action": "cardReviews",
    "version": 6,
    "params": {
        "deck": "default",
        "startID": 1594194095740
    }
  }`)
	getReviewsResponse := []byte(`{
    "result": [
        [1594194095746, 1485369733217, -1, 3,   4, -60, 2500, 6157, 0],
        [1594201393292, 1485369902086, -1, 1, -60, -60,    0, 4846, 0]
    ],
    "error": null
  }`)

	t.Run("success", func(t *testing.T) {
		defer httpmock.Reset()

		registerVerifiedPayload(t, getReviewsRequest, getReviewsResponse)

		startTime := time.UnixMilli(1594194095740)
		reviews, restErr := client.Decks.GetReviewsAfter("default", startTime)
		assert.Nil(t, restErr)
		assert.Len(t, *reviews, 2)
		firstReview := (*reviews)[0]
		assert.Equal(t, firstReview.ReviewTime, int64(1594194095746))
		assert.Equal(t, firstReview.USN, int64(-1))
		assert.Equal(t, firstReview.PreviousInterval, int64(-60))
		assert.Equal(t, firstReview.ReviewDuration, int64(6157))
	})

	t.Run("error", func(t *testing.T) {
		defer httpmock.Reset()

		registerErrorResponse(t)

		startTime := time.UnixMilli(1594194095740)
		_, restErr := client.Decks.GetReviewsAfter("default", startTime)
		assert.NotNil(t, restErr)
		assert.Equal(t, http.StatusBadRequest, restErr.StatusCode)
		assert.Equal(t, "some error message", restErr.Message)
	})
}

func TestUnmarshallCardReview(t *testing.T) {
	t.Run("notEnoughFields", func(t *testing.T) {
		result := CardReview{}
		err := result.UnmarshalJSON([]byte(`[
       0,1,2,3,4
    ]`))
		assert.NotNil(t, err)
	})
	t.Run("badJson", func(t *testing.T) {
		result := CardReview{}
		err := result.UnmarshalJSON([]byte(`{
      bad json
    }`))
		assert.NotNil(t, err)
	})

}
