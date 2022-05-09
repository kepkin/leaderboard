package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	lbtest "gihtub.com/kepkin/leaderboard/test"
)

type cmdServerInit struct {
	Path     string
	Endpoint string
	InitTime time.Duration
}

func initServer(args cmdServerInit) {
	testData, err := lbtest.NewTestData()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if args.InitTime > 0 {
		var ctxCancel context.CancelFunc
		ctx, ctxCancel = context.WithTimeout(context.Background(), args.InitTime)
		defer ctxCancel()
	}

	_, err = testData.Initialize(ctx, 1, func(score float64, user string) {
		req, err := makeRequest(args.Endpoint, user, fmt.Sprint(score))
		if err != nil {
			log.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatal(err)
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}

func makeRequest(endpoint, user, value string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodPost, endpoint+"/"+user, bytes.NewBufferString(value))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	return r, nil
}
