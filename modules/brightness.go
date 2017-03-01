package modules

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// dateTrigger Specifies the trigger.
const brightnessTrigger string = "bright"

// BrightnessModule is a module for brightness info
type BrightnessModule struct {
	colorBg string
	colorFg string

	iconBri string

	Alignment string
	Interval  time.Duration
}

// Initialize from config file. TODO
func (b *BrightnessModule) Initialize() {
	b.colorBg = "#ff282828"
	b.colorFg = "#ffebdbb2"

	b.iconBri = ""

	if b.Interval == 0 {
		b.Interval = time.Second
	}
	if b.Alignment == "" {
		b.Alignment = "r"
	}

}

// IsTriggered specifies whether test is one of its trigger keywords.
func (b *BrightnessModule) IsTriggered(test string) bool {
	return test == brightnessTrigger
}

// Info is the main wrapper function for getting Information.
func (b *BrightnessModule) Info(output chan string, trigger chan string) {
	output <- b.process()
	ticker := time.NewTicker(b.GetInterval())
	for {
		select {
		case <-ticker.C:
			output <- b.process()
		case <-trigger:
			output <- b.process()
		}
	}
}

// GetAlignment allows us to get alignment.
func (b *BrightnessModule) GetAlignment() string {
	return b.Alignment
}

// GetInterval allows us to get interval duration.
func (b *BrightnessModule) GetInterval() time.Duration {
	return b.Interval
}

func (b *BrightnessModule) process() string {

	format := "%%{F%s B%s A:%s:} %s %s %%{A F- B-}"
	dat, err := ioutil.ReadFile("/sys/class/backlight/intel_backlight/brightness")
	if err != nil {
		return fmt.Sprintf(format, b.colorFg, b.colorBg, brightnessTrigger, b.iconBri, "!!")
	}
	current, err := strconv.ParseInt(strings.TrimSpace(string(dat)), 10, 0)
	if err != nil {
		return fmt.Sprintf(format, b.colorFg, b.colorBg, brightnessTrigger, b.iconBri, "!!")
	}
	dat, err = ioutil.ReadFile("/sys/class/backlight/intel_backlight/max_brightness")
	if err != nil {
		return fmt.Sprintf(format, b.colorFg, b.colorBg, brightnessTrigger, b.iconBri, "!!")
	}
	max, err := strconv.ParseInt(strings.TrimSpace(string(dat)), 10, 0)
	if err != nil {
		return fmt.Sprintf(format, b.colorFg, b.colorBg, brightnessTrigger, b.iconBri, "!!")
	}
	return fmt.Sprintf(format, b.colorFg, b.colorBg, brightnessTrigger, b.iconBri, fmt.Sprintf("%d", current*100/max))

}

// Ensure Interface is Implemented
var _ Module = (*BrightnessModule)(nil)
