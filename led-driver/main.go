package main

import (
	"machine"
	"time"
)

var period uint64 = 1e9 / 500

func main() {

	println("Ready")
	// Configure the PWM with the given period.
	machine.PWM4.Configure(machine.PWMConfig{
		Period: period,
	})
	machine.PWM1.Configure(machine.PWMConfig{
		Period: period * 2,
	})

	LedInt, err := machine.PWM4.Channel(machine.LED)
	if err != nil {
		println(err.Error())
		return
	}
	Led2, err := machine.PWM1.Channel(machine.GP2)
	if err != nil {
		println(err.Error())
		return
	}
	Led3, err := machine.PWM1.Channel(machine.GP3)
	if err != nil {
		println(err.Error())
		return
	}

	for {
		for i := 1; i < 255; i++ {
			// This performs a stylish fade-out blink
			machine.PWM4.Set(LedInt, machine.PWM4.Top()/uint32(i))

			machine.PWM1.Set(Led2, machine.PWM1.Top()-machine.PWM1.Top()/uint32(i))
			machine.PWM1.Set(Led3, machine.PWM1.Top()/uint32(i))
			time.Sleep(time.Millisecond * 5)
		}
		println("loop")
	}
}
