package main

import (
	"fmt"
	"net/http"
	"template/controller"
)



func main() {
	http.HandleFunc("/", controller.HelloWorldHandler)
	fmt.Println("Server is listening on port 5000...")
	http.ListenAndServe(":5000", nil)
}
