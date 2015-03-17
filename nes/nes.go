package nes

import (
	"image"
	"image/color"
	"log"
)

type NES struct {
	CPU         *CPU
	APU         *APU
	PPU         *PPU
	Cartridge   *Cartridge
	Controller1 *Controller
	Controller2 *Controller
	Mapper      Mapper
	RAM         []byte
}

func NewNES(path string) (*NES, error) {
	cartridge, err := LoadNESFile(path)
	if err != nil {
		return nil, err
	}
	ram := make([]byte, 2048)
	controller1 := NewController()
	controller2 := NewController()
	nes := NES{nil, nil, nil, cartridge, controller1, controller2, nil, ram}
	nes.Mapper = NewMapper(&nes, cartridge)
	nes.CPU = NewCPU(&nes)
	nes.APU = NewAPU()
	nes.PPU = NewPPU(&nes)
	return &nes, nil
}

func (nes *NES) Step() int {
	cpuCycles := nes.CPU.Step()
	ppuCycles := cpuCycles * 3
	for i := 0; i < ppuCycles; i++ {
		nes.PPU.Step()
		nes.Mapper.Step()
	}
	return cpuCycles
}

func (nes *NES) StepFrame() {
	frame := nes.PPU.Frame
	for frame == nes.PPU.Frame {
		nes.Step()
	}
}

func (nes *NES) Buffer() *image.RGBA {
	return nes.PPU.buffer
}

func (nes *NES) BackgroundColor() color.RGBA {
	return palette[nes.PPU.readPalette(0)%64]
}

func (nes *NES) SetPressed(controller, button int, pressed bool) {
	switch controller {
	case 1:
		nes.Controller1.SetPressed(button, pressed)
	case 2:
		nes.Controller2.SetPressed(button, pressed)
	default:
		log.Fatalf("unhandled controller press: %d", controller)
	}
}
