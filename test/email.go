package main

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var Subject = "策略"
var EmailTo = []string{
	"524737353@qq.com",
}
var Port = "465"
var User = "autocoins@163.com"
var Passwd = "FHYKNBEUYNBYRGAU"
var Host = "smtp.163.com"

func SendMail(policy string, body string) error {
	message := gomail.NewMessage()
	port, err := strconv.Atoi(Port)
	if err != nil {
		logrus.Error("convert port from string to int failed. the error is ", err)
		return err
	}

	message.SetHeader("From", message.FormatAddress(User, "Autocoins量化推荐买点"))
	message.SetHeader("To", EmailTo...)
	message.SetHeader("Subject", policy+Subject)
	message.SetBody("text/html", body)
	eObject := gomail.NewDialer(Host, port, User, Passwd)

	return eObject.DialAndSend(message)
}

func main() {
	err := SendMail("测试策略", "Just a test")
	if err != nil {
		logrus.Error("send a test email failed, the err is ", err)
	} else {
		logrus.Info("send a test email success")
	}
}
