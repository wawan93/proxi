/*
ProxyPool implementation for proxy-list.download proxy provider
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

type ProxyListDownloadConfig struct {
	Type    string
	Anon    string
	Country string
}

type ProxyListDownloadPool struct {
	proxies []string
	m       sync.Mutex
	cfg     *ProxyListDownloadConfig
}

func NewProxyListDownloadPool(cfg *ProxyListDownloadConfig) *ProxyListDownloadPool {
	return &ProxyListDownloadPool{
		cfg: cfg,
	}
}

func (p *ProxyListDownloadPool) Update() error {
	query := &url.Values{}
	query.Add("type", p.cfg.Type)
	if p.cfg.Anon != "" {
		query.Add("anon", p.cfg.Anon)
	}
	if p.cfg.Country != "" {
		query.Add("country", p.cfg.Country)
	}

	resp, err := http.Get("https://www.proxy-list.download/api/v1/get?" + query.Encode())
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

func (p *ProxyListDownloadPool) Random() string {
	rand.Seed(time.Now().UnixNano())
	p.m.Lock()
	defer p.m.Unlock()
	return p.proxies[rand.Intn(len(p.proxies))]
}
