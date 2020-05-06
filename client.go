package proxi

import (
	"net/http"
	"net/url"
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
