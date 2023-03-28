# Proxi

A library to make requests through free proxies

## Usage

```go
package main

import (
	"log"
	"context"
	"net/http"
	"time"

	"github.com/wawan93/proxi"
	"github.com/wawan93/proxi/pool"
)

func main() {
	cfg := &pool.ProxyListDownloadConfig{
		Type: "https",
	}
	p := pool.NewProxyListDownloadPool(cfg)

	// It's important to call Update before use this pool, because proxies list is empty
	if err := p.Update(); err != nil {
		log.Fatalf("cannot get proxies: %v", err)
	}

	client, err := proxi.Client(p)
	if err != nil {
		log.Fatal(err)
	}

	timeout := 5 * time.Second
	retries := 5
	f := proxi.NewFetcher(client, timeout, retries)
	
	url := "https://example.com"
	data, err := f.Fetch(context.Background(), http.MethodGet, url, nil, nil)
	if err != nil {
        log.Fatal(err)
    }
    log.Println(data)
}
```