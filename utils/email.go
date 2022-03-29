package utils

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"

	"github.com/ytwxy99/autoCoins/configuration"
)

const (
	Subject = "策略"
)

func SendMail(sysConf *configuration.SystemConf, policy string, body string) error {
	message := gomail.NewMessage()
	port, err := strconv.Atoi(sysConf.Email.Port)
	if err != nil {
		logrus.Error("convert port from string to int failed. the error is ", err)
		return err
	}

	message.SetHeader("From", message.FormatAddress(sysConf.Email.User, "Autocoins量化推荐买点"))
	message.SetHeader("To", sysConf.Email.MailTo...)
	message.SetHeader("Subject", policy+Subject)
	message.SetBody("text/html", body)
	eObject := gomail.NewDialer(sysConf.Email.Host, port, sysConf.Email.User, sysConf.Email.Password)

	return eObject.DialAndSend(message)
}
