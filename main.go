package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HTMLHandler)
	http.HandleFunc("/json", JSONHandler)

	fmt.Println("Server started at localhost:8082")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		fmt.Println("Server error: ", err)
	}
}
