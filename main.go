package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var bot *tgbotapi.BotAPI
var service []Service
var webhookURLSlack = "https://hooks.slack.com/services/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"
var webhookTelegram = "1376434732:AAE6YwG6QgnHB_TCFEaM2NnTjANFUsM23dY"
var chat_id = "@rezalert"

type Service struct {
	Command    string
	Regexp     string
	Interval   int
	Name       string
	LastStatus string
}
type SlackRequestBody struct {
	Text string `json:"text"`
}

func check(i int) {
	for {
		t := time.Duration(service[i].Interval) * time.Second
		fmt.Print(service[i].Name)
		out, err := exec.Command("bash", "-c", service[i].Command).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(out), regexp.MustCompile(service[i].Regexp).Match(out))
		if service[i].LastStatus != string(out) {
			//alertTelegram()
			service[i].LastStatus = string(out)
			message := "Message to send"
			params := url.Values{
				"webhook": {webhookTelegram},
				"chat_id": {chat_id},
				"text":    {message},
			}
			fmt.Print(params)
			err := SendSlackNotification(webhookURLSlack, message)
			if err != nil {
				log.Fatal(err)
			}
			err = SendTelegramNotification(params)
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(t)
	}
}

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("service_test.db"), &gorm.Config{})
	if err != nil {
		panic("Connection failed")
	}
	db.AutoMigrate(&Service{})
	//db.Create(&Service{Command: "service ssh status | grep Active", Regexp: "", Interval: 20, Name: "status", LastStatus: ""})
	//db.Create(&Service{Command: "service2 ssh status | grep Active", Regexp: "", Interval: 30, Name: "status2", LastStatus: ""})
	db.Find(&service)
	for i := 0; i < len(service); i++ {
		go check(i)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

func SendTelegramNotification(params url.Values) error {
	var err error
	bot, err = tgbotapi.NewBotAPI(webhookTelegram)
	if err != nil {
		log.Fatal(err)
	}
	msg := tgbotapi.NewMessageToChannel("@rezalert", s)
	msg.ParseMode = "HTML"
	bot.Send(msg)
	return nil
}

func SendSlackNotification(webhookURL string, msg string) error {

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
