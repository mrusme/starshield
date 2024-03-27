.PHONY: build

build:
	arduino-cli compile --fqbn esp32:esp32:sparkfun_esp32s2_thing_plus --clean -e ./

upload:
	arduino-cli upload -p /dev/ttyUSB0 --fqbn esp32:esp32:sparkfun_esp32s2_thing_plus --input-dir ./build/esp32.esp32.sparkfun_esp32s2_thing_plus

connect:
	picocom -b 115200 -x 35000 -r /dev/ttyUSB0
