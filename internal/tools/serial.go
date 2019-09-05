package tools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/romana/rlog"
	"github.com/tarm/serial"
)

var (
	serialPortNum   = "/dev/ttyAMA0"
	serialPortSpeed = 9600
)

func WriteTotalPower(value int) error {
	return WriteScreen("page1.t4.txt", strconv.Itoa(value))
}

func WriteScreen(field string, message string) error {
	message = field + "=\"" + message + "\"" + "000"

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
		return err
	}

	rlog.Info("sending message")
	rlog.Info(messageToSend)

	_, err = s.Write(messageToSend)

	// RECEIVING (non blocking mode)
	buf := make([]byte, 128)
	n, _ := s.Read(buf)

	rlog.Info("receiving:")
	rlog.Infof("%q", buf[n])

	// CLOSING
	s.Close()
	rlog.Info("Closing serial port " + serialPortNum)

	return err
}
