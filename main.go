package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var db *gorm.DB
var bot *tgbotapi.BotAPI
var service []Service
var webhookURLSlack = "YOUR_WEBHOOK"
var webhookTelegram = "YOUR_WEBHOOK"

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
			service[i].LastStatus = string(out)
			db.Where("name = ?", service[i].Name).Update("LastStatus", service[i].LastStatus)
			message := service[i].LastStatus
			err = SendSlackNotification(webhookURLSlack, message)
			if err != nil {
				log.Fatal(err)
			}
			err = SendTelegramNotification(message)
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(t)
	}
}

func main() {
	var err error
	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("KO")
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

func SendTelegramNotification(message string) error {
	var err error
	bot, err = tgbotapi.NewBotAPI(webhookTelegram)
	if err != nil {
		log.Fatal(err)
	}
	msg := tgbotapi.NewMessageToChannel("@rezalert", message)
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
