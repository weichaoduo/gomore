package main

import (
	"fmt"

	log2 "github.com/Sirupsen/logrus"
	"github.com/weekface/mgorus"
)

func main() {
	log2.SetLevel(log2.ErrorLevel)
	log := log2.New()

	hooker, err := mgorus.NewHooker("localhost:27017", "db", "collection")
	if err == nil {
		log.Hooks.Add(hooker)
	} else {
		fmt.Println("mongodb err:", err)

	}

	log.WithFields(log2.Fields{
		"name": "zhangsan1215555555551155555",
		"age":  28225552,
	}).Info("Hello world!")

	log.Warn("2222222222221111122222")
}
