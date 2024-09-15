package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"src/database"
	"src/handlers"
	"src/internal"
	"strings"
)

func main() {

	db, err := database.InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}
	defer db.Close()

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "localhost:8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Hello World")
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {

		database.SendAllData(w, r, db)
	})

	http.HandleFunc("/api/bids/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Bad Request: only POST is allowed", http.StatusBadRequest)
			return
		}
		handlers.CreateBids(w, r, db)
	})

	http.HandleFunc("/api/bids/my", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Bad Request: only GET is allowed", http.StatusBadRequest)
			return
		}
		handlers.GetBid(w, r, db)
	})

	http.HandleFunc("/api/bids/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if strings.HasSuffix(r.URL.Path, "/status") {
				handlers.GetBidStatus(w, r, db)
			}
		case http.MethodPut:
			if strings.HasSuffix(r.URL.Path, "/status") {
				handlers.ChangeBidStatus(w, r, db)
			}
		default:
			http.Error(w, "Unsupported method", http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Bad Request: only GET is allowed", http.StatusBadRequest)
			return
		}

		handlers.PingHandler(w, r, db)
	})

	http.HandleFunc("/api/tenders/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Bad Request: only POST is allowed", http.StatusBadRequest)
			return
		}
		handlers.CreateTender(w, r, db)
	})

	http.HandleFunc("/api/tenders/my", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Bad Request: only GET is allowed", http.StatusBadRequest)
			return
		}
		handlers.GetTender(w, r, db)
	})

	http.HandleFunc("/api/tenders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Bad Request: only GET is allowed", http.StatusBadRequest)
		}
		handlers.GetTender(w, r, db)

	})

	http.HandleFunc("/api/tenders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if strings.HasSuffix(r.URL.Path, "/status") {
				handlers.GetBannerStatus(w, r, db)
			}
		case http.MethodPut:
			if strings.HasSuffix(r.URL.Path, "/status") {
				handlers.ChangeTenderStatus(w, r, db)
			} else {
				version, isDigit := internal.IntIsLast(r.URL.Path)
				if isDigit {
					if version < 1 {
						w.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(w).Encode("invalid version")
						return
					}
					handlers.TenderRollBack(w, r, db, version)
				}
			}

		case http.MethodPatch:
			if strings.HasSuffix(r.URL.Path, "/edit") {
				handlers.ChangeTender(w, r, db)
			}
		default:
			http.Error(w, "Unsupported method", http.StatusBadRequest)
			return
		}
	})
	fmt.Printf("Starting server on %s\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
