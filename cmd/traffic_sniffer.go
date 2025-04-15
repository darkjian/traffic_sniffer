package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
	"traffic_sniffer/pkg/sniffer"
)

const (
	buffSize = 2048
)

func main() {

	filename := fmt.Sprintf("./logs/%v.log", time.Now().Unix())

	// creating log file
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("creating file error: %v", err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal("close file error:", err.Error())
		}
	}()

	// setting log output to file
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// run traffic sniffer
	go sniffer.New(buffSize).Run()

	go func() {
		// starting traffic sending script
		cmd := exec.Command("/usr/local/bin/setup_and_test.sh")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("exec script error: %s, output: %s", err.Error(), string(output))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGINT, syscall.SIGHUP)
	<-c
}
