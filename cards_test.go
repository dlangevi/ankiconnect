package ankiconnect

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCardsManager_Get(t *testing.T) {
	findCardsPayload := []byte(`{
    "action": "findCards",
    "version": 6,
    "params": {
        "query": "deck:current"
    }
  }`)

	cardsInfoPayload := []byte(`{
    "action": "cardsInfo",
    "version": 6,
    "params": {
        "cards": [1498938915662, 1502098034048]
    }
  }`)

	t.Run("success", func(t *testing.T) {
		defer httpmock.Reset()

		registerMultipleVerifiedPayloads(t,
			[][2][]byte{
				// Get will do two api calls, first findCards to get the card id's
				{
					findCardsPayload,
					loadTestResult(t, ActionFindCards),
				},
				// Then cardsInfo to transform those into actual anki cards
				{
					cardsInfoPayload,
					loadTestResult(t, ActionCardsInfo),
				},
			})

		payload := "deck:current"
		notes, restErr := client.Cards.Get(payload)
		assert.Nil(t, restErr)
		assert.Equal(t, len(*notes), 2)

	})

	t.Run("errorFailSearch", func(t *testing.T) {
		defer httpmock.Reset()

		registerErrorResponse(t)

		_, restErr := client.Cards.Get("deck:current")
		assert.NotNil(t, restErr)
		assert.Equal(t, http.StatusBadRequest, restErr.StatusCode)
		assert.Equal(t, "some error message", restErr.Message)
	})
}

func TestCardsManager_GetReviews(t *testing.T) {
	findCardsPayload := []byte(`{
    "action": "findCards",
    "version": 6,
    "params": {
        "query": "deck:current"
    }
  }`)

	findCardsResult := []byte(`{
    "result": [1653613948202],
    "error": null
  }`)

	getCardReviewsPayload := []byte(`{
    "action": "getReviewsOfCards",
    "version": 6,
    "params": {
        "cards": [1653613948202]
    }
  }`)

	getCardReviewsBadResult := []byte(`{
    "result": {
        "not a number": [
            {
                "id": 1653772912146,
                "usn": 1750,
                "ease": 1,
                "ivl": -20,
                "lastIvl": -20,
                "factor": 0,
                "time": 38192,
                "type": 0
            }
        ]
    },
    "error": null
  }`)

	t.Run("success", func(t *testing.T) {
		defer httpmock.Reset()

		registerMultipleVerifiedPayloads(t,
			[][2][]byte{
				// Get will do two api calls, first findCards to get the card id's
				{
					findCardsPayload,
					findCardsResult,
				},
				// Then cardsInfo to transform those into actual anki cards
				{
					getCardReviewsPayload,
					loadTestResult(t, ActionGetReviewsOfCards),
				},
			})

		payload := "deck:current"
		reviews, restErr := client.Cards.GetReviews(payload)
		assert.Nil(t, restErr)
		assert.Equal(t, len(*reviews), 1)
		reviewsForCard := (*reviews)[1653613948202]
		assert.Equal(t, len(reviewsForCard), 2)
		assert.Equal(t, reviewsForCard[0].ID, int64(1653772912146))
	})

	t.Run("errorBadJsonKey", func(t *testing.T) {
		defer httpmock.Reset()

		registerMultipleVerifiedPayloads(t,
			[][2][]byte{
				// Get will do two api calls, first findCards to get the card id's
				{
					findCardsPayload,
					findCardsResult,
				},
				// Then cardsInfo to transform those into actual anki cards
				{
					getCardReviewsPayload,
					getCardReviewsBadResult,
				},
			})

		_, restErr := client.Cards.GetReviews("deck:current")
		assert.NotNil(t, restErr)
		assert.Equal(t, http.StatusInternalServerError, restErr.StatusCode)
		assert.Equal(t, "Internal Server Error", restErr.Message)
	})

	t.Run("errorFailSearch", func(t *testing.T) {
		defer httpmock.Reset()

		registerErrorResponse(t)

		_, restErr := client.Cards.GetReviews("deck:current")
		assert.NotNil(t, restErr)
		assert.Equal(t, http.StatusBadRequest, restErr.StatusCode)
		assert.Equal(t, "some error message", restErr.Message)
	})

}

func TestCustomUnmarshall(t *testing.T) {
	t.Run("badFormattedJson", func(t *testing.T) {
		results := ResultGetReviews{}
		err := results.UnmarshalJSON([]byte(`{
      "cards": {
        "id": "not a number"
      }
    }`))
		assert.NotNil(t, err)
	})
}
