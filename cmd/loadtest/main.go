package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	endpoint := flag.String("endpoint", "http://localhost:8080", "Base API endpoint")
	concurrency := flag.Int("concurrency", 10, "Concurrent workers")
	iterations := flag.Int("iterations", 50, "Requests per worker")
	flag.Parse()

	start := time.Now()
	var wg sync.WaitGroup
	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for j := 0; j < *iterations; j++ {
				url := fmt.Sprintf("%s/api/v1/products?limit=10", *endpoint)
				resp, err := client.Get(url)
				if err == nil {
					_ = resp.Body.Close()
				}
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Completed %d requests in %s\n", (*concurrency)*(*iterations), elapsed.String())
}
