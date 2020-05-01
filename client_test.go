package proxi_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/wawan93/proxi"
)

type p struct {
	url string
}

func (p *p) Random() string {
	return p.url
}

func TestGetClient(t *testing.T) {
	var cnt int
	var proxiServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cnt++
	}))

	p := &p{url: proxiServer.URL}
	c, err := proxi.Client(p)
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	c.Get(ts.URL)

	if cnt == 0 {
		t.Error("proxi was not used")
	}
}

func TestFetch(t *testing.T) {
	var cnt int
	var m sync.Mutex
	var proxiServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		cnt++
		m.Unlock()
	}))

	p := &p{url: proxiServer.URL}
	c, err := proxi.Client(p)
	if err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	url := ts.URL

	c, err = proxi.Client(p)
	if err != nil {
		t.Error(err)
	}

	_, err = proxi.Fetch(context.Background(), c, url, nil)
	if err != nil {
		t.Error(err)
	}
}
