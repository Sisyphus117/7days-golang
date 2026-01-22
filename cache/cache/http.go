package cache

import (
	"cache/cache/cachepb"
	"cache/cache/consthash"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultName     = "/default/"
	defaultReplicas = 50
)

type HTTPPool struct {
	baseUrl     string
	name        string
	mu          sync.Mutex
	peers       *consthash.HashMap
	httpGetters map[string]*HTTPGetter
}

type HTTPGetter struct {
	baseUrl string
}

func NewHTTPPool(baseUrl string) *HTTPPool {
	return &HTTPPool{
		baseUrl: baseUrl,
		name:    defaultName,
	}
}

func (p *HTTPPool) Run() {
	log.Fatal(http.ListenAndServe(p.baseUrl, p))
}

func (p *HTTPPool) Log(msg string) {
	log.Printf("[server %s] : %s\n", p.name, msg)
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	base := p.name
	if !strings.HasPrefix(url, base) {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	parts := strings.Split(url[len(base):], "/")
	if len(parts) < 2 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	group, key := parts[0], parts[1]

	c := groups[group]
	if c == nil {
		http.Error(w, "group "+group+"not found", http.StatusNotFound)
		return
	}
	value, err := c.Get(key)
	if err != nil {
		p.Log(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(value.b)

}

func (h *HTTPGetter) Get(in *cachepb.Request, out *cachepb.Response) error {
	u := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(in.Group), url.QueryEscape(in.Key))
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server error, status code: %d", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	out.Value = bytes
	return nil
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consthash.NewHashMap(defaultReplicas, nil)
	p.httpGetters = make(map[string]*HTTPGetter)
	for _, peer := range peers {
		p.httpGetters[peer] = &HTTPGetter{baseUrl: peer + p.baseUrl}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.name {
		return p.httpGetters[peer], true
	}

	return nil, false
}

var _ PeerGetter = (*HTTPGetter)(nil)
