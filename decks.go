package ankiconnect

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/privatesquare/bkst-go-utils/utils/errors"
)

const (
	ActionDeckNames    = "deckNames"
	ActionCreateDeck   = "createDeck"
	ActionGetDeckStats = "getDeckStats"
	ActionDeleteDecks  = "deleteDecks"
	// This is based on the deck so will put it here for now
	ActionCardReviews = "cardReviews"
)

type (
	// DecksManager describes the interface that can be used to perform operations on anki decks.
	DecksManager interface {
		GetAll() (*[]string, *errors.RestErr)
		Create(name string) *errors.RestErr
		Delete(name string) *errors.RestErr
		GetReviewsAfter(name string, startTime time.Time) (*ResultCardReviews, *errors.RestErr)
	}

	// ParamsCreateDeck represents the ankiconnect API params required for creating a new deck.
	ParamsCreateDeck struct {
		Deck string `json:"deck,omitempty"`
	}

	// ParamsDeleteDeck represents the ankiconnect API params required for deleting one or more decks
	ParamsDeleteDecks struct {
		Decks    *[]string `json:"decks,omitempty"`
		CardsToo bool      `json:"cardsToo,omitempty"`
	}

	ParamsCardReviews struct {
		Deck    string `json:"deck,omitempty"`
		StartID int64  `json:"startID"`
	}

	CardReview struct {
		ReviewTime       int64 `json:"reviewTime"`
		CardID           int64 `json:"cardID"`
		USN              int64 `json:"usn"`
		ButtonPressed    int64 `json:"buttonPressed"`
		NewInterval      int64 `json:"newInterval"`
		PreviousInterval int64 `json:"previousInterval"`
		NewFactor        int64 `json:"newFactor"`
		ReviewDuration   int64 `json:"reviewDuration"`
		ReviewType       int64 `json:"reviewType"`
	}

	ResultCardReviews []*CardReview

	// decksManager implements DecksManager.
	decksManager struct {
		Client *Client
	}
)

func (r *CardReview) UnmarshalJSON(data []byte) error {
	var raw []int64
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 9 {
		return fmt.Errorf(
			"unexpected number of fields in card review: got %d, expected 9", len(raw))
	}
	*r = CardReview{
		ReviewTime:       raw[0],
		CardID:           raw[1],
		USN:              raw[2],
		ButtonPressed:    raw[3],
		NewInterval:      raw[4],
		PreviousInterval: raw[5],
		NewFactor:        raw[6],
		ReviewDuration:   raw[7],
		ReviewType:       raw[8],
	}
	return nil
}

// GetAll retrieves all the decks from Anki.
// The result is a slice of string with the names of the decks.
// The method returns an error if:
//   - the api request to ankiconnect fails.
//   - the api returns a http error.
func (dm *decksManager) GetAll() (*[]string, *errors.RestErr) {
	result, restErr := post[[]string, ParamsDefault](dm.Client, ActionDeckNames, nil)
	if restErr != nil {
		return nil, restErr
	}
	return result, nil
}

// Create creates a new deck in Anki.
// The method returns an error if:
//   - the api request to ankiconnect fails.
//   - the api returns a http error.
func (dm *decksManager) Create(name string) *errors.RestErr {
	params := ParamsCreateDeck{
		Deck: name,
	}
	_, restErr := post[int64](dm.Client, ActionCreateDeck, &params)
	if restErr != nil {
		return restErr
	}
	return nil
}

// Delete deletes a deck from Anki
// The method returns an error if:
//   - the api request to ankiconnect fails.
//   - the api returns a http error.
func (dm *decksManager) Delete(name string) *errors.RestErr {
	params := ParamsDeleteDecks{
		Decks:    &[]string{name},
		CardsToo: true,
	}
	_, restErr := post[string](dm.Client, ActionDeleteDecks, &params)
	if restErr != nil {
		return restErr
	}
	return nil
}

func (dm *decksManager) GetReviewsAfter(name string, startTime time.Time) (
	*ResultCardReviews, *errors.RestErr) {
	params := ParamsCardReviews{
		Deck:    name,
		StartID: startTime.UnixMilli(),
	}
	reviews, restErr := post[ResultCardReviews](dm.Client, ActionCardReviews, &params)
	if restErr != nil {
		return nil, restErr
	}
	return reviews, nil
}
