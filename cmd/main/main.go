package main

import (
	"log"
	"net/http"

	"github.com/zowber/zowber-linkz-go/internal/routes"
)

func main() {
	log.Println("Coming up on port 9000!")
	router := routes.NewRouter()
	err := http.ListenAndServe(":9000", router)
	if err != nil {
		log.Fatal(err)
	}
}
