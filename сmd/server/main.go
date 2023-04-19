package main

import (
	"fmt"
	"halyklife/internal/db"
	"halyklife/internal/handler"
	"log"
	"net/http"
)

func main() {
	db := db.Init("mongodb://localhost:27017", "mydb")

	http.HandleFunc("/", handler.HandleRequest(db))
	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
