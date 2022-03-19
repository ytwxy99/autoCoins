package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/ytwxy99/autoCoins/configuration"
)

// read Key&Secret from yaml
func ReadGateAPIV4(filePath string) (*configuration.GateAPIV4, error) {
	var gateAPIV4 configuration.GateAPIV4
	value, err := os.Open(filePath)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	decode := yaml.NewDecoder(value)
	decode.Decode(&gateAPIV4)
	defer value.Close()

	return &gateAPIV4, nil
}

// read system configure
func ReadSystemConfig(filePath string) (*configuration.SystemConf, error) {
	var sysConf configuration.SystemConf
	yamlFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		logrus.Error("ReadSystemConfig -> Get err: %v", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &sysConf)
	if err != nil {
		logrus.Error("ReadSystemConfig -> Unmarshal: %v", err)
		return nil, err
	}
	return &sysConf, nil
}
