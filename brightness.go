package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const HIGH_BRIGHTNESS = 750
const LOW_BRIGHTNESS = 250
const MINIMUM_BRIGHTNESS = 1
const BRIGHTNESS_FILE = "/sys/class/backlight/gmux_backlight/brightness"

var fileReader = nativeFileReader
var fileWriter = nativeFileWriter

func main() {
	setBrightness()
}

func setBrightness() {
	if !isBrightnessValid() {
		return
	}

	if isHighBrightness() {
		setLowBrightness()
	} else {
		setHighBrightness()
	}
}

/* Test coverage starts here! */

func isBrightnessValid() bool {
	return currentBrightness() >= MINIMUM_BRIGHTNESS
}

func isHighBrightness() bool {
	return currentBrightness() == HIGH_BRIGHTNESS
}

func currentBrightness() int {
	currentBrightness := -1

	contents := fileReader(BRIGHTNESS_FILE)
	contents = strings.TrimSuffix(contents, "\n")
	value, _ := strconv.Atoi(contents)

	if(value > 0) {
		currentBrightness = value
	}

	return currentBrightness
}

func setLowBrightness() int {
	fmt.Printf("·  Setting low brightness (%d)...\n", LOW_BRIGHTNESS)
	fileWriter(BRIGHTNESS_FILE, strconv.Itoa(LOW_BRIGHTNESS))

	return LOW_BRIGHTNESS
}

func setHighBrightness() int {
	fmt.Printf("·  Setting high brightness (%d)...\n", HIGH_BRIGHTNESS)
	fileWriter(BRIGHTNESS_FILE, strconv.Itoa(HIGH_BRIGHTNESS))

	return HIGH_BRIGHTNESS
}

/* Test coverage ends here! */

func nativeFileReader(f string) string {
	content, _ := os.ReadFile(f)
	value := string(content)
	return value
}

func nativeFileWriter(filename, str string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
        panic(err)
    }

    defer f.Close()

	if _, err = f.WriteString(str); err != nil {
        panic(err)
    }

	return err
}
