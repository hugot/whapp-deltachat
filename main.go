package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	argLen := len(os.Args)

	if argLen != 2 {
		fmt.Fprintln(os.Stderr, "Usage: whapp-deltachat CONFIG_FILE")
		return
	}

	configPath := os.Args[1]

	config, err := ConfigFromFile(configPath)

	if err != nil {
		log.Fatal(err)
	}

	err = ensureDirectory(config.App.DataFolder)

	if err != nil {
		log.Fatal(err)
	}

	dcClient := BootstrapDcClientFromConfig(*config)

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt)

	for {
		select {
		case sig := <-wait:
			log.Println(sig)

			// Give dc an opportunity to perform some close-down logic
			// and close it's db etc.
			dcClient.Close()
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
