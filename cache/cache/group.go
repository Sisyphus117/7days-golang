package cache

import (
	"cache/cache/cachepb"
	"cache/cache/singleflight"
	"sync"
)

type Group struct {
	name   string
	cache  *cache
	getter Getter
	peers  PeerPicker
	single *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetFunc func(key string) ([]byte, error)

func (f GetFunc) Get(key string) ([]byte, error) {
	return f(key)
}

func NewGroup(name string, bytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		cache: &cache{
			bytes: bytes,
		},
		getter: getter,
		single: singleflight.NewGroup(),
	}

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	if g, has := groups[name]; has {
		return g
	}
	return nil
}

func (g *Group) Get(key string) (ByteView, error) {
	mu.Lock()
	defer mu.Unlock()

	value, err := g.load(key)
	if err == nil {
		return value, nil
	}
	v, err := g.getLocal(key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: v}, nil
}

func (g *Group) load(key string) (ByteView, error) {
	fn := func() (any, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				value, err := g.getFromPeer(peer, key)
				if err == nil {
					return value, nil
				}
			}
		}
		return g.cache.get(key)
	}
	bw, err := g.single.Call(key, fn)
	if err != nil {
		return ByteView{}, err
	}
	return bw.(ByteView), nil
}

func (g *Group) getLocal(key string) ([]byte, error) {
	value, err := g.getter.Get(key)
	g.storeCache(key, value)
	return value, err
}

func (g *Group) storeCache(key string, value []byte) {
	g.cache.add(key, NewByteView(value))
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("can't overwrite peers")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &cachepb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &cachepb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return NewByteView(nil), err
	}
	return NewByteView(res.Value), nil
}
