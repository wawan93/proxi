package proxi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/wawan93/proxi"
)

var cnt int
var proxiServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cnt++
}))

type pool struct {
	c int
	m sync.Mutex
}

func (p *pool) Random() string {
	p.m.Lock()
	defer p.m.Unlock()
	p.c++

	return proxiServer.URL
}

func TestGetClient(t *testing.T) {
	p := &pool{}
	c, err := proxi.Client(p)
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	c.Get(ts.URL)

	if p.c == 0 {
		t.Error("proxi was not used")
	}
	if cnt == 0 {
		t.Error("proxi was not used")
	}
}
