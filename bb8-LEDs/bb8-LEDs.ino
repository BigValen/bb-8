
int PWM_5L = 27,
    PWM_4TA  = 0,
    PWM_4TB  = 1,
    PWM_3B = 14,
    PWM_6BK2 = 15,
    PWM_5 = 16,
    PWM_6 = 26,
    PWM_7 = 25,
    PWM_8 = 24,
    LED_1FA = 4, // Flash the four 1F LEDs in sequence
    LED_1FB = 5,
    LED_1FC = 7,
    LED_1FD = 8,
    LED_6BK2 = 9,
    LED = 10,
    LED = 11,
    LED = 12,
    LED = 13,


int brightness = 0;    // how bright the LED is
int fadeAmount = 5;    // how many points to fade the LED by

// the setup routine runs once when you press reset:
void setup() {
  // declare pin 9 to be an output:
  pinMode(led, OUTPUT);
}

// the loop routine runs over and over again forever:
void loop() {

  while (Serial.available()) {
    delay(3);  //delay to allow buffer to fill 
    if (Serial.available() >0) {
      char c = Serial.read();  //gets one byte from serial buffer
      readString += c; //makes the string readString
    } 
  }

  if (readString.length() >0) {
      Serial.println(readString); //see what was received
      
      // expect a string like 07002100 containing the two servo positions      
      servo1 = readString.substring(0, 4); //get the first four characters
      servo2 = readString.substring(4, 8); //get the next four characters 
      




  // set the brightness of pin 9:
  analogWrite(led, brightness);

  // change the brightness for next time through the loop:
  brightness = brightness + fadeAmount;

  // reverse the direction of the fading at the ends of the fade:
  if (brightness <= 0 || brightness >= 255) {
    fadeAmount = -fadeAmount;
  }
  // wait for 30 milliseconds to see the dimming effect
  delay(30);
}
