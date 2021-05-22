package main

import (
	"fmt"
	"time"
)

func main() {
	hello(5*time.Second, "Hello Go")
}

func hello(d time.Duration, msg string) {
	select {
	case <-time.After(d):
		fmt.Println(msg)
	}
}

/*
func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	hello(ctx, 5*time.Second, "Hello Go")
}

func hello(ctx context.Context, d time.Duration, msg string) {
	select {
	case <-time.After(d):
		fmt.Println(msg)
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}
*/
