package proxi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
