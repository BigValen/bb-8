package main

import (
	"machine"
)

// pwmHandle facilitates PWM usage:
//
//	pwms[i].Set(pwms[i].channel, value)
type pwmHandle struct {
	// Embedded type so methods can be used directly.
	PWM
	channel uint8
}

// SetPercent sets the PWM with a percentage of the max duty cycle.
// percent must be in range 0..100.
func (ph *pwmHandle) SetPercent(percent uint8) {
	top := ph.PWM.Top()
	ph.PWM.Set(ph.channel, uint32(percent)*top/100)
}

type PWM interface {
	Top() uint32
	Set(ch uint8, value uint32)
	Channel(pin machine.Pin) (uint8, error)
	Configure(machine.PWMConfig) error
}

func GetPWM(pin machine.Pin) (pwm PWM, channel uint8, err error) {
	slice, err := machine.PWMPeripheral(pin)
	if err != nil {
		return pwm, channel, err
	}
	pwm = pwmFromSlice(slice)
	channel, err = pwm.Channel(pin)
	if err != nil {
		return pwm, channel, err
	}
	return pwm, channel, nil
}

func pwmFromSlice(i uint8) PWM {
	if i > 7 {
		panic("PWM out of range")
	}
	pwms := [...]PWM{
		machine.PWM0, machine.PWM1, machine.PWM2,
		machine.PWM3, machine.PWM4, machine.PWM5,
		machine.PWM6, machine.PWM7,
	}
	return pwms[i]
}

// GetPWM acquires a unconfigured [PWM instance]. The returned PWM
// should be configured before use. e.g:
//
//	pwm, channel, err := GetPWM(machine.GP20)
//	if err != nil {
//		panic(err)
//	}
//	err = pwm.Configure(machine.PWMConfig{Period: 1e9 / 200}) // 200Hz
//	if err != nil {
//		panic(err) // On rp2040 only error is for bad Period.
//	}
//	pwm.Set(channel, pwm.Top()/4) // 25% duty cycle.
//
// [PWM instance]: https://tinygo.org/docs/tutorials/pwm/
