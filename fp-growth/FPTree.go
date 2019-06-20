package growth

// 代码参考:https://github.com/enaeseth/python-fp-growth/blob/master/fp_growth.py
import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	// DEBUG 是否处于debug状态
	DEBUG = true
)

// Route 结构保存了包含某个物品的路径的首尾两个节点
type Route [2]*FPNode

// FPTree 是一个用来存储所有交易记录的树结构
type FPTree struct {
	Root   *FPNode
	Routes map[ItemType]Route
}

// BuildTransactions 读取输入文件并解析,返回交易列表和成功读取的交易数目
func BuildTransactions(filename string, lines int) ([]Transaction, int) {
	result := make([]Transaction, 0)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b := make([]byte, 0, 64*1024)
	scanner := bufio.NewScanner(f)
	scanner.Buffer(b, 1024*1024)
	count := 0
	for scanner.Scan() {
		t := make([]ItemType, 0)
		ss := strings.Split(scanner.Text(), " ")
		for _, s := range ss {
			// log.Println(s)
			i1, err := strconv.Atoi(s)
			if err != nil {
				log.Println(err)
				continue
			}
			t = append(t, ItemType(i1))
		}
		result = append(result, t)
		count++
		if lines > 0 && count >= lines {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result, count
}

// NewFPTree 生成一个空的FPTree
func NewFPTree() *FPTree {
	root := &FPNode{Item: NilItem, Count: NilCount, Children: make(map[ItemType]*FPNode)}
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

// Add 向FPTree插入交易数据
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

// Nodes 返回所有包含目标物品的节点
func (t *FPTree) Nodes(item ItemType) []*FPNode {
	if _, ok := t.Routes[item]; !ok {
		return []*FPNode{}
	}
	route := t.Routes[item]
	result := make([]*FPNode, 0)
	node := route[0]
	for node != nil {
		result = append(result, node)
		node = node.Neighbor
	}
	return result
}

// PrefixPaths 返回所有以给定物品为结尾的所有路径
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

// SaveTreeToFile 保存树结构到文件中
func SaveTreeToFile(r *FPTree, filename string) bool {
	gob.Register(FPTree{})
	treeFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
	}
	enc := gob.NewEncoder(treeFile)
	err2 := enc.Encode(r)
	if err2 != nil {
		fmt.Println(err2)
		return false
	}
	return true
}

// FindFrequentItemsets 挖掘频繁项集
func FindFrequentItemsets(transactions []Transaction, minSuppRatio float64, treeSaveFile string) [][]ItemType {
	minSuppCount := int(float64(len(transactions)) * minSuppRatio)
	fmt.Printf("最小支持出现次数: %d \n", minSuppCount)
	items := make(map[ItemType]int)
	for _, trans := range transactions {
		for _, item := range trans {
			items[item]++
		}
	}
	log.Printf("扫描了%d个商品", len(items))
	for item, count := range items {
		if count < minSuppCount {
			delete(items, item)
		}
	}
	log.Printf("清洗后剩余%d个商品", len(items))
	// log.Println(items)
	// 根据supp阈值对所有交易进行清洗,删除低于阈值的所有物品并从大到小进行排序
	cleanTrans := func(trans Transaction) Transaction {
		newTrans := make([]ItemType, 0)
		for _, item := range trans {
			if _, ok := items[item]; ok {
				newTrans = append(newTrans, item)
			}
		}
		sort.Slice(newTrans, func(i, j int) bool {
			return items[newTrans[i]] > items[newTrans[j]]
		})
		return newTrans
	}

	// 构造FP-Tree
	master := NewFPTree()
	for _, trans := range transactions {
		master.Add(cleanTrans(trans))
	}
	log.Printf("Tree build complete with %d children\n", len(master.Root.Children))
	// 保存Tree结果到本地
	// if treeSaveFile != "" {
	// 	log.Printf("开始保存Tree结构到本地...\n")
	// 	if ok := SaveTreeToFile(master, treeSaveFile); ok {
	// 		log.Printf("Tree结构已经保存到 %s", treeSaveFile)
	// 	}
	// }

	// 利用管道来返回结果
	itemSetChan := make(chan []ItemType, 100)
	n := &sync.WaitGroup{}
	var findWithSuffix func(*FPTree, []ItemType, chan<- []ItemType, *sync.WaitGroup)
	findWithSuffix = func(t *FPTree, suffix []ItemType, result chan<- []ItemType, n *sync.WaitGroup) {
		for item, nodes := range t.Items() {

			n.Add(1)
			go func(item ItemType, nodes []*FPNode, result chan<- []ItemType, n *sync.WaitGroup) {
				defer n.Done()
				support := 0
				for _, node := range nodes {
					support += node.Count
				}
				// log.Printf("Support: %d\n", support)
				inFlag := false
				if support >= minSuppCount {
					for _, v := range suffix {
						if item == v {
							inFlag = true
							break
						}
					}
					if !inFlag {
						foundSet := []ItemType{item}
						foundSet = append(foundSet, suffix...)
						result <- foundSet
						// log.Println("已找到一个集合并放入channel")
						condTree := ConditionalTreeFromPaths(t.PrefixPaths(item))
						findWithSuffix(condTree, foundSet, result, n)
					}

				}
			}(item, nodes, result, n)
		}

	}
	findWithSuffix(master, []ItemType{}, itemSetChan, n)
	// 等待所有goroutine都结束,关闭结果channel
	go func() {
		n.Wait()
		close(itemSetChan)
	}()
	results := make([][]ItemType, 0)
	for itemSet := range itemSetChan {

		results = append(results, itemSet)
	}
	return results
}

// ConditionalTreeFromPaths 根据条件模式基构建条件FPTree
func ConditionalTreeFromPaths(paths [][]*FPNode) *FPTree {
	tree := NewFPTree()
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
