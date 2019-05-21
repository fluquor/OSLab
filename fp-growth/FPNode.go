package growth

type ItemType string
type Transaction []ItemType

type FPNode struct {
	Item     ItemType
	Count    int //支持度计数
	Parent   *FPNode
	Children map[ItemType]*FPNode
}

func NewFPNode(item ItemType, p *FPNode) *FPNode {
	return &FPNode{Item: item, Count: 1, Parent: p}
}

func (n *FPNode) Add(child *FPNode) {
	if _, ok := n.Children[child.Item]; !ok {
		n.Children[child.Item] = child
		child.Parent = n
	}
}

func (n FPNode) Search(item ItemType) *FPNode {
	if _, ok := n.Children[item]; !ok {
		return n.Children[item]
	}
	return nil
}

func (n *FPNode) Increment() {
	if n.Count == -1 {
		panic("Root nodes have no associated count.")
	} else {
		n.Count++
	}
}

func (n FPNode) IsRoot() bool {
	return false
}
