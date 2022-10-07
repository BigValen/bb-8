#include <SoftwareSerial.h>

int
    LED_1F_1 = 4, // Flash the four 1F LEDs in sequence
    LED_1F_2 = 5,
    LED_1F_3 = 7,
    LED_1F_4 = 8,

    PWM_2R_1 = 16,  // Blue Pulse
    PWM_2R_2 = 26,  // Red Pulse

    PWM_3B = 14,   // Blue pulse
    
    PWM_4T_1  = 0,  // Blue pulse
    PWM_4T_2  = 1,  // Red pulse

    PWM_5L = 27,   // blue pulse

    LED_6BK_1 = 9,  // yellow
    PWM_6BK_2 = 15, // red pulse
    LED_6BK_3 = 10, // red solid on
    PWM_6BK_4 = 25, // blue
    PWM_6BK_5 = 24, // red pulse
    LED_6BK_6 = 11, // red solid on

    LED_spare1 = 12,
    LED_spare2 = 13,
    CMD_RX = 18,
    CMD_TX = 19;
;



int brightness = 0;    // how bright the LED is
int fadeAmount = 5;    // how many points to fade the LED by

// the setup routine runs once when you press reset:
void setup() {
  SoftwareSerial command(CMD_RX, CMD_TX); // RX, TX
  Serial.begin(57600;)
  command.begin(9600);
  while (!Serial) { ; }


  command.println("ready.");
  Serial.println("ready.");

// Setup solid LEDs for output, PWM doesn't need that
  pinMode(LED_1F_1, OUTPUT);
  pinMode(LED_1F_2, OUTPUT);
  pinMode(LED_1F_3, OUTPUT);
  pinMode(LED_1F_4, OUTPUT);
  pinMode(LED_6BK_1, OUTPUT);
  pinMode(LED_6BK_3, OUTPUT);
  pinMode(LED_6BK_6, OUTPUT);

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
