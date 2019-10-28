package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

var speed = 9600
var port = "/dev/ttyAMA0"
var duration = time.Millisecond * 135

/*
Baudrate	2400	4800	9600	19200	38400	57600	115200
Delay(ms)	447		239		135		83		57		48		39
*/

func main() {

	var capsulePath string

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.StringVar(&capsulePath, "input", "", "Path to the update Screen capsule")
	flag.StringVar(&capsulePath, "i", "", "Path to the update Screen capsule")
	flag.Parse()

	if _, err := os.Stat(capsulePath); os.IsNotExist(err) {
		log.Fatal("Screen capsule " + capsulePath + " doesn't exist")
	}

	c := &serial.Config{Name: port, Baud: speed, ReadTimeout: time.Millisecond * 5000}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	n, err := s.Write(generateInstruction(""))
	time.Sleep(duration)
	n, err = s.Write(generateInstruction("DRAKJHSUYDGBNCJHGJKSHBDN"))
	time.Sleep(duration)
	n, err = s.Write(generateInstruction("connect"))
	time.Sleep(duration)
	n, err = s.Write(generateInstruction("ÿÿconnect"))
	time.Sleep(duration)

	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", buf[:n])

	// get file size
	fi, err := os.Stat(capsulePath)
	if err != nil {
		fmt.Println("Erreur reading file")
	}
	size := fi.Size()
	fmt.Println("file size is ", size)

	chunkNb := int(size) / 4096

	// send the upload firmware command
	n, err = s.Write(generateInstruction("whmi-wri " + strconv.FormatInt(size, 10) + "," + strconv.Itoa(speed) + ",res0"))
	time.Sleep(time.Millisecond * 200)
	buf = make([]byte, 128)
	n, err = s.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", buf[:n])

	// read file
	f, err := os.Open(capsulePath)
	check(err)
	for i := 0; i < chunkNb; i++ {
		offset := int64(i * 4096)
		fmt.Println("processing chunk ", i, " on ", chunkNb, " starting at offset ", offset)

		chunk := make([]byte, 4096)
		_, err = f.Read(chunk)
		check(err)
		fmt.Println("control ", chunk[0], chunk[1], chunk[2], chunk[3], chunk[4], chunk[5], chunk[6], chunk[7], chunk[8], chunk[9])
		n, err = s.Write(chunk)
		time.Sleep(time.Millisecond * 200)
		buf = make([]byte, 128)
		n, err = s.Read(buf)
		for err != nil {
			fmt.Println("Wait for acknowledgment")
			time.Sleep(time.Millisecond * 200)
			buf = make([]byte, 128)
			n, err = s.Read(buf)
		}
		fmt.Printf("%q\n", buf[:n])

	}

	// last chunk
	offset := int64(chunkNb * 4096)
	fmt.Println("processing last chunk starting at offset ", offset)
	chunk := make([]byte, (size - offset))
	fmt.Println("last chunk size is ", size-offset, " bytes, starting at offset ", offset, " on ", size)
	_, err = f.Read(chunk)
	check(err)
	fmt.Println("control ", chunk[0], chunk[1], chunk[2], chunk[3], chunk[4], chunk[5], chunk[6], chunk[7], chunk[8], chunk[9])
	n, err = s.Write(chunk)
	time.Sleep(time.Millisecond * 200)
	fmt.Println("end of update process, ", size, " bytes written")

	buf = make([]byte, 128)
	n, err = s.Read(buf)
	for err != nil {
		fmt.Println("Wait for acknowledgment")
		time.Sleep(time.Millisecond * 200)
		buf = make([]byte, 128)
		n, err = s.Read(buf)
	}
	fmt.Printf("%q\n", buf[:n])

}

func generateInstruction(inst string) []byte {
	stringSlice := []string{inst, "fff"}
	stringByte := strings.Join(stringSlice, "")
	instruction := []byte(stringByte)
	instruction[len(instruction)-1] = 0xff
	instruction[len(instruction)-2] = 0xff
	instruction[len(instruction)-3] = 0xff
	fmt.Println(instruction)
	return instruction

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
