package lru

import "container/list"

type LRU struct {
	//maxBytes==0 for no limit
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	onEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewLRU(maxBytes int64, onEvicted func(key string, value Value)) *LRU {
	return &LRU{
		maxBytes:  maxBytes,
		nbytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func NewEntry(key string, value Value) *entry {
	return &entry{key, value}
}

func (e *entry) calcMemo() int64 {
	return int64(len(e.key)) + int64(e.value.Len())
}

func (r *LRU) Get(key string) (Value, bool) {
	if elem, has := r.cache[key]; has {
		r.ll.MoveToFront(elem)
		v := elem.Value.(*entry).value
		return v, true
	}
	return nil, false

}

func (r *LRU) RemoveOldest() {
	elem := r.ll.Back()
	if elem == nil {
		return
	}

	entry := elem.Value.(*entry)
	key := entry.key
	value := entry.value
	delete(r.cache, key)
	r.nbytes -= entry.calcMemo()
	r.ll.Remove(elem)
	if r.onEvicted != nil {
		r.onEvicted(key, value)
	}
}

func (r *LRU) Add(key string, value Value) {
	if elem, has := r.cache[key]; has {
		entry := elem.Value.(*entry)
		r.nbytes -= entry.calcMemo()
		entry.value = value
		r.nbytes += entry.calcMemo()
		r.ll.MoveToFront(elem)
		return
	}
	entry := NewEntry(key, value)
	r.cache[key] = r.ll.PushFront(entry)
	r.nbytes += entry.calcMemo()
	if r.maxBytes > 0 {
		for r.nbytes > r.maxBytes {
			r.RemoveOldest()
		}
	}
}

func (r *LRU) Len() int {
	return r.ll.Len()
}
