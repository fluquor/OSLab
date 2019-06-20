package growth

import (
	"log"
	"math/rand"
)

// 一个频繁项集最多只能产生一条和用户数据匹配的关联规则
func parseSet(set []ItemType, user []ItemType) ([]ItemType, ItemType, bool) {
	setA := make(map[ItemType]bool)
	for _, item := range set {
		setA[item] = true
	}
	setU := make(map[ItemType]bool)
	for _, item := range user {
		setU[item] = true
	}
	if len(setU) < len(setA) {
		return []ItemType{}, NilItem, false
	}
	preItems := make([]ItemType, 0)
	afterItem := NilItem
	for item := range setA {
		if _, ok := setU[item]; !ok {
			if afterItem == NilItem {
				afterItem = item
			} else {
				return preItems, afterItem, false
			}

		} else {
			preItems = append(preItems, item)
		}
	}
	if afterItem == NilItem || len(preItems) == 0 {
		return preItems, afterItem, false
	}
	return preItems, afterItem, true
}

// RecommendItem 根据输入的频繁项集和用户的特征(购买商品向量) 来随机选择一个商品推荐
func RecommendItem(freSets [][]ItemType, user []ItemType) (ItemType, bool) {
	afterSets := make([]ItemType, 0)
	for _, freSet := range freSets {
		if _, afterItem, ok := parseSet(freSet, user); ok {
			afterSets = append(afterSets, afterItem)
		}
	}
	if len(afterSets) == 0 {
		log.Println("没有找到合适的推荐商品")
		return NilItem, false
	}
	_ = rand.Intn(len(afterSets))
	return afterSets[0], true
}
