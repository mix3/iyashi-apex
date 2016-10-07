package main

import (
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	iyashi, err := NewIyashi()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	iyashi.Run()
}
