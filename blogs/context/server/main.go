//Server Main
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("Handler Started")
	defer fmt.Println("Handler Ended")

	select {
	case <-time.After(5 * time.Second):
		fmt.Fprint(w, "Hello Go")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

}
