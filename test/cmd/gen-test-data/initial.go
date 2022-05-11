package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gonum.org/v1/gonum/stat/distuv"
)

type cmdInit struct {
	Path string
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
