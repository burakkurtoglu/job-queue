package main

import (
	"job-queue/database"
	"job-queue/internal"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	//http.HandleFunc("/jobs", handler.JobHandler())
	godotenv.Load()
	err := database.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/jobs", internal.InsertDbHandler())

	http.HandleFunc("/jobs/{id}", internal.StatusHandler())

	internal.StartWorkerPool()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
