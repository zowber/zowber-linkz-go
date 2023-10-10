package main

import (
	"log"
	"net/http"

	"github.com/zowber/zowber-linkz-go/internal/routes"
)

func main() {
	log.Println("Sup")

	router := routes.NewRouter()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
