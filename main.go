package main

import (
	growth "asahi/OSLab/fp-growth"
	"flag"
	"fmt"
)

const (
	defaultDataFile = "data/webdocs.dat"
	MinSuppRatio    = 0.092
)

func main() {
	trainDataFile := flag.String("data", defaultDataFile, "Input train data file")
	flag.Parse()
	fmt.Printf("Parsing data file: %s\n", *trainDataFile)
	trans, count := growth.BuildTransactions(*trainDataFile)
	fmt.Printf("Parsed %d trans \n", count)

	itemSetChan := make(chan []growth.ItemType, 100)
	go growth.FindFrequentItemsets(trans, MinSuppRatio, itemSetChan)
	fCount := 0
	for range itemSetChan {
		fCount++
	}
	fmt.Printf("Get %d frequend itemsets\n", fCount)
}
