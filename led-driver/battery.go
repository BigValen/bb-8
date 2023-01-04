package main

// accept a battery level from 0 to 100, and set 4 lights based on it
func battery(level uint8) {
	if level < 0 || level > 100 {
		panic("battery level out of bounds")
	}

}
