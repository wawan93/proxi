package proxi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/wawan93/proxi"
)

func TestFetch(t *testing.T) {
	var cnt int
	var m sync.Mutex
	var proxiServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		cnt++
		m.Unlock()
	}))
	defer proxiServer.Close()

	p := &p{url: proxiServer.URL}
	c, err := proxi.Client(p)
	if err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(nil)
	defer ts.Close()
	url := ts.URL

	tries := 5

	f := proxi.NewFetcher(c, time.Second, tries)
	_, err = f.Fetch(context.Background(), http.MethodGet, url, nil, nil)
	if err != nil {
		t.Error(err)
	}

	m.Lock()
	total := cnt
	m.Unlock()
	if total == 0 {
		t.Error("proxy was not used")
	}
}
