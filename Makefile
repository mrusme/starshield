.PHONY: build

build:
	arduino-cli compile --fqbn esp32:esp32:sparkfun_esp32c6_qwiic_pocket --libraries ./lib --clean -e ./

upload:
	arduino-cli upload -p /dev/ttyACM1 --fqbn esp32:esp32:sparkfun_esp32c6_qwiic_pocket --input-dir ./build/esp32.esp32.sparkfun_esp32c6_qwiic_pocket

connect:
	picocom -b 115200 -r /dev/ttyACM1
