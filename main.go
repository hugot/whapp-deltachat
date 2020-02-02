package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hugot/whapp-deltachat/bridge"
	"github.com/hugot/whapp-deltachat/core"
)

func main() {
	argLen := len(os.Args)

	if argLen != 2 {
		fmt.Fprintln(os.Stderr, "Usage: whapp-deltachat CONFIG_FILE")
		return
	}

	configPath := os.Args[1]

	config, err := core.ConfigFromFile(configPath)

	if err != nil {
		log.Fatal(err)
	}

	ensureDirectoryOrDie(config.App.DataFolder)
	ensureDirectoryOrDie(config.App.DataFolder + "/tmp")

	bridge := &bridge.Bridge{}
	err = bridge.Init(config)

	if err != nil {
		log.Fatal(err)
	}

	defer bridge.Close()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case sig := <-wait:
			log.Println(sig)
			return
		}
	}

}

func ensureDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}

	}

	err := os.Chmod(dir, 0700)
	if err != nil {
		return err
	}

	return nil
}

func ensureDirectoryOrDie(dir string) {
	err := ensureDirectory(dir)

	if err != nil {
		log.Fatal(err)
	}
}
