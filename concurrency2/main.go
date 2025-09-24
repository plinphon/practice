package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"
)

func toHash(filePaths []string) []string {
	var result []string

	for _, path := range filePaths {
		data := []byte(path)
		hash := sha256.Sum256(data)
		result = append(result, hex.EncodeToString(hash[:]))
	}

	return result
}

func worker(ctx context.Context, jobs <-chan []string, results chan<- []string) {
	for {
		select {
		case <-ctx.Done():
			return
		case filepath, ok := <-jobs:
			if !ok {
				return
			}
			results <- toHash(filepath)
		}
	}

}

func main() {

	groups := [][]string{
		{"file1.txt", "file2.txt"},
		{"image1.png", "image2.png", "image3.png"},
		{"doc1.pdf"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jobs := make(chan []string, 10)
	results := make(chan []string, 10)

	maxWorkers := 5
	for i := 0; i < maxWorkers; i++ {
		go worker(ctx, jobs, results)
	}

	for _, group := range groups {
		jobs <- group
	}
	close(jobs)

	for i := 0; i < len(groups); i++ {
		log.Println(<-results)
	}
}
