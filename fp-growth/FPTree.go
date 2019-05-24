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
			nextPoint = &FPNode{Item: item}
			point.Add(nextPoint)
			t.updateRoute(nextPoint)
		}
		point = nextPoint
	}
}
func (t FPTree) Items()map[ItemType][]*ItemType  {
	result:=make(map[ItemType][]*ItemType)
	for key := range t.routes{
		result[key]=t.Nodes(key)
	}
	return result
}
func (t *FPTree) Nodes(item ItemType) []*ItemType{
	if route,ok:=t.routes[item];!ok{
		return []*ItemType{}
	}else{
		result:=make([]*ItemType,0)
		node:=route[0]
		result=append(result,node)
		for node!=nil {
			node=node.Neighber
		}
		return result
	}
}

func FindFrequentItemsets(transactions []Transaction, min_supp float64) {
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

	// TODO:
	var findWithSuffix func(*FPTree, []ItemType)
	findWithSuffix = func(t *FPTree, suffix []ItemType) {
		for item,nodes := range t.Items()
	}
	_ = findWithSuffix
}
