package proxi

import (
	"net/http"
	"net/url"
)

type ProxyPool interface {
	Random() string
}

func Client(pool ProxyPool) (*http.Client, error) {
	// rand.Seed(time.Now().UnixNano())
	// proxies, err := pool.Update()
	// if err != nil {
	// 	return nil, err
	// }
	// if len(proxies) == 0 {
	// 	return nil, fmt.Errorf("proxy pool is empty")
	// }

	c := &http.Client{
		Transport: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				// proxy := proxies[rand.Intn(len(proxies))]
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
