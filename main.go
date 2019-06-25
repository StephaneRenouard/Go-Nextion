package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {

	var serialPortNum = "COM9"
	var serialPortSpeed = 9600
	var messageVariable = "33"
	var message = "page1.t4.txt=\"" + messageVariable + "\"" + "000"
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
}
