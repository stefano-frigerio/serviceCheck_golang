package rezbot

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/iesreza/foundation/system"
	"github.com/iesreza/gutil/log"
)

var db = system.GetDBO()
var bot *tgbotapi.BotAPI

func Schedule() {
	var err error
	bot, err = tgbotapi.NewBotAPI("1376434732:AAE6YwG6QgnHB_TCFEaM2NnTjANFUsM23dY")
	if err != nil {
		log.Fatal(err)
	}

	for {
		list := []Endpoint{}
		db.Where("active = ?", true).Find(&list)
		for _, item := range list {
			if err := item.CheckHealth(); err != nil {
				Error("App internal Error >\r\n" + err.Error())
			}
		}
		time.Sleep(30 * time.Minute)
	}

}

func Info(s string) {
	s = "**INFO:**\r\n" + s
	message(s)
}
func Error(s string) {
	s = "<b>ERROR:</b>\r\n" + s
	message(s)
}
func message(s string) {
	msg := tgbotapi.NewMessageToChannel("@rezalert", s)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func ErrorStatus(e *Endpoint, status Status, msg string) {
	Error("<b>" + msg + "</b>" + Parse("<b>Endpoint</b>: "+e.Name+"\n<b>IP..........</b>:"+e.IP) + "\n" + Parse(PrettyPrint(status)) + "\n" + "http://" + e.IP + ":83/api/health?key=" + e.Key)
}

func PrettyPrint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		s := strings.Replace(string(b), "\"", "", -1)
		s = strings.Replace(s, ",", "", -1)
		s = strings.Replace(s, "{", "", -1)
		s = strings.Replace(s, "}", "", -1)
		return s
	}
	return ""
}

func Parse(s string) string {
	lines := strings.Split(s, "\n")
	output := ""
	maxwidth := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts[0]) > maxwidth {
			maxwidth = len(strings.TrimSpace(parts[0]))
		}

	}
	maxwidth += 4

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			output += "\n"
			continue
		}
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts[0]) < maxwidth {
			for len(parts[0]) < maxwidth+1 {
				parts[0] += "."
			}
		}
		output += "\n" + parts[0]
		if len(parts) > 1 {
			v := strings.Join(parts[1:], ":")
			if _, err := strconv.Atoi(v); err == nil {
				v += "."
			}
			output += " " + strings.TrimSpace(v)
		}
	}

	return "<pre>" + output + "</pre>"
}
