package growth

type Route [2]*FPNode
type FPTree struct {
	Root   *FPNode
	routes map[ItemType]Route
}

func (t *FPTree) updateRoute(point *FPNode) {
	if route, ok := t.routes[point.Item]; ok {
		route[1].Neighber = point
		t.routes[point.Item] = Route{route[0], point}
	} else {
		t.routes[point.Item] = Route{point, point}
	}
}

func (t *FPTree) Add(trans Transaction) {
	point := t.Root
	for _, item := range trans {
		nextPoint := point.Search(item)
		if nextPoint != nil {
			nextPoint.Increment()
		} else {
			nextPoint = NewFPNode(item, point)
			point.Add(nextPoint)
			t.updateRoute(nextPoint)
		}
		point = nextPoint
	}
}
func (t FPTree) Items() map[ItemType][]*FPNode {
	result := make(map[ItemType][]*FPNode)
	for key := range t.routes {
		result[key] = t.Nodes(key)
	}
	return result
}
func (t *FPTree) Nodes(item ItemType) []*FPNode {
	if route, ok := t.routes[item]; !ok {
		return []*FPNode{}
	} else {
		result := make([]*FPNode, 0)
		node := route[0]
		result = append(result, node)
		for node != nil {
			node = node.Neighber
		}
		return result
	}
}

func (t *FPTree) PrefixPaths(item ItemType) [][]*FPNode {
	collectPath := func(t *FPNode) {
		path := make([]*FPNode)
		for t != nil && !t.IsRoot() {
			path = append(path, t)
			t = t.Parent
		}
		for i := len(path)/2 - 1; i >= 0; i-- {
			opp := len(path) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		return path
	}
	result := make([][]*FPNode)
	for _, node := range t.Nodes(item) {
		append(result, collectPath(node))
	}
	return result
}
func FindFrequentItemsets(transactions []Transaction, min_supp float64, itemSetChan chan<- []ItemType) {
	min_count := int(float64(len(transactions)) * min_supp)
	items := make(map[ItemType]int)
	for _, trans := range transactions {
		for _, item := range trans {
			items[item]++
		}
	}
	for item, count := range items {
		if count < min_count {
			delete(items, item)
		}
	}

	cleanTrans := func(trans Transaction) Transaction {
		newTrans := make([]ItemType, 0)
		for _, item := range trans {
			if _, ok := items[item]; ok {
				newTrans = append(newTrans, item)
			}
		}
		return newTrans
	}

	master := FPTree{}
	for _, trans := range transactions {
		master.Add(cleanTrans(trans))
	}

	// 利用管道来返回结果
	var findWithSuffix func(*FPTree, []ItemType, chan<- []ItemType)
	findWithSuffix = func(t *FPTree, suffix []ItemType, result chan<- []ItemType) {
		for item, nodes := range t.Items() {
			support := 0
			for _, node := range nodes {
				support += node.Count
			}
			if support > min_count {
				for _, v := range suffix {
					if item == v {
						continue
					} else {
						foundSet := []ItemType{item}
						foundSet = append(foundSet, suffix...)
						result <- foundSet
					}
				}

			}

		}
	}
}

func ConditionalTreeFromPaths(paths [][]*FPNode) {
	tree := FPTree{}
	var condItem = NilItem
	items := make(map[ItemType]bool)
	for _, path := range paths {
		if condItem == NilItem {
			condItem = path[len(path)-1].Item
		}
		point := tree.Root
		for _, node := range path {
			nextPoint := point.Search(node.Item)
			if nextPoint == nil {
				items[node.Item] = true
				count := 0
				if node.Item == condItem {
					count = node.Count
				}
				nextPoint = &FPNode{Tree: &tree, Item: node.Item, Count: 1}
				point.Add(nextPoint)
				tree.updateRoute()
			}
		}
	}
}
