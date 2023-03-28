/*
ProxyPool implementation for sslproxies.org proxy provider
*/
package pool

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type SSLProxiesConfig struct {
	Email       string
	Pass        string
	Pid         string
	ShowCountry bool
	Https       bool
}

type SSLProxiesPool struct {
	proxies []string
	m       sync.Mutex
	cfg     *SSLProxiesConfig
}

func NewSSLProxiesPool(cfg *SSLProxiesConfig) *SSLProxiesPool {
	return &SSLProxiesPool{
		cfg: cfg,
	}
}

func (p *SSLProxiesPool) Update() error {
	query := &url.Values{}
	query.Add("email", p.cfg.Email)
	query.Add("pass", p.cfg.Pass)
	query.Add("pid", p.cfg.Pid)
	showcountry := "no"
	if p.cfg.ShowCountry {
		showcountry = "yes"
	}
	query.Add("showcountry", showcountry)
	https := "no"
	if p.cfg.Https {
		https = "yes"
	}
	query.Add("https", https)

	resp, err := http.Get("http://list.didsoft.com/get?" + query.Encode())
	if err != nil {
		return fmt.Errorf("cannot get proxies: %v", err)
	}
	defer resp.Body.Close()

	proxiesBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot parse proxies: %v", err)
	}
	proxies := strings.Split(string(proxiesBytes), "\n")

	result := make([]string, 0)
	for _, proxy := range proxies {
		result = append(result, "http://"+strings.TrimSpace(proxy))
	}
	log.Printf("fetched %d proxies", len(result))

	if len(result) == 0 {
		return fmt.Errorf("proxy pool is empty")
	}

	p.m.Lock()
	p.proxies = result
	p.m.Unlock()
	return nil
}

func (p *SSLProxiesPool) Random() string {
	rand.Seed(time.Now().UnixNano())
	p.m.Lock()
	defer p.m.Unlock()
	return p.proxies[rand.Intn(len(p.proxies))]
}
