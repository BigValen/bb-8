package main

import (
	"machine"
	"time"
)

type gpioLed struct {
	name     string
	pin      machine.Pin
	duration time.Duration
}

type pwmLed struct {
	PWM
	name     string
	solid    bool
	index    uint8
	pin      machine.Pin
	duration time.Duration
}

var (
	gpioLeds = [...]gpioLed{
		{name: "ONBOARD", pin: machine.LED, duration: 2 * time.Second},      // not used for PWM - 0 is the UART
		{name: "BK_1_YELLOW", pin: machine.GP16, duration: 3 * time.Second}, // not used for PWM - 0 is the UART
	}

	pwmLeds = [...]pwmLed{
		// Front battery indicator lights - on solid, unless battery drops
		{name: "F_1_BLUE", pin: machine.GP2, duration: 3 * time.Second},
		{name: "F_2_BLUE", pin: machine.GP3, duration: 4 * time.Second},
		{name: "F_3_BLUE", pin: machine.GP4, duration: 5 * time.Second},
		{name: "F_4_BLUE", pin: machine.GP5, duration: 6 * time.Second},

		{name: "B_BLUE", pin: machine.GP6, duration: 7 * time.Second},

		{name: "R_1_BLUE", pin: machine.GP7, duration: 8 * time.Second},
		{name: "R_2_RED", pin: machine.GP8, duration: 9 * time.Second},

		{name: "T_1_BLUE", pin: machine.GP9, duration: 2 * time.Second},
		{name: "T_2_RED", pin: machine.GP10, duration: 10 * time.Second},

		{name: "BK_2_RED", pin: machine.GP11, duration: 2 * time.Second},
		{name: "BK_3_RED", pin: machine.GP12, duration: 5 * time.Second},
		{name: "BK_4_BLUE", pin: machine.GP13, duration: 6 * time.Second},
		{name: "BK_5_RED", pin: machine.GP14, duration: 9 * time.Second},
		{name: "BK_6_RED", pin: machine.GP15, duration: 3 * time.Second},
	}
)

const numPwmLeds = len(pwmLeds)
const pwmPeriod uint64 = 1e9 / 200
const delayBetweenFades = time.Millisecond * 300

var pwmHandles [numPwmLeds]pwmHandle

func main() {

	println("Setting up serial")
	machine.UART0.Configure(machine.UARTConfig{BaudRate: 9600, TX: machine.GP0, RX: machine.GP1})
	time.Sleep(2 * time.Second)

	println("Ready")

	var err error

	for _, led := range gpioLeds {
		println("Kicked off blink on ", led.name)
		go blink(led)
	}

	for i, led := range pwmLeds {
		led.index = uint8(i)
		pwmLeds[i] = led
		pwmHandles[i].PWM, pwmHandles[i].channel, err = GetPWM(led.pin)
		if err != nil {
			panic(err)
		}
		err = pwmHandles[i].Configure(machine.PWMConfig{
			Period: pwmPeriod,
		})
		if err != nil {
			panic(err)
		}
		go fade(led)
	}
	go readUart(machine.UART0)
	for {
		time.Sleep(time.Second / 100)
	}
}

// blink the LED with given duration
func blink(led gpioLed) {
	println("blinking ", led.name)
	led.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.pin.Low()
		time.Sleep(led.duration)
		led.pin.High()
		time.Sleep(led.duration)
	}
}

// fade an LED with a small 'empty' delay
func fade(led pwmLed) {
	delay := (led.duration - 2*delayBetweenFades) / (2 * 100)
	println("fading ", led.name, " over ", led.duration/time.Second, "s as ", led.index)
	for {
		for percent := 0; percent < 100; percent++ {
			pwmHandles[led.index].SetPercent(uint8(percent))
			time.Sleep(delay)
		}
		time.Sleep(delayBetweenFades)
		for percent := 100; percent > 0; percent-- {
			pwmHandles[led.index].SetPercent(uint8(percent))
			time.Sleep(delay)
		}
		time.Sleep(delayBetweenFades)
	}
}

// Read a command from the UART
// TODO: Make this actually work
func readUart(uart *machine.UART) {
	var buffer []byte
	for {
		n, err := uart.Read(buffer)
		if err != nil {
			println("Got error reading serial:", err)
			panic(err)
		}
		if n > 0 {
			println("Read ", n, "bytes ", buffer)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
