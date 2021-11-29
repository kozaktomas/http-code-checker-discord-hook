package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const discordMessage = "__**URL check**__\n**URL:** %s\n**CODE:** %d"

var (
	url     = kingpin.Arg("url", "Requested URL address.").Required().String()
	code    = kingpin.Arg("code", "Expected HTTP status code.").Required().Int()
	discord = kingpin.Arg("discord_hook_url", "Discord webhook url.").Required().String()
	sleep   = kingpin.Flag("sleep", "Duration between checks (2s, 5m, 10h, 2d). Default 5m.").Short('s').Default("5m").String()
	verbose = kingpin.Flag("verbose", "Verbose mode. Default false.").Short('v').Bool()
)

func main() {
	kingpin.Parse()

	// first check (endpoint validity)
	_, err := checkUrl()
	if err != nil {
		fmt.Println(fmt.Sprintf("could not get response from %s\n", *url))
		return
	}

	// parse sleep duration - time between checks
	duration, err := time.ParseDuration(*sleep)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not parse sleep duration %s", *sleep))
		return
	}

	// generate tick for every check
	ticker := make(chan interface{})
	go func() {
		for {
			ticker <- true
			log("ticker - new check initialized")
			time.Sleep(duration)
		}
	}()

	// receive cancellation signals using channel
	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// process ticks and cancellation signals
	for {
		select {
		case _ = <-cancel:
			// graceful shutdown
			log("The end is near!")
			return
		case _ = <-ticker:
			log("start new check")
			match := checkAndValidate()
			log("check finished")
			if match {
				log("Expected status code returned. Quiting.")
				return
			}
		}
	}
}

// checkAndValidate checks requested url and return if the condition pass
// returns true when condition pass and program should quit
func checkAndValidate() bool {
	c, err := checkUrl()
	if err != nil {
		log("could not check url")
		return false
	}

	// expected code match
	if c == *code {
		log("match found ")
		err = sendWebhook()
		if err != nil {
			log(err.Error())
			return false
		}

		return true
	}

	return false
}

// checkUrl checks requested url and returns status_code and error
func checkUrl() (int, error) {
	resp, err := http.Get(*url)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

// sendWebhook sends Discord webhook
// Discord shows message in channel related to the webhook
func sendWebhook() error {
	body := struct {
		Content string `json:"content"`
	}{
		Content: fmt.Sprintf(discordMessage, *url, *code),
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not marshal data for Discord webhook")
	}
	data := bytes.NewBuffer(jsonData)

	resp, err := http.Post(*discord, "application/json", data)
	if err != nil {
		return fmt.Errorf("could not send Discord webhook: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid response code from Discord webhook: %d", resp.StatusCode)
	}

	return nil
}

func log(msg string) {
	if *verbose {
		fmt.Println(msg)
	}
}
