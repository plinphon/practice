package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

func fetch(httpe string) string {
	resp, err := http.Get(httpe)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	log.Printf(sb)

	return sb
}

func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- string) {
	defer wg.Done()
	for url := range jobs {
		results <- fetch(url)
	}
}

func main() {

	urls := []string{
		"https://blog.logrocket.com/making-http-requests-in-go/",
		"https://golang.org",
	}

	jobs := make(chan string, 10)
	results := make(chan string, 10)

	var wg sync.WaitGroup

	for w := 0; w < 2; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		log.Printf("Result length: %d\n", len(res))
	}

}
