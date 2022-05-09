package ldtest

import (
	"fmt"
	"math/rand"
)

type ActionType int

const (
	ActionUpdate ActionType = iota
	ActionInsert
)

type Request struct {
	User   string
	Action ActionType
	Value  float64
}

func MakeRequest(lastNewUser int) (Request, int) {
	req := Request{
		User:   fmt.Sprintf("%v", rand.Intn(lastNewUser)),
		Action: ActionType(rand.Intn(int(ActionInsert))),
		Value:  rand.Float64(),
	}

	if req.Action == ActionInsert {
		req.User = fmt.Sprintf("%v", lastNewUser)
		lastNewUser = lastNewUser + 1
	}

	return req, lastNewUser
}
