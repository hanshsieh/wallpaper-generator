package main

import (
	"github.com/you/hello/pkg/consumer"
	"github.com/you/hello/pkg/producer"
	"log"
	"os"
	"strconv"
)

func main() {
	srcDir := os.Args[1]
	dstDir := os.Args[2]
	width, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal("invalid width")
	}
	height, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatal("invalid height")
	}
	imgProducer := producer.NewProducer(srcDir)
	imgConsumer := consumer.NewConsumer(
		dstDir,
		float64(width) / float64(height))
	err = imgProducer.Start()
	if err != nil {
		log.Fatalf("faild to start producer: %v", err)
	}
	done := false
	for !done {
		select {
		case entry := <- imgProducer.Entries():
			log.Printf("Name: %s", entry.Name)
			err := imgConsumer.PutEntry(entry)
			if err != nil {
				log.Printf("failed to process image %q: %v", entry.Name, err)
			}
		case err := <- imgProducer.Errors():
			log.Printf("error: %v", err)
		case <-imgProducer.Done():
			log.Print("Done")
			done = true
		}
	}
}