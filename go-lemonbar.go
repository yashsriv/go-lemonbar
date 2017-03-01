package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/yashsriv/go-lemonbar/modules"
)

// List of activated modules
var activated = []modules.Module{
	&modules.WorkspaceModule{
		Alignment: "l",
		Interval:  time.Second * 60,
	},
	&modules.ShellModule{
		Alignment: "l",
		Icon:      "",
		Interval:  time.Second * 1,
		MaxLen:    20,
		Script:    "/home/yash/.i3/bin/current_window",
		Trigger:   "window",
	},
	&modules.ShellModule{
		Alignment: "c",
		Icon:      "",
		Interval:  time.Second * 90,
		Script:    "/home/yash/.i3/bin/fortune_lemonbar",
		Trigger:   "fortune",
	},
	&modules.BrightnessModule{
		Interval: time.Second * 1,
	},
	&modules.ShellModule{
		Alignment: "r",
		Icon:      "",
		Interval:  time.Second * 1,
		Script:    "/home/yash/.i3/bin/volume",
		Trigger:   "volume",
	},
	&modules.IpModule{
		Interval: time.Second,
	},
	&modules.BatteryModule{
		Interval: time.Second * 10,
	},
	&modules.DateModule{
		Interval: time.Second * 10,
	},
}

func main() {
	colorBack := "#ff282828"
	colorFore := "#ffebdbb2"

	lemonbar := exec.Command("lemonbar",
		"-p",
		"-f", "SauceCodePro Nerd Font-12:style=Semibold",
		"-B", colorBack,
		"-F", colorFore,
		"-u", "8",
		"-o", "-2",
		"-a", fmt.Sprintf("%d", len(activated)+9))
	input, err := lemonbar.StdinPipe()
	if err != nil {
		panic(err)
	}
	output, err := lemonbar.StdoutPipe()
	if err != nil {
		panic(err)
	}
	err = lemonbar.Start()
	if err != nil {
		panic(err)
	}

	initializeModules()
	number := len(activated)
	triggerChans := make([]chan string, number)
	for i := range triggerChans {
		triggerChans[i] = make(chan string)
	}
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		lemonbar.Process.Kill()
		os.Exit(0)
	}()
	go feedbackLoop(&output, triggerChans)
	loop(&input, triggerChans)
}

func initializeModules() {
	for _, module := range activated {
		module.Initialize()
	}
}

func loop(output *io.WriteCloser, triggerChans []chan string) {
	number := len(activated)
	channels := make([]chan string, number)
	for i := range channels {
		channels[i] = make(chan string)
	}
	for i, module := range activated {
		go module.Info(channels[i], triggerChans[i])
	}
	cases := make([]reflect.SelectCase, len(channels))
	values := make([]string, len(channels))
	for i, ch := range channels {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	for {
		chosen, value, _ := reflect.Select(cases)
		values[chosen] = value.String()
		var outputString string
		lfound := false
		cfound := false
		rfound := false
		for i, module := range activated {
			if !lfound && module.GetAlignment() == "l" {
				outputString += "%{l}"
				lfound = true
			}
			if !cfound && module.GetAlignment() == "c" {
				outputString += "%{c}"
				lfound = true
				cfound = true
			} else if !rfound && module.GetAlignment() == "r" {
				outputString += "%{r}"
				lfound = true
				cfound = true
				rfound = true
			}
			outputString += fmt.Sprintf("%s", values[i])
		}
		outputString += "\n"
		// log.Print(outputString)
		(*output).Write([]byte(outputString))
	}
}

func feedbackLoop(input *io.ReadCloser, triggerChans []chan string) {
	scanner := bufio.NewScanner(*input)
	for scanner.Scan() {
		line := scanner.Text()
		for i, module := range activated {
			if module.IsTriggered(line) {
				triggerChans[i] <- line
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("Reading Standard Input: ", err)
	}
}
