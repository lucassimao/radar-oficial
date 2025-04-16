package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	crawler "radaroficial.app/internal"
	"radaroficial.app/internal/diarios"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		crawler.FetchAndUploadDiarios()

		// today := time.Now().Format("2006-01-02") // "2006-01-02" is the layout string to get YYYY-MM-DD
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		resp, err := diarios.FetchGovernoPiauiDiarios(yesterday, 10)
		if err != nil {
			log.Fatal("Error fetching Diários:", err)
		}

		for i, row := range resp.Data {
			fmt.Printf("[%d] Tipo: %s | Data Publicação: %s | Última Modificação: %s", i+1, row[1], row[2], row[3])
		}

		fmt.Fprint(w, "pong")
	})

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
