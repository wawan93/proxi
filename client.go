package proxi

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type ProxyPool interface {
	Random() string
}

func Client(pool ProxyPool) (*http.Client, error) {
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				proxy := pool.Random()
				u, err := url.Parse(proxy)
				if err != nil {
					return nil, err
				}
				return u, nil
			},
		},
	}

	return c, nil
}

func Fetch(ctx context.Context, client *http.Client, url string, headers http.Header) (data string, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %v", err)
	}
	req.WithContext(ctx)
	delete(headers, "Host")
	req.Header = headers

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c := make(chan string, 5)

	for i := 0; i < 5; i++ {
		go request(client, req, c)
	}

	select {
	case data = <-c:
		return data, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("request timeout")
	}
}
func request(client *http.Client, req *http.Request, c chan string) {
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	c <- string(b)
}
