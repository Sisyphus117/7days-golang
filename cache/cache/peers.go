package cache

import "cache/cache/cachepb"

type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

type PeerGetter interface {
	Get(in *cachepb.Request, out *cachepb.Response) error
}
