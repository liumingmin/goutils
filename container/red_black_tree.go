package container

type RBColor bool

const (
	BLACK RBColor = true
	RED   RBColor = false
)

type Node struct {
	Key         string
	Value       interface{}
	Color       RBColor
	Left, Right *Node
}

type RedBlackTree struct {
	root *Node
	size int
}

func isRed(node *Node) bool {
	if node == nil {
		return false
	}
	return node.Color == RED
}

func rotateLeft(node *Node) *Node {
	x := node.Right
	node.Right = x.Left
	x.Left = node
	x.Color = node.Color
	node.Color = RED
	return x
}

func rotateRight(node *Node) *Node {
	x := node.Left
	node.Left = x.Right
	x.Right = node
	x.Color = node.Color
	node.Color = RED
	return x
}

func flipColors(node *Node) {
	node.Color = RED
	node.Left.Color = BLACK
	node.Right.Color = BLACK
}

func (t *RedBlackTree) Put(key string, value interface{}) {
	t.root = t.put(t.root, key, value)
	t.root.Color = BLACK
}

func (t *RedBlackTree) put(node *Node, key string, value interface{}) *Node {
	if node == nil {
		t.size++
		return &Node{Key: key, Value: value, Color: RED}
	}

	if key < node.Key {
		node.Left = t.put(node.Left, key, value)
	} else if key > node.Key {
		node.Right = t.put(node.Right, key, value)
	} else {
		node.Value = value
	}

	if isRed(node.Right) && !isRed(node.Left) {
		node = rotateLeft(node)
	}
	if isRed(node.Left) && isRed(node.Left.Left) {
		node = rotateRight(node)
	}
	if isRed(node.Left) && isRed(node.Right) {
		flipColors(node)
	}

	return node
}

func (t *RedBlackTree) Get(key string) interface{} {
	node := t.root
	for node != nil {
		if key < node.Key {
			node = node.Left
		} else if key > node.Key {
			node = node.Right
		} else {
			return node.Value
		}
	}
	return nil
}
