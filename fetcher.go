package proxi

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type fetcher struct {
	c       *http.Client
	timeout time.Duration
	tries   int
}

func NewFetcher(client *http.Client, timeout time.Duration, tries int) *fetcher {
	return &fetcher{client, timeout, tries}
}

func (f *fetcher) Fetch(ctx context.Context, url string, headers http.Header) (data string, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %v", err)
	}
	req.WithContext(ctx)
	delete(headers, "Host")
	req.Header = headers

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := make(chan string)
	for i := 0; i < f.tries; i++ {
		go f.request(req, c)
	}

	select {
	case data = <-c:
		return data, nil
	case <-time.After(f.timeout):
		return "", fmt.Errorf("request timeout")
	}
}

func (f *fetcher) request(req *http.Request, c chan string) {
	resp, err := f.c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	select {
	case <-req.Context().Done():
	case c <- string(b):
	}
}
