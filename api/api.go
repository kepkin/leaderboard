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

func (p *Points) MarshalJSON() ([]byte, error) {
	return (*decimal.Decimal)(p).MarshalJSON()
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


func GetAdjacentLeaders[K leaderboard.Ordered, V leaderboard.Comparable](itr *leaderboard.Iter[K,V], visiter func(leaderboard.BTreeLeaf[K, V]), before int, after int) {
	backItr := itr
	for ; before > 0 && backItr != nil; before -= 1 {
		backItr = backItr.Prev()
	}

	for i := 0; i < after && itr != nil; i+=1 {
		visiter(itr.Value())
		itr = itr.Next()
	}
}

func (p LeaderBoardServiceImpl) PostResults(in PostResultsRequest, c *gin.Context) {
	itr, err := p.store.Insert(leaderboard.BTreeLeaf[Decimal, UserID]{Decimal(in.Body.JSON), UserID(in.Path.UserID)})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	data := make(UserPointsArr, 0, 21)
	GetAdjacentLeaders(itr, func(v leaderboard.BTreeLeaf[Decimal, UserID]) {
		data = append(data, UserPoints{Points: Points(v.OrderKey), UserID: string(v.Value)})
	}, 10, 10)

	c.JSON(http.StatusOK, data)
}
