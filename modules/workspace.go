package modules

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/proxypoke/i3ipc"
)

// workspaceTrigger Specifies the trigger.
const workspaceTrigger string = "wsp"

// WorkspaceModule is a module for workspaces info
type WorkspaceModule struct {
	colorBg   string
	colorFg   string
	colorBgHl string

	icon string

	currentWorkspace string

	Alignment string
	Interval  time.Duration
}

// Initialize from config file. TODO
func (wsp *WorkspaceModule) Initialize() {
	wsp.colorBg = "#ff282828"
	wsp.colorBgHl = "#ff665c54"
	wsp.colorFg = "#ffebdbb2"

	wsp.icon = ""

	if wsp.Interval == 0 {
		wsp.Interval = time.Second * 60
	}
	if wsp.Alignment == "" {
		wsp.Alignment = "l"
	}

}

// IsTriggered specifies whether test is one of its trigger keywords.
func (wsp *WorkspaceModule) IsTriggered(test string) bool {
	return strings.HasPrefix(test, workspaceTrigger)
}

// Info is the main wrapper function for getting Information.
func (wsp *WorkspaceModule) Info(output chan string, trigger chan string) {
	ipcsocket, err := i3ipc.GetIPCSocket()
	if err != nil {
		panic(err)
	}
	wsEvents, err := i3ipc.Subscribe(i3ipc.I3WorkspaceEvent)
	if err != nil {
		panic(err)
	}
	output <- wsp.process(ipcsocket)
	ticker := time.NewTicker(wsp.GetInterval())
	for {
		select {
		case <-ticker.C:
			output <- wsp.process(ipcsocket)
		case <-wsEvents:
			output <- wsp.process(ipcsocket)
		case value := <-trigger:
			if value != fmt.Sprintf("%s%s", workspaceTrigger, wsp.currentWorkspace) {
				success, err := ipcsocket.Command(fmt.Sprintf("workspace %s", strings.TrimPrefix(value, workspaceTrigger)))
				if err != nil {
					log.Println(err)
				}
				if success {
					output <- wsp.process(ipcsocket)
				}
			}
		}
	}
}

// GetAlignment allows us to get alignment.
func (wsp *WorkspaceModule) GetAlignment() string {
	return wsp.Alignment
}

// GetInterval allows us to get interval duration.
func (wsp *WorkspaceModule) GetInterval() time.Duration {
	return wsp.Interval
}

func (wsp *WorkspaceModule) process(ipc *i3ipc.IPCSocket) string {

	workspaceString := fmt.Sprintf("%%{F%s B%s} %s ", wsp.colorFg, wsp.colorBg, wsp.icon)

	workspaces, err := ipc.GetWorkspaces()
	if err != nil {
		panic(err)
	}
	for _, workspace := range workspaces {
		if workspace.Focused {
			wsp.currentWorkspace = workspace.Name
			workspaceString += fmt.Sprintf("%%{F%s B%s A:%s%s:} %s %%{F%s B%s A}", wsp.colorBg, wsp.colorBgHl, workspaceTrigger, workspace.Name, strings.TrimLeft(workspace.Name, "1234567890"), wsp.colorFg, wsp.colorBg)
		} else {
			workspaceString += fmt.Sprintf("%%{A:%s%s:} %s %%{A}", workspaceTrigger, workspace.Name, strings.TrimLeft(workspace.Name, "0123456789"))
		}
	}
	workspaceString += fmt.Sprintf("%%{F%s B%s} ", wsp.colorFg, wsp.colorBg)
	return workspaceString

}

// Ensure Interface is Implemented
var _ Module = (*WorkspaceModule)(nil)
