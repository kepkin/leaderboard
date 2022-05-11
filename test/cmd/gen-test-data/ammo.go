package main

import (
	"fmt"
	"os"
	"log"
	"strings"

	"github.com/alexflint/go-arg"
	"gonum.org/v1/gonum/stat/distuv"

	lbtest "gihtub.com/kepkin/leaderboard/test"
)

type cmdAmmo struct {
	Path string
}

func genAmmoData(args *cmdAmmo) {
	file, err := os.Create(args.Path) // For read access.
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Fprintln(file, "[Host: example.org]")
	fmt.Fprintln(file, "[Connection: close]")
	fmt.Fprintln(file, "[User-Agent: Tank]")
	fmt.Fprintln(file, "[Content-type: application/json]")

	lastNewUser := InitialDataN
	for i := 0; i < RequestsN; i++ {
		var req lbtest.Request
		req, lastNewUser = lbtest.MakeRequest(lastNewUser)

		val := fmt.Sprintf("%v", req.Value)
		fmt.Fprintf(file, "%v /results/%v\n", len(val), req.User)
		fmt.Fprintf(file, "%v\n", val)
	}

	file.Close()
}
