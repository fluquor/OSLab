package main

import (
	growth "asahi/OSLab/fp-growth"
	"flag"
	"fmt"
	"log"
)

const (
	defaultDataFile = "data/trans_10W.txt"
	defaultTreeFile = "data/tree.gob"
	defaultNums     = 10000
	MinSuppRatio    = 0.093
)

func PrintSets(s [][]growth.ItemType) {
	for _, set := range s {
		fmt.Println(set)
	}
}

func main() {
	trainDataFile := flag.String("data", defaultDataFile, "Input train data file")
	TreeFile := flag.String("tree", defaultTreeFile, "Save the tree gob file")
	flag.Parse()
	log.Printf("开始读取输入文件: %s\n", *trainDataFile)
	trans, count := growth.BuildTransactions(*trainDataFile, defaultNums)
	log.Printf("读取了 %d 条交易记录 \n", count)
	result := growth.FindFrequentItemsets(trans, MinSuppRatio, *TreeFile)
	log.Printf("生成频繁项集规模: %d\n", len(result))
	// PrintSets(result)
}
