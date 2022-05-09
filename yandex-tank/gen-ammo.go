package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const maxint = int(^uint(0) >> 1)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("[Host: example.org]")
	fmt.Println("[Connection: close]")
	fmt.Println("[User-Agent: Tank]")
	fmt.Println("[Content-type: application/json]")

	for i := 0; i < maxint; i += 1 {
		user, _ := uuid.NewUUID()
		value := fmt.Sprintf("%v.%v", rand.Intn(10), rand.Intn(10))

		fmt.Printf("%v /results/%v\n", len(value), user.String())
		fmt.Println(value)
	}
}
