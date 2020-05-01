/*
ProxyPool implementation for freeproxy-list.ru proxy provider
*/
package pool

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FreeproxyListConfig struct {
	Token         string
	Count         int
	Country       string
	Port          []int
	Accessibility int
	Anonymity     bool
}

type FreeproxyListPool struct {
	proxies []string
	m       sync.Mutex
	cfg     *FreeproxyListConfig
}

func NewFreeproxyListPool(cfg *FreeproxyListConfig) *FreeproxyListPool {
	return &FreeproxyListPool{
		cfg: cfg,
	}
}

func (p *FreeproxyListPool) Update() error {
	query := &url.Values{}
	query.Add("token", p.cfg.Token)
	if p.cfg.Count != 0 {
		query.Add("count", strconv.Itoa(p.cfg.Count))
	}
	if p.cfg.Country != "" {
		query.Add("country", p.cfg.Country)
	}
	if p.cfg.Accessibility != 0 {
		query.Add("accessibility", strconv.Itoa(p.cfg.Accessibility))
	}
	if p.cfg.Anonymity {
		query.Add("anonymity", "true")
	}

	resp, err := http.Get("https://www.freeproxy-list.ru/api/proxy?" + query.Encode())
	if err != nil {
		return fmt.Errorf("cannot get proxies: %v", err)
	}
	defer resp.Body.Close()

	proxiesBytes, err := ioutil.ReadAll(resp.Body)
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

func (p *FreeproxyListPool) Random() string {
	rand.Seed(time.Now().UnixNano())
	p.m.Lock()
	defer p.m.Unlock()
	return p.proxies[rand.Intn(len(p.proxies))]
}
