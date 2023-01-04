package main

import (
	"fmt"
	"machine"
	"strconv"
	"strings"
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
		// Front battery indicator lights - on solid, unless battery drops - don't use the first two, during debugging
		//{name: "F_1_BLUE", pin: machine.GP0, duration: 3 * time.Second, solid: true},
		//{name: "F_2_BLUE", pin: machine.GP1, duration: 4 * time.Secon, solid: true},
		{name: "F_3_BLUE", pin: machine.GP2, duration: 5 * time.Second, solid: true},
		{name: "F_4_BLUE", pin: machine.GP3, duration: 6 * time.Second, solid: true},

		{name: "B_BLUE", pin: machine.GP4, duration: 7 * time.Second, solid: false},

		{name: "R_1_BLUE", pin: machine.GP5, duration: 8 * time.Second, solid: false},
		{name: "R_2_RED", pin: machine.GP6, duration: 9 * time.Second, solid: false},

		{name: "T_1_BLUE", pin: machine.GP7, duration: 2 * time.Second, solid: false},
		{name: "T_2_RED", pin: machine.GP10, duration: 10 * time.Second, solid: false},

		{name: "BK_2_RED", pin: machine.GP11, duration: 2 * time.Second, solid: false},
		{name: "BK_3_RED", pin: machine.GP12, duration: 5 * time.Second, solid: false},
		{name: "BK_4_BLUE", pin: machine.GP13, duration: 6 * time.Second, solid: false},
		{name: "BK_5_RED", pin: machine.GP14, duration: 9 * time.Second, solid: false},
		{name: "BK_6_RED", pin: machine.GP15, duration: 3 * time.Second, solid: false},
	}
)

const numPwmLeds = len(pwmLeds)
const pwmPeriod uint64 = 1e9 / 200
const delayBetweenFades = time.Millisecond * 300

var pwmHandles [numPwmLeds]pwmHandle

func main() {

	println("Setting up serial")
	machine.UART1.Configure(machine.UARTConfig{BaudRate: 9600, TX: machine.GP8, RX: machine.GP9})
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
		if led.solid != true {
			go fade(led)
		} else {
			solid(led, 100)
		}
	}
	go readUart(machine.UART1)
	for {
		time.Sleep(time.Second / 100)
	}
}

// Set the LED to a specific level
func solid(led pwmLed, level uint8) {
	if level < 0 || level > 100 {
		panic("led level out of bounds")
	}
	// If level is 0 or 100, then don't use PWM, so it doesn't interfere with any lights on the same block
	if level == 0 {
		led.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		led.pin.Low()
	} else if level == 100 {
		led.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		led.pin.High()
	} else {
		err := pwmHandles[i].Configure(machine.PWMConfig{
			Period: pwmPeriod,
		})
		if err != nil {
			panic(err)
		}
		pwmHandles[led.index].SetPercent(uint8(level))
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
	var input string
	for {
		b, err := uart.ReadByte()
		if err != nil {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if b == 13 {
			if len(input) < 1 {
				uart.Write([]byte("\n\r> "))
				continue
			}
			command := strings.Fields(input)
			switch command[0] {
			case "battery":
				if len(command) != 2 {
					uart.Write([]byte("\n\rerror: battery <0-100>"))
				} else {
					level, err := strconv.Atoi(command[1])
					if err != nil || level < 0 || level > 100 {
						uart.Write([]byte("\n\rerror: battery <0-100>"))
					} else {
						uart.Write([]byte("set battery level: " + fmt.Sprint(level) + "\n\r"))
					}
				}
			default:
				uart.Write([]byte("unknown command: " + input + "\n\r"))
			}
			input = ""
			uart.Write([]byte("\n\r> "))
		} else {
			uart.Write([]byte{b})
			input += string(b)
		}
		//out := "Read " + fmt.Sprint(b)
		//uart.Write([]byte(out))
		time.Sleep(20 * time.Millisecond)

		/*
			var buffer []byte
			n, err := uart.Read(buffer)
			if err != nil {
				out := string("Got error reading serial:" + string(err.Error()))
				uart.Write([]byte(out))
				panic(err)
			}
			if n > 0 {
				out := "Read " + fmt.Sprint(n) + "bytes " + string(buffer)
				uart.Write([]byte(out))
			} else {
				uart.Write([]byte(fmt.Sprint(n) + "."))
			}
		*/

	}
}
