package ankiconnect

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/privatesquare/bkst-go-utils/utils/errors"
)

const (
	ActionFindCards         = "findCards"
	ActionCardsInfo         = "cardsInfo"
	ActionGetReviewsOfCards = "getReviewsOfCards"
)

type (
	// Notes manager describes the interface that can be used to perform operation on the notes in a deck.
	CardsManager interface {
		Search(query string) (*[]int64, *errors.RestErr)
		Get(query string) (*[]ResultCardsInfo, *errors.RestErr)
		// Note that the current released anki-connect api does not yet support this
		// api call https://github.com/FooSoft/anki-connect/issues/378
		GetReviews(query string) (*ResultGetReviews, *errors.RestErr)
	}

	// notesManager implements NotesManager.
	cardsManager struct {
		Client *Client
	}

	ParamsFindCards struct {
		Query string `json:"query,omitempty"`
	}

	ResultCardsInfo struct {
		Answer     string               `json:"answer,omitempty"`
		Question   string               `json:"question,omitempty"`
		DeckName   string               `json:"deckName,omitempty"`
		ModelName  string               `json:"modelName,omitempty"`
		FieldOrder int64                `json:"fieldOrder,omitempty"`
		Fields     map[string]FieldData `json:"fields,omitempty"`
		Css        string               `json:"css,omitempty"`
		CardId     int64                `json:"cardId,omitempty"`
		Interval   int64                `json:"interval,omitempty"`
		Note       int64                `json:"note,omitempty"`
		Ord        int64                `json:"ord,omitempty"`
		Type       int64                `json:"type,omitempty"`
		Queue      int64                `json:"queue,omitempty"`
		Due        int64                `json:"due,omitempty"`
		Reps       int64                `json:"reps,omitempty"`
		Lapses     int64                `json:"lapses,omitempty"`
		Left       int64                `json:"left,omitempty"`
		Mod        int64                `json:"mod,omitempty"`
	}

	// ParamsCardsInfo represents the ankiconnect API params for getting card info.
	ParamsCardsInfo struct {
		Cards *[]int64 `json:"cards,omitempty"`
	}

	ResultGetReviews map[int64][]ReviewData

	ReviewData struct {
		ID       int64 `json:"id"`
		USN      int64 `json:"usn"`
		Ease     int64 `json:"ease"`
		Interval int64 `json:"ivl"`
		LastIvl  int64 `json:"lastIvl"`
		Factor   int64 `json:"factor"`
		Time     int64 `json:"time"`
		Type     int64 `json:"type"`
	}
)

func (m *ResultGetReviews) UnmarshalJSON(data []byte) error {
	tmp := make(map[string][]ReviewData)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*m = ResultGetReviews{}
	for k, v := range tmp {
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid key %q: %s", k, err)
		}
		(*m)[i] = v
	}
	return nil
}

func (cm *cardsManager) Search(query string) (*[]int64, *errors.RestErr) {
	findParams := ParamsFindCards{
		Query: query,
	}
	return post[[]int64](cm.Client, ActionFindCards, &findParams)
}

func (cm *cardsManager) Get(query string) (*[]ResultCardsInfo, *errors.RestErr) {
	cardIds, restErr := cm.Search(query)
	if restErr != nil {
		return nil, restErr
	}
	infoParams := ParamsCardsInfo{
		Cards: cardIds,
	}
	return post[[]ResultCardsInfo](cm.Client, ActionCardsInfo, &infoParams)
}

func (cm *cardsManager) GetReviews(query string) (
	*ResultGetReviews, *errors.RestErr) {
	cardIds, restErr := cm.Search(query)
	if restErr != nil {
		return nil, restErr
	}
	getReviewsParams := ParamsCardsInfo{
		Cards: cardIds,
	}
	return post[ResultGetReviews](
		cm.Client, ActionGetReviewsOfCards, &getReviewsParams)
}
