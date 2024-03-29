package container

type RBColor bool

const (
	BLACK RBColor = true
	RED   RBColor = false
)

type RBNode struct {
	Key         string
	Value       interface{}
	Color       RBColor
	Left, Right *RBNode
}

type RedBlackTree struct {
	root *RBNode
	size int
}

func (t *RedBlackTree) isRed(node *RBNode) bool {
	if node == nil {
		return false
	}
	return node.Color == RED
}

func (t *RedBlackTree) rotateLeft(node *RBNode) *RBNode {
	x := node.Right
	node.Right = x.Left
	x.Left = node
	x.Color = node.Color
	node.Color = RED
	return x
}

func (t *RedBlackTree) rotateRight(node *RBNode) *RBNode {
	x := node.Left
	node.Left = x.Right
	x.Right = node
	x.Color = node.Color
	node.Color = RED
	return x
}

func (t *RedBlackTree) flipColors(node *RBNode) {
	node.Color = RED
	node.Left.Color = BLACK
	node.Right.Color = BLACK
}

func (t *RedBlackTree) Put(key string, value interface{}) {
	t.root = t.put(t.root, key, value)
	t.root.Color = BLACK
}

func (t *RedBlackTree) put(node *RBNode, key string, value interface{}) *RBNode {
	if node == nil {
		t.size++
		return &RBNode{Key: key, Value: value, Color: RED}
	}

	if key < node.Key {
		node.Left = t.put(node.Left, key, value)
	} else if key > node.Key {
		node.Right = t.put(node.Right, key, value)
	} else {
		node.Value = value
	}

	if t.isRed(node.Right) && !t.isRed(node.Left) {
		node = t.rotateLeft(node)
	}
	if t.isRed(node.Left) && t.isRed(node.Left.Left) {
		node = t.rotateRight(node)
	}
	if t.isRed(node.Left) && t.isRed(node.Right) {
		t.flipColors(node)
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
