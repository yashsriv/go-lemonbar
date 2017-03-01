package modules

import (
	"fmt"
	"time"

	"github.com/distatus/battery"
)

// BatteryTrigger Specifies the trigger for BatteryInfo.
const batteryTrigger string = "battery"

// BatteryModule is a module for battery info
type BatteryModule struct {
	colorBg string
	colorFg string

	iconCha string
	iconBat []string

	Alignment string
	Interval  time.Duration
}

// Initialize from config file. TODO
func (b *BatteryModule) Initialize() {
	b.colorBg = "#ff665c54"
	b.colorFg = "#ffebdbb2"

	b.iconCha = ""
	b.iconBat = []string{"", "", "", "", ""}

	if b.Interval == 0 {
		b.Interval = time.Second
	}
	if b.Alignment == "" {
		b.Alignment = "r"
	}

}

// IsTriggered specifies whether test is one of its trigger keywords.
func (b *BatteryModule) IsTriggered(test string) bool {
	return test == batteryTrigger
}

// Info is the main wrapper function for getting Battery Information.
func (b *BatteryModule) Info(output chan string, trigger chan string) {
	output <- b.batteryProcess()
	ticker := time.NewTicker(b.GetInterval())
	for {
		select {
		case <-ticker.C:
			output <- b.batteryProcess()
		case <-trigger:
			output <- b.batteryProcess()
		}
	}
}

// GetAlignment allows us to get b's alignment.
func (b *BatteryModule) GetAlignment() string {
	return b.Alignment
}

// GetInterval allows us to get b's interval duration.
func (b *BatteryModule) GetInterval() time.Duration {
	return b.Interval
}

func (b *BatteryModule) batteryProcess() string {

	batteries, err := battery.GetAll()
	format := "%%{F%s B%s A:%s:} %s %1.0f%% %%{A F- B-}"
	if err != nil {
		return fmt.Sprintf(format, b.colorFg, b.colorBg, batteryTrigger, b.iconCha+" !! ", 0.0)
	}

	for _, bat := range batteries {
		charge := bat.Current * 100 / bat.Full
		if bat.State == battery.Charging {
			return fmt.Sprintf(format, b.colorFg, b.colorBg, batteryTrigger, b.iconCha, charge)
		}
		if charge < 15 {
			return fmt.Sprintf(format, b.colorFg, b.colorBg, batteryTrigger, b.iconBat[0], charge)
		}
		if charge < 45 {
			return fmt.Sprintf(format, b.colorFg, b.colorBg, batteryTrigger, b.iconBat[1], charge)
		}
		if charge < 65 {
			return fmt.Sprintf(format, b.colorFg, b.colorBg, batteryTrigger, b.iconBat[2], charge)
		}
		if charge < 95 {
			return fmt.Sprintf(format, b.colorFg, b.colorBg, batteryTrigger, b.iconBat[3], charge)
		}
		return fmt.Sprintf(format, b.colorFg, b.colorBg, batteryTrigger, b.iconBat[4], charge)
	}
	return fmt.Sprintf(format, b.colorFg, b.colorBg, batteryTrigger, b.iconCha+" !! ", 0.0)

}

// Ensure Interface is Implemented
var _ Module = (*BatteryModule)(nil)
