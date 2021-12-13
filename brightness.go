package main

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const HIGH_BRIGHTNESS = 750
const LOW_BRIGHTNESS = 250
const MINIMUM_BRIGHTNESS = 1
const BRIGHTNESS_FILE = "/sys/class/backlight/gmux_backlight/brightness"

var userId = nativeUserId
var groupId = nativeGroupId
var fileReader = nativeFileReader
var filePermissionCheck = nativeFilePermissionCheck
var fileWriter = nativeFileWriter

func main() {
	setBrightness()
}

func setBrightness() bool {
	if !isBrightnessValid() {
		fmt.Println("System does not report a current brightness level.")
		return false
	}

	if !canWriteBrightness() {
		fmt.Println("Cannot update the brightness level; you may need to run with elevated privileges.")
		return false
	}

	if isHighBrightness() {
		return setLowBrightness() > 0
	} else {
		return setHighBrightness() > 0
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

	if value > 0 {
		currentBrightness = value
	}

	return currentBrightness
}

func canWriteBrightness() bool {
	return fileWriteBit() == "w"
}

func fileWriteBit() string {
	uid := userId()
	gid := groupId()
	fmt.Printf("Running as: UID·%v GID·%v\n", uid, gid)

	fileUid, fileGid, mode := filePermissionCheck(BRIGHTNESS_FILE)
	filePerms := string(mode.String())
	fmt.Printf("File ownership: UID·%v GID·%v %v\n", fileUid, fileGid, filePerms)

	return fileWriteBitString(uid, gid, fileUid, fileGid, string(filePerms[2]), string(filePerms[5]), string(filePerms[8]))
}

func fileWriteBitString(uid, gid, fileUid, fileGid int, userWritePerm, groupWritePerm, worldWritePerm string) string {
	if uid == fileUid {
		fmt.Printf("Write permission is: User %v\n", userWritePerm)
		return userWritePerm
	} else if gid == fileGid {
		fmt.Printf("Write permission is: Group %v\n", groupWritePerm)
		return groupWritePerm
	} else {
		fmt.Printf("Write permission is: World %v\n", worldWritePerm)
		return worldWritePerm
	}
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

func nativeUserId() int {
	return syscall.Getuid()
}
func nativeGroupId() int {
	return syscall.Getgid()
}

func nativeFileReader(f string) string {
	content, _ := os.ReadFile(f)
	value := string(content)
	return value
}

func nativeFilePermissionCheck(f string) (int, int, fs.FileMode) {
	fmt.Printf("File: %s\n", BRIGHTNESS_FILE)
	info, err := os.Stat(f)

	if err != nil {
		panic(err)
	}

	stat := info.Sys().(*syscall.Stat_t)

	return int(stat.Uid), int(stat.Gid), info.Mode()
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
