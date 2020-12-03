package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var service []Service

type Service struct {
	Command    string
	Regexp     string
	Interval   int
	Name       string
	LastStatus string
}

func check(i int) {
	for {
		t := time.Duration(service[i].Interval) * time.Second
		fmt.Print(service[i].Name)
		/*out, err := exec.Command(service[i].Command).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("Return", out)
		if out !=service[i].LastStatus {
			service[i].LastStatus = string(out)
			alertTelegram();
		}
		*/
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
		fmt.Println(i)
		go check(i)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}
