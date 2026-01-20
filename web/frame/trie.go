package frame

import (
	"fmt"
	"strings"
)

// trie node
type Node struct {
	children map[string]*Node
	//param for path like a/:name
	param     *Node
	paramName string
	//wildcard for path like a/b/*filepath or a/b/*
	iswild   bool
	wildName string
}

// a router maintains a single trie
type Trie struct {
	root *Node
	keys map[string]struct{}
}

func NewTrie() *Trie {
	return &Trie{root: NewNode(), keys: make(map[string]struct{})}
}

func NewNode() *Node {
	return &Node{make(map[string]*Node), nil, "", false, ""}
}

// create a trie node from a path key
func (t *Trie) addKey(key string) (string, error) {
	cur := t.root
	if key == "/" {
		t.keys[key] = struct{}{}
		return "/", nil
	}
	parts := strings.Split(key[1:], "/")
	added := false
	for i, part := range parts {
		if strings.HasPrefix(part, "*") {
			cur.iswild = true
			cur.wildName = part[1:]
			if len(cur.wildName) == 0 {
				cur.wildName = "subpath"
			}
			if i == 0 {
				key = "/*"
			} else {
				key = "/" + strings.Join(parts[:i], "/") + "/*"
			}
			t.keys[key] = struct{}{}
			return key, nil
		}
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
	}
	if added {
		t.keys[key] = struct{}{}
		return key, nil
	}
	return "", fmt.Errorf("key conflicted")

}

// parse api and params from a path
func (t *Trie) getKey(path string) (string, map[string]string, error) {
	cur := t.root
	parts := strings.Split(path[1:], "/")
	params := make(map[string]string)
	if path == "/" {
		return path, params, nil
	}
	for i, part := range parts {
		if next, has := cur.children[part]; has {
			cur = next
		} else if cur.param != nil {
			params[cur.paramName[1:]] = part
			parts[i] = cur.paramName
			cur = cur.param
			continue
		} else if cur.iswild {
			params[cur.wildName] = strings.Join(parts[i:], "/")
			api := "/" + strings.Join(parts[:i], "/") + "/*"
			return api, params, nil
		} else {
			return "", nil, fmt.Errorf("invalid path")
		}
	}
	api := "/" + strings.Join(parts, "/")
	return api, params, nil
}
