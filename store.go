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
	if key[0] == '/' {
		key = key[1:]
	}
	p0 := strings.Index(key, ":")
	if p0 == -1 {
		// no param tokens found: create a static node
		child := &node{
			static:        true,
			key:           key,
			minOrder:      order,
			children:      make([]*node, 256),
			paramChildren: make([]*node, 0),
			paramIdx:      n.paramIdx,
			params:        n.params,
			data:          data,
			order:         order,
		}
		n.children[key[0]] = child
		if n.data == nil {
			// if the node is already a data node, ignore the new data since we only care about the first matched node
			n.data = data
			n.order = order
			return n.paramIdx + 1
		}
		return child.paramIdx + 1
	}
	// param token found: create a static node for characters before the param token
	child := &node{
		static:        true,
		key:           key[:p0],
		minOrder:      order,
		children:      make([]*node, 256),
		paramChildren: make([]*node, 0),
		paramIdx:      n.paramIdx,
		params:        n.params,
	}
	n.children[key[0]] = child
	n = child
	key = key[p0:]

	// add param node for the current param token
	p1 := strings.Index(key, "/")
	if p1 == -1 {
		// the param token is at the end of the key
		p1 = len(key)
	}
	pname := key[1:p1]
	pattern, err := regexp.Compile("[^/]+")
	if err != nil {
		// invalid param regex
		return -1
	}
	child = &node{
		static:        false,
		key:           pname,
		minOrder:      order,
		children:      make([]*node, 256),
		paramChildren: make([]*node, 0),
		paramIdx:      n.paramIdx + len(n.paramChildren) + 1,
		params:        append(n.params, pname),
		regex:         pattern,
	}
	n.paramChildren = append(n.paramChildren, child)

	if p1 == len(key) {
		// the param token is at the end of the key
		child.data = data
		child.order = order
		return child.paramIdx + 1
	}

	// process the rest of the key recursively
	n = child
	key = key[p1:]
	return n.addChild(key, data, order)
}

func (n *node) get(key string, pvalues []string) (data interface{}, pnames []string, order int) {
	order = math.MaxInt32
	for len(key) > 0 {
		if n.static {
			// check if the node key is a prefix of the given key
			// a slightly optimized version of strings.HasPrefix
			nkl := len(n.key)
			if nkl > len(key) || n.key != key[:nkl] {
				return
			}
			key = key[nkl:]
		} else if n.regex != nil {
			// param node with regular expression
			match := n.regex.FindStringIndex(key)
			if match == nil || match[0] != 0 {
				return
			}
			if n.paramIdx >= len(pvalues) {
				pvalues = append(pvalues, make([]string, n.paramIdx-len(pvalues)+1)...)
			}
			pvalues[n.paramIdx] = key[0:match[1]]
			key = key[match[1]:]
		} else {
			// param node matching non-"/" characters
			i := strings.IndexByte(key, '/')
			if i == -1 {
				if n.paramIdx >= len(pvalues) {
					pvalues = append(pvalues, make([]string, n.paramIdx-len(pvalues)+1)...)
				}
				pvalues[n.paramIdx] = key
				key = ""
			} else {
				if n.paramIdx >= len(pvalues) {
					pvalues = append(pvalues, make([]string, n.paramIdx-len(pvalues)+1)...)
				}
				pvalues[n.paramIdx] = key[:i]
				key = key[i:]
			}
		}

		// find a static child that can match the rest of the key
		if child := n.children[key[0]]; child != nil {
			if len(n.paramChildren) == 0 {
				// use iteration instead of recursion when there are no param children
				n = child
				continue
			}
			data, pnames, order = child.get(key, pvalues)
		}

		break
	}

	// capture data from this node, if any
	if n.data != nil && (len(key) == 0 || len(n.paramChildren) == 0 && n.static) {
		if n.order < order {
			data, pnames, order = n.data, n.params, n.order
		}
	}
	// try matching param children
	for _, child := range n.paramChildren {
		if child.minOrder >= order {
			continue
		}
		tvalues := make([]string, len(pvalues))
		copy(tvalues, pvalues)
		if d, p, s := child.get(key, tvalues); d != nil && s < order {
			data, pnames, order = d, p, s
			copy(pvalues[child.paramIdx:], tvalues[child.paramIdx:])
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
