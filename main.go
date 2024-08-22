package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Function to stream files from a directory
func streamFiles(dirPath string, fileCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCh <- path
		}
		return nil
	})

	if err != nil {
		log.Printf("Error retrieving files from directory %s: %v", dirPath, err)
	}

	close(fileCh) // Close the channel after processing
}

// Function to upload a file to S3
func uploadFile(svc *s3.S3, filePath string, bucketName string, keyPrefix string, baseDir string, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		errCh <- fmt.Errorf("failed to open file %s: %v", filePath, err)
		return
	}
	defer file.Close()

	// Calculate the relative path to preserve folder structure
	relativePath, err := filepath.Rel(baseDir, filePath)
	if err != nil {
		errCh <- fmt.Errorf("failed to calculate relative path for file %s: %v", filePath, err)
		return
	}

	// Generate the S3 key with the relative path
	key := filepath.Join(keyPrefix, relativePath)
	key = strings.ReplaceAll(key, "\\", "/") // Ensure S3 keys use forward slashes

	// Determine Content-Type
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream" // Default value
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		errCh <- fmt.Errorf("failed to upload file %s: %v", filePath, err)
		return
	}

	fmt.Printf("File successfully uploaded: %s\n", key)
}

func main() {
	// Command-line arguments
	accessKeyId := flag.String("accessKeyId", "", "AWS Access Key ID")
	secretAccessKey := flag.String("secretAccessKey", "", "AWS Secret Access Key")
	region := flag.String("region", "auto", "AWS Region")
	endpoint := flag.String("endpoint", "https://fly.storage.tigris.dev", "AWS Endpoint")
	bucket := flag.String("bucket", "", "S3 Bucket Name")
	folder := flag.String("folder", "", "Path to the folder to upload")
	prefix := flag.String("prefix", "", "Optional prefix for keys in S3 (Subfolder)")
	workers := flag.Int("workers", 10, "Number of concurrent workers")

	flag.Parse()

	if *accessKeyId == "" || *secretAccessKey == "" || *bucket == "" || *folder == "" {
		log.Fatal("All parameters --accessKeyId, --secretAccessKey, --bucket, and --folder are required")
	}

	// AWS S3 configuration
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(*region),
		Credentials: credentials.NewStaticCredentials(*accessKeyId, *secretAccessKey, ""),
		Endpoint:    aws.String(*endpoint),
	})

	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	svc := s3.New(sess)

	// Channels for streaming files and error handling
	fileCh := make(chan string)
	errCh := make(chan error, *workers)

	var wg sync.WaitGroup

	// Start goroutine to retrieve files
	wg.Add(1)
	go streamFiles(*folder, fileCh, &wg)

	// Parallel file uploading
	var uploadWg sync.WaitGroup
	sem := make(chan struct{}, *workers)

	for file := range fileCh {
		uploadWg.Add(1)
		sem <- struct{}{} // Block semaphore while workers are available
		go func(file string) {
			defer func() { <-sem }() // Release semaphore after completion
			uploadFile(svc, file, *bucket, *prefix, *folder, &uploadWg, errCh)
		}(file)
	}

	// Close the error channel when all files are processed
	go func() {
		uploadWg.Wait()
		close(errCh)
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Print all errors
	for err := range errCh {
		log.Println("Error:", err)
	}

	fmt.Println("Upload completed.")
}

/*

AWS_ACCESS_KEY_ID=tid_zXeLSU_wRFnOxeAWYlBGsiwaKcrHKzNdGrLGVAAcUXXDcQtJtK
AWS_SECRET_ACCESS_KEY=tsec_EYz6vk_9z4cf9rQFxADYyBVPoHfWbgUlRQNxzs_PqMMFBZ-SgCGQGNUTT+bQQ+WGRGnR2Z
AWS_ENDPOINT_URL_S3=https://fly.storage.tigris.dev
AWS_REGION=auto

*/
