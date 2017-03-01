package modules

import (
	"fmt"
	"time"
)

// dateTrigger Specifies the trigger.
const dateTrigger string = "date"

// DateModule is a module for date info
type DateModule struct {
	colorBg string
	colorFg string

	iconCal string
	iconClo string

	Alignment string
	Interval  time.Duration
}

// Initialize from config file. TODO
func (d *DateModule) Initialize() {
	d.colorBg = "#ff665c54"
	d.colorFg = "#ffebdbb2"

	d.iconCal = ""
	d.iconClo = ""

	if d.Interval == 0 {
		d.Interval = time.Second
	}
	if d.Alignment == "" {
		d.Alignment = "r"
	}

}

// IsTriggered specifies whether test is one of its trigger keywords.
func (d *DateModule) IsTriggered(test string) bool {
	return test == dateTrigger
}

// Info is the main wrapper function for getting Information.
func (d *DateModule) Info(output chan string, trigger chan string) {
	output <- d.dateProcess()
	ticker := time.NewTicker(d.GetInterval())
	for {
		select {
		case <-ticker.C:
			output <- d.dateProcess()
		case <-trigger:
			output <- d.dateProcess()
		}
	}
}

// GetAlignment allows us to get alignment.
func (d *DateModule) GetAlignment() string {
	return d.Alignment
}

// GetInterval allows us to get interval duration.
func (d *DateModule) GetInterval() time.Duration {
	return d.Interval
}

func (d *DateModule) dateProcess() string {

	format := fmt.Sprintf("Mon 2 Jan  %s 15:04", d.iconClo)
	t := time.Now()
	date := t.Format(format)
	date = fmt.Sprintf("%%{F%s B%s A:%s:} %s %s %%{A F- B-}", d.colorFg, d.colorBg, dateTrigger, d.iconCal, date)
	return date

}

// Ensure Interface is Implemented
var _ Module = (*DateModule)(nil)
