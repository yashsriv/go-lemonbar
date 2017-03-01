package modules

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ShellModule is a module for getting info via shell scripts
type ShellModule struct {
	colorBg string
	colorFg string

	Alignment string
	Icon      string
	Interval  time.Duration
	MaxLen    int
	Script    string
	Trigger   string
}

// Initialize from config file. TODO
func (s *ShellModule) Initialize() {
	s.colorBg = "#ff282828"
	s.colorFg = "#ffebdbb2"

	if s.Interval == 0 {
		s.Interval = time.Second
	}
	if s.Alignment == "" {
		s.Alignment = "r"
	}
	if s.Trigger == "" {
		s.Trigger = s.Script
	}

}

// IsTriggered specifies whether test is one of its trigger keywords.
func (s *ShellModule) IsTriggered(test string) bool {
	return test == s.Trigger
}

// Info is the main wrapper function for getting Shell Information.
func (s *ShellModule) Info(output chan string, trigger chan string) {
	output <- s.shellProcess()
	ticker := time.NewTicker(s.GetInterval())
	for {
		select {
		case <-ticker.C:
			output <- s.shellProcess()
		case <-trigger:
			output <- s.shellProcess()
		}
	}
}

// GetAlignment allows us to get b's alignment.
func (s *ShellModule) GetAlignment() string {
	return s.Alignment
}

// GetInterval allows us to get b's interval duration.
func (s *ShellModule) GetInterval() time.Duration {
	return s.Interval
}

func (s *ShellModule) shellProcess() string {

	var format string
	format = "%%{F%s B%s A:%s:} %s %s %%{A F- B-}"
	out, err := exec.Command(s.Script).Output()
	if err != nil {
		return fmt.Sprintf(format, s.colorFg, s.colorBg, s.Trigger, s.Icon, "!!")
	}
	outputString := strings.TrimSpace(string(out))
	var finalOutput string
	if s.MaxLen == 0 || len(outputString) < s.MaxLen {
		finalOutput = outputString
	} else {
		finalOutput = outputString[0:s.MaxLen]
	}

	return fmt.Sprintf(format, s.colorFg, s.colorBg, s.Trigger, s.Icon, finalOutput)

}

// Ensure Interface is Implemented
var _ Module = (*ShellModule)(nil)
