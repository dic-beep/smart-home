#include <Arduino.h>
#if defined(ESP32)
  #include <WiFi.h>
#elif defined(ESP8266)
  #include <ESP8266WiFi.h>
#endif
#include <Firebase_ESP_Client.h>

//Provide the token generation process info.
#include "addons/TokenHelper.h"
//Provide the RTDB payload printing info and other helper functions.
#include "addons/RTDBHelper.h"

// Insert your network credentials
#define WIFI_SSID "POCO M3"
#define WIFI_PASSWORD "cencen123"

// Insert Firebase project API Key
#define API_KEY "AIzaSyD1mKj4AnEfcJyJQjMp5QLarrcLxIgJZUY"

// Insert RTDB URLefine the RTDB URL */
#define DATABASE_URL "https://smart-home-bdf74-default-rtdb.firebaseio.com/" 

//Define Firebase Data object
FirebaseData fbdo;
FirebaseAuth auth;
FirebaseConfig config;

//Define led pin
#define LED_PIN 2

unsigned long sendDataPrevMillis = 0;
bool signupOK = false;
bool ledStatus = false;

void setup(){
  Serial.begin(115200);
  //setup led as an output
  pinMode(LED_PIN, OUTPUT);

  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  Serial.print("Connecting to Wi-Fi");
  while (WiFi.status() != WL_CONNECTED){
    Serial.print(".");
    delay(300);
  }
  Serial.println();
  Serial.print("Connected with IP: ");
  Serial.println(WiFi.localIP());
  Serial.println();

  /* Assign the api key (required) */
  config.api_key = API_KEY;

  /* Assign the RTDB URL (required) */
  config.database_url = DATABASE_URL;

  /* Sign up */
  if (Firebase.signUp(&config, &auth, "", "")){
    Serial.println("ok");
    signupOK = true;
  }
  else{
    Serial.printf("%s\n", config.signer.signupError.message.c_str());
  }

  /* Assign the callback function for the long running token generation task */
  config.token_status_callback = tokenStatusCallback; //see addons/TokenHelper.h
  
  Firebase.begin(&config, &auth);
  Firebase.reconnectWiFi(true);
}

void loop(){
  if (Firebase.ready() && signupOK && (millis() - sendDataPrevMillis > 15000 || sendDataPrevMillis == 0)){
    sendDataPrevMillis = millis();

  }
    if(Firebase.RTDB.getBool(&fbdo, "/LED/digital")){
      if(fbdo.dataType() == "boolean"){
        ledStatus = fbdo.boolData();
        digitalWrite(LED_PIN, (int)ledStatus);
        Serial.println("Succsess full read from: "+ fbdo.dataPath() + ": " + (int)ledStatus + "("+fbdo.dataType()+")");
      }else{
        Serial.println("Failed: " + fbdo.errorReason());
      }
    }
}

  