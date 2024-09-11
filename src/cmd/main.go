package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
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
	http.ListenAndServe(":8080", nil)
}
