package main

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Service struct {
	Command    string
	Regexp     string
	Interval   int
	Name       string
	LastStatus string
}

func main() {
	var err error
	var service []Service
	db, err = gorm.Open(sqlite.Open("service_test.db"), &gorm.Config{})
	if err != nil {
		panic("Connection failed")
	}
	db.AutoMigrate(&Service{})
	//db.Create(&Service{Command: "service ssh status | grep Active", Regexp: "", Interval: 20, Name: "status", LastStatus: ""})
	//db.Create(&Service{Command: "service2 ssh status | grep Active", Regexp: "", Interval: 30, Name: "status2", LastStatus: ""})
	db.Find(&service)
	var t time.Duration
	t = time.Duration(service[1].Interval) * time.Second
	fmt.Println(t)
	go func() {
		for {
			fmt.Print("CIAO")
			time.Sleep(t * time.Second)
		}
	}()
	fmt.Println("SDC")
	for {
		time.Sleep(1 * time.Second)
	}
	/*
		for i := 0; i < len(service); i++ {
			t := service[i].Interval
			go func() {
				for {
					fmt.Print("CIAO")
					time.Sleep(t * time.Second)
				}
			}()
		}
	*/

}
