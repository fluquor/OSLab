package growth

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Route [2]*FPNode
type FPTree struct {
	Root   *FPNode
	Routes map[ItemType]Route
}

func BuildTransactions(filename string) ([]Transaction, int) {
	result := make([]Transaction, 0)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		t := make([]ItemType, 0)
		count++
		ss := strings.Split(scanner.Text(), " ")
		for _, s := range ss {
			i1, err := strconv.Atoi(s)
			if err != nil {
				log.Fatalln(err)
				continue
			}
			t = append(t, ItemType(i1))
		}
		result = append(result, t)
	}
	return result, count
}

func NewFPTree() *FPTree {
	root := &FPNode{Children: make(map[ItemType]*FPNode)}
	root.Parent = root
	return &FPTree{Root: root, Routes: make(map[ItemType]Route)}
}
func (t *FPTree) updateRoute(point *FPNode) {
	if route, ok := t.Routes[point.Item]; ok {
		route[1].Neighbor = point
		t.Routes[point.Item] = Route{route[0], point}
	} else {
		t.Routes[point.Item] = Route{point, point}
	}
}

func (t *FPTree) Add(trans Transaction) {
	point := t.Root
	for _, item := range trans {
		nextPoint := point.Search(item)
		if nextPoint != nil {
			nextPoint.Increment()
		} else {
			nextPoint = &FPNode{Tree: t, Item: item, Count: 1, Parent: point, Children: make(map[ItemType]*FPNode)}

			point.Add(nextPoint)
			t.updateRoute(nextPoint)
		}
		point = nextPoint
	}
}
func (t FPTree) Items() map[ItemType][]*FPNode {
	result := make(map[ItemType][]*FPNode)
	for key := range t.Routes {
		result[key] = t.Nodes(key)
	}
	return result
}
func (t *FPTree) Nodes(item ItemType) []*FPNode {
	if _, ok := t.Routes[item]; !ok {
		return []*FPNode{}
	}
	route := t.Routes[item]
	result := make([]*FPNode, 0)
	node := route[0]
	result = append(result, node)
	for node != nil {
		node = node.Neighbor
	}
	return result
}

func (t *FPTree) PrefixPaths(item ItemType) [][]*FPNode {
	collectPath := func(t *FPNode) []*FPNode {
		path := make([]*FPNode, 0)
		for t != nil && !t.IsRoot() {
			path = append(path, t)
			t = t.Parent
		}
		for i := len(path)/2 - 1; i >= 0; i-- {
			opp := len(path) - 1 - i
			path[i], path[opp] = path[opp], path[i]
		}
		return path
	}
	result := make([][]*FPNode, 0)
	for _, node := range t.Nodes(item) {
		result = append(result, collectPath(node))
	}
	return result
}
func FindFrequentItemsets(transactions []Transaction, minSuppRatio float64, itemSetChan chan<- []ItemType) {
	minSuppCount := int(float64(len(transactions)) * minSuppRatio)
	items := make(map[ItemType]int)
	for _, trans := range transactions {
		for _, item := range trans {
			items[item]++
		}
	}
	for item, count := range items {
		if count < minSuppCount {
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

	master := NewFPTree()
	for _, trans := range transactions {
		master.Add(cleanTrans(trans))
	}
	fmt.Printf("Tree build complete with %d children\n", len(master.Root.Children))
	// 利用管道来返回结果
	var findWithSuffix func(*FPTree, []ItemType, chan<- []ItemType)
	findWithSuffix = func(t *FPTree, suffix []ItemType, result chan<- []ItemType) {
		for item, nodes := range t.Items() {
			support := 0
			for _, node := range nodes {
				support += node.Count
			}
			if support > minSuppCount {
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

	findWithSuffix(master, []ItemType{}, itemSetChan)
	close(itemSetChan)
}

func ConditionalTreeFromPaths(paths [][]*FPNode) *FPTree {
	tree := &FPTree{}
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
				nextPoint = &FPNode{Tree: tree, Item: node.Item, Count: count, Children: make(map[ItemType]*FPNode)}
				point.Add(nextPoint)
				tree.updateRoute(nextPoint)
			}
			point = nextPoint
		}
	}

	for _, path := range tree.PrefixPaths(condItem) {
		count := path[len(path)-1].Count
		for i := 0; i < len(path)-1; i++ {
			path[i].Count += count
		}
	}
	return tree
}
