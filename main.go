package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/energieip/common-components-go/pkg/dswitch"
	"github.com/romana/rlog"
	"github.com/tarm/serial"
)

func main() {

	var serialPortNum = "/dev/ttyAMA0"
	var serialPortSpeed = 9600
	var messageVariable = "0"
	var message = "page1.t4.txt=\"" + messageVariable + "\"" + "000"

	var serverURL = "127.0.0.1"
	var serverPORT = "8888"

	var logLevel = "debug"

	os.Setenv("RLOG_LOG_LEVEL", logLevel)
	os.Setenv("RLOG_LOG_NOTIME", "yes")
	rlog.UpdateEnv()

	for {

		url := "https://" + serverURL + ":" + serverPORT + "/v1.0/status/consumptions"

		req, _ := http.NewRequest("GET", url, nil)
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
		}
		client := &http.Client{Transport: transCfg}
		resp, err := client.Do(req)

		if err != nil {
			rlog.Error(err.Error())
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		rlog.Info(string(body))

		var m dswitch.SwitchConsumptions
		err = json.Unmarshal(body, &m)

		rlog.Info(m.TotalPower)

		messageVariable = strconv.Itoa(m.TotalPower)
		message = "page1.t4.txt=\"" + messageVariable + "\"" + "000"

		var messageToSend = []byte(message)

		// add terminaison
		messageToSend[len(messageToSend)-1] = 255
		messageToSend[len(messageToSend)-2] = 255
		messageToSend[len(messageToSend)-3] = 255

		fmt.Println("Opening serial port on " + serialPortNum)

		// OPENING (non blocking mode)
		c := &serial.Config{Name: serialPortNum, Baud: serialPortSpeed, ReadTimeout: time.Second}

		s, err := serial.OpenPort(c)

		if err != nil {
			log.Fatal(err)
		}

		rlog.Info("sending message")
		rlog.Info(messageToSend)

		n, err := s.Write(messageToSend)

		// RECEIVING (non blocking mode)
		buf := make([]byte, 128)
		n, _ = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		rlog.Info("receiving:")
		rlog.Infof("%q", buf[n])

		//time.Sleep(time.Millisecond * 10)

		// CLOSING
		s.Close()
		rlog.Info("Closing serial port " + serialPortNum)

		// Wait forever
		//for {
		//	time.Sleep(1 * time.Second)
		//}

		time.Sleep(1 * time.Second)

	}
}
