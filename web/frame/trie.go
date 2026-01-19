package frame

import (
	"fmt"
	"strings"
)

type Node struct {
	children  map[string]*Node
	param     *Node
	paramName string
}
type Trie struct {
	root *Node
	keys map[string]struct{}
}

func NewTrie() *Trie {
	return &Trie{root: NewNode(), keys: make(map[string]struct{})}
}

func NewNode() *Node {
	return &Node{make(map[string]*Node), nil, ""}
}

func (t *Trie) addKey(key string) error {
	cur := t.root
	parts := strings.Split(key[1:], "/")
	added := false
	if len(parts) == 0 {
		return nil
	}
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			if cur.param == nil {
				param := NewNode()
				cur.paramName = part
				cur.param = param
				added = true
			}
			cur = cur.param
			continue
		}
		if next, has := cur.children[part]; !has {
			cur.children[part] = NewNode()
			cur = cur.children[part]
			added = true
		} else {
			cur = next
		}
		if added {
			t.keys[key] = struct{}{}
		} else {
			return fmt.Errorf("key conflicted")
		}
	}
	return nil
}

func (t *Trie) getKey(path string) (string, []string, error) {
	cur := t.root
	parts := strings.Split(path[1:], "/")
	params := make([]string, 0)
	for i, part := range parts {
		if next, has := cur.children[part]; has {
			cur = next
		} else if cur.param != nil {
			params = append(params, parts[i])
			parts[i] = cur.paramName
			next = cur.param
			continue
		} else {
			return "", nil, fmt.Errorf("invalid path")
		}
	}
	api := "/" + strings.Join(parts, "/")
	return api, params, nil
}
