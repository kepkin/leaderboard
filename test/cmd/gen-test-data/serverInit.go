package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"io"
	"io/ioutil"
	"time"

	lbtest "gihtub.com/kepkin/leaderboard/test"
)

type cmdServerInit struct {
	Path     string
	Endpoint string
	InitTime time.Duration
}

func httpClient() *http.Client {
    client := &http.Client{
        Transport: &http.Transport{
            MaxIdleConnsPerHost: 20,
        },
        Timeout: 10 * time.Second,
    }

    return client
}

func serverInit(args *cmdServerInit) {
	c := httpClient()
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
		res, err := c.Do(req)
		if err != nil {
			log.Fatal(fmt.Errorf("couldn't make request: %w", err))
		}
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Fatal(res.StatusCode)
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
