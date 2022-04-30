package main

import (
	"github.com/ytwxy99/autocoins/pkg/system"
	"github.com/ytwxy99/autocoins/pkg/utils"
)

func main() {
	authConf, _ := utils.ReadGateAPIV4("./etc/auth.yml")
	sysConf, _ := utils.ReadSystemConfig("./etc/autoCoin.yml")
	system.Init(authConf, sysConf)
}
