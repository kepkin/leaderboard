package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"

	lb "gihtub.com/kepkin/leaderboard"
	"gihtub.com/kepkin/leaderboard/api"
)

func main() {
	var args struct {
		Port int `default:"8080"`
	}
	arg.MustParse(&args)

	store := lb.NewBtreeStore[api.Decimal, api.UserID](
		api.DecimalLess,
		api.DecimalEquals,
		lb.StdLess[api.UserID],
		lb.StdEquals[api.UserID],
	)

	r := gin.New()
	r.Use(gin.Logger())
	apiImpl := api.NewLeaderBoardService(store)
	api.RegisterRoutes(r, apiImpl)

	r.Run(fmt.Sprintf(":%v", args.Port))
}
