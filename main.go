package main

import (
	growth "asahi/OSLab/fp-growth"
	"flag"
	"log"
)

const (
	defaultDataFile = "data/test.dat"
	defaultTreeFile = "data/tree.gob"
	MinSuppRatio    = 0.0
)

func main() {
	trainDataFile := flag.String("data", defaultDataFile, "Input train data file")
	TreeFile := flag.String("tree", defaultTreeFile, "Save the tree gob file")
	flag.Parse()
	log.Printf("开始读取输入文件: %s\n", *trainDataFile)
	trans, count := growth.BuildTransactions(*trainDataFile)
	log.Printf("读取了 %d 条交易记录 \n", count)
	growth.FindFrequentItemsets(trans, MinSuppRatio, *TreeFile)
	fCount := 0
	log.Printf("Get %d frequend itemsets\n", fCount)
}
