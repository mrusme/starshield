#include <bsec2.h>
#include <Wire.h>
#include <Adafruit_SI1145.h>
#include <SparkFun_u-blox_GNSS_Arduino_Library.h>
#include <SparkFun_BNO08x_Arduino_Library.h>
#include <ArduinoJson.h>

#define PANIC_LED   LED_BUILTIN
#define ERROR_DUR   1000

// #define BNO08X_INT  A4
#define BNO08X_INT  -1
// #define BNO08X_RST  A5
#define BNO08X_RST  -1
#define BNO08X_ADDR 0x4B

#define VB_ERR  "err"  // Error
#define VB_ERRC "errc" // Error code
#define VB_IAQ  "iaq"  // IAQ
#define VB_IAQA "iaqa" // IAQ accuracy
#define VB_TEMP "temp" // Temperature
#define VB_PRES "pres" // Pressure
#define VB_HUMY "humy" // Humidity
#define VB_GASR "gasr" // Gas resistance
#define VB_STBS "stbs" // Stanilization status
#define VB_RUNS "runs" // Run in status
#define VB_VISB "visb" // Visible light
#define VB_INFR "infr" // Infrared
#define VB_ULVI "ulvi" // UV index
#define VB_GLAT "glat" // GPS Latitude (deg * 10^7)
#define VB_GLON "glon" // GPS Longitude (deg * 10^-7)
#define VB_GALT "galt" // GPS Altitude (mm)
#define VB_GALM "galm" // GPS Altitude MSL (mm)
#define VB_GRSP "grsp" // GPS Ground speed (mm/s)
#define VB_GHED "ghed" // GPS Heading (deg * 10^-5)
#define VB_GSIV "gsiv" // GPS SIV
#define VB_GUXT "guxt" // GPS UNIX Epoch Time
#define VB_IMSC "imsc" // IMU stability classifier
#define VB_IMAC "imac" // IMU activity classifier
#define VB_IMSP "imsp" // IMU step counter


JsonDocument valbuf;

Bsec2 envSensor;
Adafruit_SI1145 si1145;
SFE_UBLOX_GNSS gnss;
BNO08x imu;

void err(void);

void setupBsec(void);
void setupSI(void);
void setupGNSS(void);
void setupIMU(void);

void bsecCheckStatus(Bsec2 bsec);
void bsecDataCallback(const bme68xData data, const bsecOutputs outputs, Bsec2 bsec);

void siRefresh(void);
void gnssRefresh(void);
void imuRefresh(void);

void setup(void) {
  Serial.begin(115200);
  Wire.begin();
  pinMode(PANIC_LED, OUTPUT);

  setupBsec();
  setupSI();
  setupGNSS();
  // setupIMU();
}

void setupBsec(void) {
  bsecSensor sensorList[] = {
    BSEC_OUTPUT_IAQ,
    BSEC_OUTPUT_RAW_TEMPERATURE,
    BSEC_OUTPUT_RAW_PRESSURE,
    BSEC_OUTPUT_RAW_HUMIDITY,
    BSEC_OUTPUT_RAW_GAS,
    BSEC_OUTPUT_STABILIZATION_STATUS,
    BSEC_OUTPUT_RUN_IN_STATUS
  };

  if(!envSensor.begin(BME68X_I2C_ADDR_HIGH, Wire)) {
    bsecCheckStatus(envSensor);
  }

  if(!envSensor.updateSubscription(sensorList,
        ARRAY_LEN(sensorList), BSEC_SAMPLE_RATE_ULP)) {
    bsecCheckStatus(envSensor);
  }

  envSensor.attachCallback(bsecDataCallback);
}

void setupSI(void) {
  si1145  = Adafruit_SI1145();
  if(!si1145.begin()) {
    valbuf[VB_ERR] = F("SI1145 not found");
    valbuf[VB_ERRC] = 1;
    err();
  }
}

void setupGNSS(void) {
  if(!gnss.begin()) {
    valbuf[VB_ERR] = F("GNSS not found");
    valbuf[VB_ERRC] = 1;
    err();
  }

  gnss.setI2COutput(COM_TYPE_UBX);
  gnss.saveConfigSelective(VAL_CFG_SUBSEC_IOPORT);
}

void setupIMU(void) {
  if(!imu.begin(BNO08X_ADDR, Wire, BNO08X_INT, BNO08X_RST)) {
    valbuf[VB_ERR] = F("IMU not found");
    valbuf[VB_ERRC] = 1;
    err();
  }
}

void loop(void) {
  if(!envSensor.run()) {
    bsecCheckStatus(envSensor);
  }
  siRefresh();
  gnssRefresh();
  // imuRefresh();

  serializeJson(valbuf, Serial);
  delay(1000);
}

void err(void) {
  while(1) {
    digitalWrite(PANIC_LED, HIGH);
    delay(ERROR_DUR);
    digitalWrite(PANIC_LED, LOW);
    delay(ERROR_DUR);
    serializeJson(valbuf, Serial);
  }
}

void bsecDataCallback(const bme68xData data, const bsecOutputs outputs, Bsec2 bsec) {
  if(!outputs.nOutputs) {
    return;
  }

  for(uint8_t i = 0; i < outputs.nOutputs; i++) {
    const bsecData output  = outputs.output[i];
    switch(output.sensor_id) {
      case BSEC_OUTPUT_IAQ:
        valbuf[VB_IAQ] = output.signal;
        valbuf[VB_IAQA] = output.accuracy;
        break;
      case BSEC_OUTPUT_RAW_TEMPERATURE:
        valbuf[VB_TEMP] = output.signal;
        break;
      case BSEC_OUTPUT_RAW_PRESSURE:
        valbuf[VB_PRES] = output.signal;
        break;
      case BSEC_OUTPUT_RAW_HUMIDITY:
        valbuf[VB_HUMY] = output.signal;
        break;
      case BSEC_OUTPUT_RAW_GAS:
        valbuf[VB_GASR] = output.signal;
        break;
      case BSEC_OUTPUT_STABILIZATION_STATUS:
        valbuf[VB_STBS] = output.signal;
        break;
      case BSEC_OUTPUT_RUN_IN_STATUS:
        valbuf[VB_RUNS] = output.signal;
        break;
      default:
        break;
    }
  }
}

void bsecCheckStatus(Bsec2 bsec) {
  if(bsec.status < BSEC_OK) {
    valbuf[VB_ERR] = F("BSEC error");
    valbuf[VB_ERRC] = bsec.status;
    err();
  } else if(bsec.status > BSEC_OK) {
    valbuf[VB_ERR] = F("BSEC warning");
    valbuf[VB_ERRC] = bsec.status;
  }

  if(bsec.sensor.status < BME68X_OK) {
    valbuf[VB_ERR] = F("BSEC sensor error");
    valbuf[VB_ERRC] = bsec.sensor.status;
    err();
  } else if(bsec.sensor.status > BME68X_OK) {
    valbuf[VB_ERR] = F("BSEC sensor warning");
    valbuf[VB_ERRC] = bsec.sensor.status;
  }
}

void siRefresh(void) {
  valbuf[VB_VISB] = si1145.readVisible();
  valbuf[VB_INFR] = si1145.readIR();
  valbuf[VB_ULVI] = (si1145.readUV() / 100.0);
}

void gnssRefresh(void) {
  valbuf[VB_GLAT] = gnss.getLatitude();
  valbuf[VB_GLON] = gnss.getLongitude();
  valbuf[VB_GALT] = gnss.getAltitude();
  valbuf[VB_GALM] = gnss.getAltitudeMSL();
  valbuf[VB_GRSP] = gnss.getGroundSpeed();
  valbuf[VB_GHED] = gnss.getHeading();
  valbuf[VB_GUXT] = gnss.getUnixEpoch();
  valbuf[VB_GSIV] = gnss.getSIV();
}

void imuRefresh(void) {
  if(imu.wasReset()) {
    if(!imu.enableStabilityClassifier()) {
      valbuf[VB_ERR] = F("IMU stability classifier not available");
      valbuf[VB_ERRC] = 1;
      err();
    }
  }

  if(imu.getSensorEvent()) {
    switch(imu.getSensorEventID()) {
      case SENSOR_REPORTID_STABILITY_CLASSIFIER:
        valbuf[VB_IMSC] = imu.getStabilityClassifier();
        break;
      case SENSOR_REPORTID_PERSONAL_ACTIVITY_CLASSIFIER:
        valbuf[VB_IMAC] = imu.getActivityClassifier();
        break;
      case SENSOR_REPORTID_STEP_COUNTER:
        valbuf[VB_IMSP] = imu.getStepCount();
        break;
      default:
        break;
    }
  }
}
