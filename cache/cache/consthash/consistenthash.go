package consthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(b []byte) uint32

type HashMap struct {
	replicas int
	keys     []int
	maps     map[int]string
	hash     Hash
}

func NewHashMap(replicas int, hash Hash) *HashMap {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &HashMap{
		replicas: replicas,
		keys:     make([]int, 0),
		maps:     make(map[int]string),
		hash:     hash,
	}
}

func (m *HashMap) Add(names ...string) {
	for _, name := range names {
		for i := range m.replicas {
			key := strconv.Itoa(i) + name
			id := int(m.hash([]byte(key)))
			m.keys = append(m.keys, id)
			m.maps[id] = name
		}
	}
	sort.Ints(m.keys)
}

func (m *HashMap) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hashKey := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hashKey
	})
	target := m.maps[m.keys[idx%len(m.keys)]]
	return target
}
