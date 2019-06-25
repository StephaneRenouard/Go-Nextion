package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/romana/rlog"
	"github.com/tarm/serial"
)

type arrayString []string

func (i *arrayString) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayString) Set(value string) error {
	if value == "" {
		return nil
	}

	var list []string
	for _, in := range strings.Split(value, ",") {
		list = append(list, in)
	}

	*i = arrayString(list)
	return nil
}

func (i *arrayString) Get() interface{} { return []string(*i) }

type arrayInt []int

func (i *arrayInt) Set(val string) error {
	if val == "" {
		return nil
	}

	var list []int
	for _, in := range strings.Split(val, ",") {
		i, err := strconv.Atoi(in)
		if err != nil {
			return err
		}

		list = append(list, i)
	}

	*i = arrayInt(list)
	return nil
}

func (i *arrayInt) Get() interface{} { return []int(*i) }

func (i *arrayInt) String() string {
	var list []string
	for _, in := range *i {
		list = append(list, strconv.Itoa(in))
	}
	return strings.Join(list, ",")
}

func main() {

	var serialPortNum = "/dev/ttyAMA0"
	var serialPortSpeed = 9600
	var messageVariable = "0"
	var message = "page1.t4.txt=\"" + messageVariable + "\"" + "000"
	var messageToSend = []byte(message)

	var serverURL = "192.168.0.2"
	var serverPORT = "8888"

	var logLevel = "debug"

	type Message struct {
		totalPower     int
		lightningPower int
		blindPower     int
		hvacPower      int
	}

	os.Setenv("RLOG_LOG_LEVEL", logLevel)
	os.Setenv("RLOG_LOG_NOTIME", "yes")
	rlog.UpdateEnv()

	for {

		requestBody, err := json.Marshal("totalPower")
		if err != nil {
			rlog.Error(err.Error())
			os.Exit(1)
		}

		url := "https://" + serverURL + ":" + serverPORT + "/v1.0/status/consumptions"

		req, _ := http.NewRequest("GET", url, bytes.NewBuffer(requestBody))
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

		fmt.Println(string(body))

		var m Message
		err = json.Unmarshal(body, &m)

		fmt.Println(m.totalPower)

		messageVariable = string(m.totalPower)
		message = "page1.t4.txt=\"" + messageVariable + "\"" + "000"

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

		fmt.Println("sending message")
		fmt.Printf("%s\n", messageToSend)

		n, err := s.Write(messageToSend)

		// RECEIVING (non blocking mode)
		buf := make([]byte, 128)
		n, _ = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("receiving:")
		log.Printf("%q", buf[n])

		//time.Sleep(time.Millisecond * 10)

		// CLOSING
		s.Close()
		fmt.Println("Closing serial port " + serialPortNum)

		// Wait forever
		//for {
		//	time.Sleep(1 * time.Second)
		//}

		time.Sleep(1 * time.Second)

	}
}
