package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Go-Lang/go-PostgreSQL/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on the port 8080....")
	log.Fatal(http.ListenAndServe(":8000", r))

}
