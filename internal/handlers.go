package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "job-queue/database"
	"log"
	"net/http"
	"strconv"
)

type JobRequest struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

func JobHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		fmt.Printf("Job queued!")

		decoder := json.NewDecoder(r.Body)
		var jb JobRequest
		if err := decoder.Decode(&jb); err != nil {
			http.Error(w, "Broken JSON format", http.StatusBadRequest)
			return
		}
		fmt.Println(jb.Name, jb.Payload)
	}
}

func InsertDbHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var jb JobRequest
		if err := decoder.Decode(&jb); err != nil {
			log.Println("JSON decode error:", err)
			http.Error(w, "Couldn't decode the request", http.StatusBadRequest)
			return
		}

		err := db.InsertDb(jb.Name, jb.Payload)
		if err != nil {
			log.Println("DB insert error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Println("Data inserted into DB!")

	}
}

func StatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		idSTR := r.PathValue("id")
		id, err := strconv.Atoi(idSTR)
		if err != nil {
			http.Error(w, "Id couldn't found", http.StatusBadRequest)
		}
		data, err := db.GetJobDb(id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Job not found", http.StatusNotFound)
				return
			}
			log.Println("GetJobDb err:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

	}
}
