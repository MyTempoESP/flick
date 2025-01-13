package flick

import (
	"bufio"
	"fmt"
	"time"

	"github.com/MyTempoESP/serial"
)

type SerialForth struct {
	config  *serial.Config
	port    *serial.Port
	scanner *bufio.Scanner

	serialReader *SerialReader
}

func (SerialForth) GetBytes(s string) (fixed string) {

	rns := []rune(s) // convert to rune
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		// swap the letters of the string,
		// like first with last and so on.
		rns[i], rns[j] = rns[j], rns[i]
	}

	for _, c := range rns {

		fixed = fmt.Sprintf("%s %d", fixed, c)
	}

	fixed = fmt.Sprintf("%s %d", fixed, len(s))

	return
}

func NewSerialForth(device string) (forth SerialForth, err error) {

	// Configure the serial port
	forth.config = &serial.Config{
		Name: device,
		Baud: 115200,
	}

	// Open the serial port
	forth.port, err = serial.OpenPort(forth.config)

	if err != nil {

		return
	}

	// Allow time for Arduino to reset
	time.Sleep(2 * time.Second)

	// defining the API functions
	err = forth.Run("BYE")
	err = forth.Run(": a 1 API ;")
	err = forth.Run(": m 2 API ;")
	err = forth.Run(": d 3 API ;")

	// shorthand for end of line
	err = forth.Run(": $ 15 ;")
	//err = forth.Run("1 TRC")

	if err != nil {

		return
	}

	forth.scanner = bufio.NewScanner(forth.port)
	err = forth.scanner.Err()

	if err != nil {

		return
	}

	forth.serialReader = NewSerialReader()
	forth.serialReader.Start(forth.scanner)

	return
}

func (forth *SerialForth) Close() {

	forth.serialReader.Stop()
	forth.port.Close()
}

func (forth *SerialForth) Query(msg string) (output string, err error) {

	err = forth.Run(msg)

	output = forth.serialReader.GetLatest()

	return
}

func (forth *SerialForth) Run(msg string) (err error) {

	_, err = forth.port.Write(append([]byte(msg), '\n'))

	time.Sleep(500 * time.Millisecond)

	return
}
