package main

import (
	"os"
	"time"

	"github.com/StephaneRenouard/Go-Nextion/internal/tools"
	"github.com/romana/rlog"
)

func main() {
	var logLevel = "debug"

	os.Setenv("RLOG_LOG_LEVEL", logLevel)
	os.Setenv("RLOG_LOG_NOTIME", "yes")
	rlog.UpdateEnv()

	for {

		m, err := tools.GetSwitchConsumption()
		if err != nil {
			rlog.Info(m.TotalPower)

			tools.WriteTotalPower(m.TotalPower)
		} else {
			tools.WriteTotalPower(-1)
		}

		// Wait forever
		//for {
		//	time.Sleep(1 * time.Second)
		//}

		time.Sleep(1 * time.Second)

	}
}
