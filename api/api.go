package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"gihtub.com/kepkin/leaderboard"
)

//go:generate go run ./gen/main.go


func (p *Points) UnmarshalJSON(decimalBytes []byte) error {
	return (*decimal.Decimal)(p).UnmarshalJSON(decimalBytes)
}

var _ LeaderBoardService = (*LeaderBoardServiceImpl)(nil)

type Decimal decimal.Decimal

func (d Decimal) Less(other leaderboard.Ordered) bool {
	return decimal.Decimal(d).LessThan(decimal.Decimal(other.(Decimal)))
}

func (d Decimal) Equals(other leaderboard.Comparable) bool {
	return decimal.Decimal(d).Equal(decimal.Decimal(other.(Decimal)))
}

type UserID string

func (u UserID) Less(other leaderboard.Ordered) bool {
	return u < other.(UserID)
}

func (u UserID) Equals(other leaderboard.Comparable) bool {
	return u == other.(UserID)
}

type LeaderBoardServiceImpl struct {
	store *leaderboard.Store[Decimal, UserID]
}

func NewLeaderBoardService(store *leaderboard.Store[Decimal, UserID]) LeaderBoardService {
	return &LeaderBoardServiceImpl{
		store: store,
	}
}

// Errors processing

func (p *LeaderBoardServiceImpl) ProcessMakeRequestErrors(c *gin.Context, errors []FieldError) {
	c.JSON(http.StatusBadRequest, fmt.Sprintf("parse request error: %+v", errors))
}

func (p *LeaderBoardServiceImpl) ProcessValidateErrors(c *gin.Context, errors []FieldError) {
	c.JSON(http.StatusBadRequest, fmt.Sprintf("validate request error: %+v", errors))
}



func (p LeaderBoardServiceImpl) GetLeaderBoard(in GetLeaderBoardRequest, c *gin.Context) {

}

func (p LeaderBoardServiceImpl) PostResults(in PostResultsRequest, c *gin.Context) {
	p.store.Insert(leaderboard.BTreeLeaf[Decimal, UserID]{Decimal(in.Body.JSON), UserID(in.Path.UserID)})
}

