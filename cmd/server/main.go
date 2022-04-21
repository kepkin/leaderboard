package main

import (
	"fmt"
	
	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"

	"gihtub.com/kepkin/leaderboard"
	"gihtub.com/kepkin/leaderboard/api"
)

func main() {
	var args struct {
		Port int `default:"8080"`
	}
	arg.MustParse(&args)


	store := leaderboard.NewStore[api.Decimal, api.UserID]()

	r := gin.New()
	apiImpl := api.NewLeaderBoardService(store)
	api.RegisterRoutes(r, apiImpl)

	r.Run(fmt.Sprintf(":%v", args.Port))
}
