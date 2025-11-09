package main

import (
	"io"
	"log"
	"os"

	"machine"

	"github.com/nobonobo/gamepad-emulator/service"
)

const (
	LED1 machine.Pin = 25
	LED2 machine.Pin = 14
	LED3 machine.Pin = 15
	SW1  machine.Pin = 24
	SW2  machine.Pin = 23
	SW3  machine.Pin = 22
)

type w struct {
	io.Reader
	io.WriteCloser
}

func init() {
	LED1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LED2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LED3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LED1.High()
	LED2.High()
	LED3.High()
	SW1.Configure(machine.PinConfig{Mode: machine.PinInput})
	SW2.Configure(machine.PinConfig{Mode: machine.PinInput})
	SW3.Configure(machine.PinConfig{Mode: machine.PinInput})
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	srv := service.New()
	if err := srv.Run(w{Reader: os.Stdin, WriteCloser: os.Stdout}); err != nil {
		log.Fatal(err)
	}
}
