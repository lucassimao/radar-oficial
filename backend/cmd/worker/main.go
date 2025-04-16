package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	_ = godotenv.Load()

	log.Println("Starting Diário fetcher worker...")

	err := FetchAndUploadDiarios()
	if err != nil {
		log.Fatalf("Worker failed: %v", err)
	}

	log.Println("Worker completed successfully.")
}

func FetchAndUploadDiarios() error {
	// Example only – you’ll implement real logic here
	today := time.Now().Format("2006-01-02")
	fmt.Printf("Downloading Diário Oficial for %s...\n", today)

	// TODO: Download file
	// TODO: Parse / process content
	// TODO: Upload to DO Spaces

	return nil
}

func uploadToSpaces() error {
	endpoint := os.Getenv("DO_SPACES_ENDPOINT")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return err
	}

	filePath := "diario_2025_04_15.pdf"
	bucket := "diarios-oficiais"

	_, err = client.FPutObject(context.Background(), bucket, filePath, "/tmp/"+filePath, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	return err
}
