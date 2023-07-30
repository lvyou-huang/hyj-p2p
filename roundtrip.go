package server

import (
	"fmt"
	"net/http"
	"time"
)

type roundTripper struct {
	rt http.RoundTripper
}

func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := rt.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	fmt.Printf("Request took %s\n", elapsed)
	fmt.Printf("Response size: %d bytes\n", resp.ContentLength)
	bandwidth := float64(resp.ContentLength) / elapsed.Seconds()
	fmt.Printf("Bandwidth: %.2f bytes/s\n", bandwidth)
	return resp, nil
}

/*
func main() {
	client := &http.Client{
		Transport: roundTripper{http.DefaultTransport},
	}
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	client.Do(req)
}*/
