package main

import (
	"io/fs"
	"testing"
)

const testUserId = 500
const testGroupId = 500
const rootUserId = 0
const rootGroupId = 0

func TestIsBrightnessValid(t *testing.T) {
	t.Run("isBrightnessValid", func(t *testing.T) {
		t.Run("when the current brightness is above the minimum", func(t *testing.T) {
			t.Run("is true", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "2\n"
				}

				actual := isBrightnessValid()
				expected := true

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})

		t.Run("when the current brightness is at the minimum", func(t *testing.T) {
			t.Run("is true", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "1\n"
				}

				actual := isBrightnessValid()
				expected := true

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})

		t.Run("when the current brightness is below the minimum", func(t *testing.T) {
			t.Run("is false", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "0\n"
				}

				actual := isBrightnessValid()
				expected := false

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})

		t.Run("when the current brightness is negative", func(t *testing.T) {
			t.Run("is false", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "-1\n"
				}

				actual := isBrightnessValid()
				expected := false

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})
	})
}

func TestIsHighBrightness(t *testing.T) {
	t.Run("isHighBrightness", func(t *testing.T) {
		t.Run("when the current brightness is high", func(t *testing.T) {
			t.Run("is true", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "750\n"
				}

				actual := isHighBrightness()
				expected := true

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})

		t.Run("when the current brightness is low (not high)", func(t *testing.T) {
			t.Run("is false", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "123\n"
				}

				actual := isHighBrightness()
				expected := false

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})

		t.Run("when the current brightness is strange", func(t *testing.T) {
			t.Run("is false", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "Hello!\n"
				}

				actual := isHighBrightness()
				expected := false

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})
	})
}

func TestCurrentBrightness(t *testing.T) {
	t.Run("currentBrightness", func(t *testing.T) {
		t.Run("when the file is read", func(t *testing.T) {
			t.Run("returns the parsed value", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "123\n"
				}

				actual := currentBrightness()
				expected := 123

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})

		t.Run("when the file cannot be read", func(t *testing.T) {
			t.Run("returns a fallback (-1)", func(t *testing.T) {
				fileReader = func(_ string) string {
					return "Sorry!\n"
				}

				actual := currentBrightness()
				expected := -1

				if expected != actual {
					t.Fatalf("expected: %v · got: %v", expected, actual)
				}
			})
		})
	})
}

func TestSetLowBrightness(t *testing.T) {
	t.Run("setLowBrightness", func(t *testing.T) {
		t.Run("tells the system to set the brightness to low", func(t *testing.T) {
			var fileWritten string
			var contentWritten string

			fileWriter = func(filename, str string) error {
				fileWritten = filename
				contentWritten = str

				return nil
			}

			setLowBrightness()

			if fileWritten != "/sys/class/backlight/gmux_backlight/brightness" {
				t.Fatalf("did not write to expected file; wrote to: %s", fileWritten)
			}
			if contentWritten != "250" {
				t.Fatalf("did not write expected content; wrote: %s", contentWritten)
			}
		})
	})
}

func TestSetHighBrightness(t *testing.T) {
	t.Run("setHighBrightness", func(t *testing.T) {
		t.Run("tells the system to set the brightness to high", func(t *testing.T) {
			var fileWritten string
			var contentWritten string

			fileWriter = func(filename, str string) error {
				fileWritten = filename
				contentWritten = str

				return nil
			}

			setHighBrightness()

			if fileWritten != "/sys/class/backlight/gmux_backlight/brightness" {
				t.Fatalf("did not write to expected file; wrote to: %s", fileWritten)
			}
			if contentWritten != "750" {
				t.Fatalf("did not write expected content; wrote: %s", contentWritten)
			}
		})
	})
}

func TestCanWriteBrightness(t *testing.T) {
	t.Run("canWriteBrightness", func(t *testing.T) {
		userId = func() int {
			return testUserId
		}
		groupId = func() int {
			return testGroupId
		}

		t.Run("when the brightness file is user-writable and a bit more", func(t *testing.T) {
			t.Run("is true", func(t *testing.T) {
				filePermissionCheck = func(_ string) (int, int, fs.FileMode) {
					return testUserId, testGroupId, fs.FileMode(uint32(0700))
				}

				actual := canWriteBrightness()
				expected := true

				if expected != actual {
					t.Fatalf("expected %v; got: %v", expected, actual)
				}
			})
		})
		t.Run("when the brightness file is user-writable", func(t *testing.T) {
			t.Run("is true", func(t *testing.T) {
				filePermissionCheck = func(_ string) (int, int, fs.FileMode) {
					return testUserId, testGroupId, fs.FileMode(uint32(0600))
				}

				actual := canWriteBrightness()
				expected := true

				if expected != actual {
					t.Fatalf("expected %v; got: %v", expected, actual)
				}
			})
		})

		t.Run("when the brightness file is user-readable", func(t *testing.T) {
			t.Run("is true", func(t *testing.T) {
				filePermissionCheck = func(_ string) (int, int, fs.FileMode) {
					return testUserId, testGroupId, fs.FileMode(uint32(0400))
				}

				actual := canWriteBrightness()
				expected := false

				if expected != actual {
					t.Fatalf("expected %v; got: %v", expected, actual)
				}
			})
		})

		t.Run("when the brightness file is user-inaccessible", func(t *testing.T) {
			t.Run("is true", func(t *testing.T) {
				filePermissionCheck = func(_ string) (int, int, fs.FileMode) {
					return rootUserId, rootGroupId, fs.FileMode(uint32(0000))
				}

				actual := canWriteBrightness()
				expected := false

				if expected != actual {
					t.Fatalf("expected %v; got: %v", expected, actual)
				}
			})
		})
	})
}
