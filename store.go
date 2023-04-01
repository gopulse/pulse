package routing

import (
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

	children []*node // child static nodes, indexed by the first byte of each child key

	regex    *regexp.Regexp // regular expression for a param node containing regular expression key
	paramIdx int            // the parameter index, meaningful only for param node
	params   []string       // the parameter names collected from the root till this node
}

// newStore returns a new store
func newStore() *store {
	return &store{
		root: &node{
			static:   true,
			children: make([]*node, 256),
			paramIdx: -1,
			params:   []string{},
		},
	}
}

func (s *store) Add(key string, data interface{}) int {
	s.count++
	return s.root.add(key, data)
}

func (s *store) Get(path string) (data interface{}) {
	data = s.root.get(path)
	return
}

func (s *store) Count() int {
	return s.count
}

func (n *node) add(key string, data interface{}) int {
	if n.static {
		return n.addStatic(key, data)
	}
	return n.addParam(key, data)
}

func (n *node) addStatic(key string, data interface{}) int {
	if len(key) == 0 {
		n.data = data
		return n.paramIdx
	}

	c := key[0]
	if n.children[c] == nil {
		n.children[c] = &node{
			static:   true,
			children: make([]*node, 256),
			paramIdx: -1,
			params:   append(n.params, n.key),
		}
	}

	return n.children[c].addStatic(key[1:], data)
}

func (n *node) addParam(key string, data interface{}) int {
	if len(key) == 0 {
		n.data = data
		return n.paramIdx
	}

	if key[0] == '/' {
		key = key[1:]
	}

	idx := strings.IndexByte(key, '/')
	if idx == -1 {
		idx = len(key)
	}

	var child *node
	pathPart := key[:idx]

	// Check if a child node with the path part exists
	for _, c := range n.children {
		if c.static && c.key == pathPart {
			child = c
			break
		}
	}

	// If no child node exists, create a new one
	if child == nil {
		if pathPart[0] == ':' { // parameterized segment
			child = &node{
				static:   false,
				children: make([]*node, 1),
				paramIdx: n.paramIdx + 1,
				params:   append(n.params, pathPart[1:]),
			}
		} else { // static segment
			child = &node{
				static:   true,
				key:      pathPart,
				children: make([]*node, 256),
			}
		}
		n.children = append(n.children, child)
	}

	// Recurse into the child node with the remaining part of the key
	return child.addParam(key[idx:], data)
}

func (n *node) get(path string) (data interface{}) {
	if n.static {
		return n.getStatic(path)
	}
	return n.getParam(path)
}

func (n *node) getStatic(path string) (data interface{}) {
	if len(path) == 0 {
		return n.data
	}

	c := path[0]
	if n.children[c] == nil {
		return nil
	}

	return n.children[c].getStatic(path[1:])
}

func (n *node) getParam(path string) (data interface{}) {
	if len(path) == 0 {
		return n.data
	}

	if n.children[0] == nil {
		return nil
	}

	return n.children[0].getParam(path)
}

func (n *node) String() string {
	return n.string("")
}

func (n *node) string(indent string) string {
	s := indent + n.key + "\n"
	for _, child := range n.children {
		if child != nil {
			s += child.string(indent + "  ")
		}
	}
	return s
}
