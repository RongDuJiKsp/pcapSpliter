package main

import (
	"fmt"
	"os"
	"pcapSpliter/flowmaker"
	"pcapSpliter/utils/fileutil"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()
	pathHd, err := fileutil.Mkdir(file)
	if err != nil {
		panic(err)
	}
	if err := flowmaker.MakeSession(file, pathHd); err != nil {
		panic(err)
	}
	fmt.Println("Done")
	os.Exit(0)
}
