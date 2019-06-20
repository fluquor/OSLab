package main

import (
	growth "asahi/OSLab/fp-growth"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	defaultDataFile = "data/trans_10W.txt"
	defaultTreeFile = "data/tree.gob"
	defaultUserFile = "data/test_asahi.txt"

	setsFile         = "data/result_sets.txt"
	defautMaxRecords = 10000
	defaultMaxUsers  = 10
	MinSuppRatio     = 0.25
)

func PrintSets(s [][]growth.ItemType) {
	for _, set := range s {
		fmt.Println(set)
	}
}

// func saveSets(sets [][]ItemType) {

// }

func main() {
	trainDataFile := flag.String("data", defaultDataFile, "Input train data file")
	userDataFile := flag.String("user", defaultUserFile, "Input user data")
	TreeFile := flag.String("tree", defaultTreeFile, "Save the tree gob file")
	flag.Parse()

	// 生成频繁项集
	log.Printf("开始读取输入文件: %s\n", *trainDataFile)
	trans, count := growth.BuildTransactions(*trainDataFile, defautMaxRecords)
	log.Printf("读取了 %d 条交易记录 \n", count)
	result := growth.FindFrequentItemsets(trans, MinSuppRatio, *TreeFile)
	log.Printf("生成频繁项集规模: %d\n", len(result))

	// 保存频繁项集结果到文件
	f, err := os.OpenFile(setsFile, os.O_CREATE|os.O_RDWR, 0755)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	w := bufio.NewWriter(f)
	for _, line := range result {
		for _, item := range line {
			w.WriteString(fmt.Sprintf("%d ", item))

		}
		w.WriteString(fmt.Sprintf("\n"))
	}

	// 读取用户数据并推荐
	userRecors, _ := growth.BuildTransactions(*userDataFile, defaultMaxUsers)
	for i, user := range userRecors {
		if recommendItem, ok := growth.RecommendItem(result, user); ok {
			log.Printf("为第%d个用户推荐商品为:%d", i+1, recommendItem)
		} else {
			log.Printf("为第%d个用户 无推荐商品!")
		}

	}
}
