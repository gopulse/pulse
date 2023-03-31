package routing

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type store struct {
	root  *node
	count int // number of nodes in the store
}

type node struct {
	static bool // whether the node is a static node or param node

	key  string      // the key identifying this node
	data interface{} // the data associated with this node. nil if not a data node.

	order    int // the order at which the data was added. used to be pick the first one when matching multiple
	minOrder int // minimum order among all the child nodes and this node

	children      []*node // child static nodes, indexed by the first byte of each child key
	paramChildren []*node // child param nodes

	regex    *regexp.Regexp // regular expression for a param node containing regular expression key
	paramIdx int            // the parameter index, meaningful only for param node
	params   []string       // the parameter names collected from the root till this node
}

// newStore returns a new store
func newStore() *store {
	return &store{
		root: &node{
			static:        true,
			children:      make([]*node, 256),
			paramChildren: make([]*node, 0),
			paramIdx:      -1,
			params:        []string{},
		},
	}
}

func (s *store) Add(key string, data interface{}) int {
	s.count++
	return s.root.add(key, data, s.count)
}

func (s *store) Get(path string, pvalues []string) (data interface{}, pnames []string) {
	data, pnames, _ = s.root.get(path, pvalues)
	return
}

// Add adds a new node to the store with the given key and data
func (n *node) add(key string, data interface{}, order int) int {
	matched := 0

	// find the common prefix
	for ; matched < len(key) && matched < len(n.key); matched++ {
		if key[matched] != n.key[matched] {
			break
		}
	}

	if matched == len(n.key) {
		if matched == len(key) {
			// the node key is the same as the key: make the current node as data node
			// if the node is already a data node, ignore the new data since we only care the first matched node
			if n.data == nil {
				n.data = data
				n.order = order
			}
			return n.paramIdx + 1
		}

		// the node key is a prefix of the key: create a child node
		newKey := key[matched:]

		// try adding to a static child
		if child := n.children[newKey[0]]; child != nil {
			if pn := child.add(newKey, data, order); pn >= 0 {
				return pn
			}
		}
		// try adding to a param child
		for _, child := range n.paramChildren {
			if pn := child.add(newKey, data, order); pn >= 0 {
				return pn
			}
		}

		return n.addChild(newKey, data, order)
	}

	if matched == 0 || !n.static {
		// no common prefix, or partial common prefix with a non-static node: should skip this node
		return -1
	}

	// the node key shares a partial prefix with the key: split the node key
	n1 := &node{
		static:        true,
		key:           n.key[matched:],
		data:          n.data,
		order:         n.order,
		minOrder:      n.minOrder,
		paramChildren: n.paramChildren,
		children:      n.children,
		paramIdx:      n.paramIdx,
		params:        n.params,
	}

	n.key = key[0:matched]
	n.data = nil
	n.paramChildren = make([]*node, 0)
	n.children = make([]*node, 256)
	n.children[n1.key[0]] = n1

	return n.add(key, data, order)
}

// addChild creates static and param nodes to store the given data
func (n *node) addChild(key string, data interface{}, order int) int {
	// find the first occurrence of a param token
	p0, p1 := -1, -1
	for i := 0; i < len(key); i++ {
		if p0 < 0 && key[i] == '<' {
			p0 = i
		}
		if p0 >= 0 && key[i] == '>' {
			p1 = i
			break
		}
	}

	if p0 > 0 && p1 > 0 || p1 < 0 {
		// param token occurs after a static string, or no param token: create a static node
		child := &node{
			static:        true,
			key:           key,
			minOrder:      order,
			children:      make([]*node, 256),
			paramChildren: make([]*node, 0),
			paramIdx:      n.paramIdx,
			params:        n.params,
		}
		n.children[key[0]] = child
		if p1 > 0 {
			// param token occurs after a static string
			child.key = key[:p0]
			n = child
		} else {
			// no param token: done adding the child
			child.data = data
			child.order = order
			return child.paramIdx + 1
		}
	}

	// add param node
	child := &node{
		static:        false,
		key:           key[p0 : p1+1],
		minOrder:      order,
		children:      make([]*node, 256),
		paramChildren: make([]*node, 0),
		paramIdx:      n.paramIdx,
		params:        n.params,
	}
	pattern := ""
	pname := key[p0+1 : p1]
	for i := p0 + 1; i < p1; i++ {
		if key[i] == ':' {
			pname = key[p0+1 : i]
			pattern = key[i+1 : p1]
			break
		}
	}
	if pattern != "" {
		// the param token contains a regular expression
		child.regex = regexp.MustCompile("^" + pattern)
	}
	pnames := make([]string, len(n.params)+1)
	copy(pnames, n.params)
	pnames[len(n.params)] = pname
	child.params = pnames
	child.paramIdx = len(pnames) - 1
	n.paramChildren = append(n.paramChildren, child)

	if p1 == len(key)-1 {
		// the param token is at the end of the key
		child.data = data
		child.order = order
		return child.paramIdx + 1
	}

	// process the rest of the key
	return child.addChild(key[p1+1:], data, order)
}

func (n *node) get(key string, pvalues []string) (data interface{}, pnames []string, order int) {
	order = math.MaxInt32

repeat:
	if n.static {
		// check if the node key is a prefix of the given key
		// a slightly optimized version of strings.HasPrefix
		nkl := len(n.key)
		if nkl > len(key) {
			return
		}
		for i := nkl - 1; i >= 0; i-- {
			if n.key[i] != key[i] {
				return
			}
		}
		key = key[nkl:]
	} else if n.regex != nil {
		// param node with regular expression
		if n.regex.String() == "^.*" {
			pvalues[n.paramIdx] = key
			key = ""
		} else if match := n.regex.FindStringIndex(key); match != nil {
			pvalues[n.paramIdx] = key[0:match[1]]
			key = key[match[1]:]
		} else {
			return
		}
	} else {
		// param node matching non-"/" characters
		i, kl := 0, len(key)
		for ; i < kl; i++ {
			if key[i] == '/' {
				pvalues[n.paramIdx] = key[0:i]
				key = key[i:]
				break
			}
		}
		if i == kl {
			pvalues[n.paramIdx] = key
			key = ""
		}
	}

	if len(key) > 0 {
		// find a static child that can match the rest of the key
		if child := n.children[key[0]]; child != nil {
			if len(n.paramChildren) == 0 {
				// use goto to avoid recursion when no param children
				n = child
				goto repeat
			}
			data, pnames, order = child.get(key, pvalues)
		}
	} else if n.data != nil {
		// do not return yet: a param node may match an empty string with smaller order
		data, pnames, order = n.data, n.params, n.order
	}

	// try matching param children
	tvalues := pvalues
	allocated := false
	for _, child := range n.paramChildren {
		if child.minOrder >= order {
			continue
		}
		if data != nil && !allocated {
			tvalues = make([]string, len(pvalues))
			allocated = true
		}
		if d, p, s := child.get(key, tvalues); d != nil && s < order {
			if allocated {
				for i := child.paramIdx; i < len(p); i++ {
					pvalues[i] = tvalues[i]
				}
			}
			data, pnames, order = d, p, s
		}
	}

	return
}

func (n *node) print(level int) string {
	r := fmt.Sprintf("%v{key: %v, regex: %v, data: %v, order: %v, minOrder: %v, paramIdx: %v, params: %v}\n", strings.Repeat(" ", level<<2), n.key, n.regex, n.data, n.order, n.minOrder, n.paramIdx, n.params)
	for _, child := range n.children {
		if child != nil {
			r += child.print(level + 1)
		}
	}
	for _, child := range n.paramChildren {
		r += child.print(level + 1)
	}
	return r
}
