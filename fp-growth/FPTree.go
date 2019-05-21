package growth

type FPTree struct {
	Root *FPNode
}

func (t *FPTree) Add(trans Transaction) {
	point := t.Root
	for _, item := range trans {
		nextPoint := point.Search(item)
		if nextPoint != nil {
			nextPoint
		}
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
}
