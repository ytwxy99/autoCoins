package policy

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Cointegration struct{}

// find buy point by doing cointegration
func (*Cointegration) Target(args ...interface{}) interface{} {
	// convert specified type
	dbPath := args[0].(string)
	scriptPath := args[1].(string)
	coinCsv := args[2].(string)

	cmd := exec.Command("python3", scriptPath, dbPath, coinCsv)
	output, err := cmd.Output()
	if err != nil {
		logrus.Error("run cointegration python srcipt error:", err)
		return err
	}

	fmt.Println(string(output))

	return nil
}
