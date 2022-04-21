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

	for i := 0; i < maxint; i += 1 {
		user, _ := uuid.NewUUID()
		value := fmt.Sprintf("%v.%v", rand.Intn(999999), rand.Intn(999999))

		fmt.Printf("%v /results/%v\n", len(value), user.String())
		fmt.Println(value)
	}
}
