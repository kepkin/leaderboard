package main

import (
	"fmt"
	// "math"
	"os"
	// "math/rand"
	"log"
	"strings"
	// "encoding/gob"

	"github.com/alexflint/go-arg"
	"gonum.org/v1/gonum/stat/distuv"

	lbtest "gihtub.com/kepkin/leaderboard/test"
)

const InitialDataN = 10_000_000
const RequestsN = 10_000_000
const UpperBoundRand = 20

type cmdInit struct {
	Path string
}

type cmdAmmo struct {
	Path string
}

type cmdArgs struct {
	Init       *cmdInit       `arg:"subcommand"`
	Ammo       *cmdAmmo       `arg:"subcommand"`
	ServerInit *cmdServerInit `arg:"subcommand"`
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

func getInitial(args *cmdInit) {
	file, err := os.Create(args.Path) // For read access.
	if err != nil {
		log.Fatal(err)
		return
	}

	variants := []distuv.LogNormal{
		{Mu: 0, Sigma: 1},
		{Mu: 0, Sigma: 0.25},
		{Mu: 0, Sigma: 0.5},
		{Mu: 1, Sigma: 1},
		{Mu: 1, Sigma: 0.25},
		{Mu: 1, Sigma: 0.5},
	}

	headers := make([]string, 0)
	headers = append(headers, "user")

	for _, v := range variants {
		headers = append(headers, fmt.Sprintf("m%vs%v", v.Mu, v.Sigma))
	}
	fmt.Fprintln(file, strings.Join(headers, ","))

	for i := 0; i < InitialDataN; i++ {
		values := make([]string, len(variants)+1)
		values[0] = fmt.Sprintf("\"%v\"", i)
		for j, v := range variants {
			vr := v.Rand()
			for vr > UpperBoundRand {
				vr = v.Rand()
			}

			values[j+1] = fmt.Sprintf("%v", vr)
		}
		fmt.Fprintln(file, strings.Join(values, ","))
	}
}

func main() {
	var args cmdArgs
	arg.MustParse(&args)

	switch {
	case args.Ammo != nil:
		genAmmoData(args.Ammo)
	case args.Init != nil:
		getInitial(args.Init)
	case args.ServerInit != nil:
		serverInit(args.ServerInit)
	}
}
