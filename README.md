# go-lemonbar
Lemonbar Configuration via Go.

![Preview](./screenshots/full.png)

Looks and feels inspired by [@CopperBadger's Config](https://github.com/CopperBadger/dotfiles/tree/master/dots/.i3/lemonbar-new)

Completely Modular. Each module runs in its own goroutine and handles
click signals as well as workspace switching and updates automatically
in certain time intervals.

Requires:

* `lemonbar-xft` - Taken from AUR
* `go`

Install Using:

```
go install
```

Current Modules:

* DateModule - For Current Date and Time. Simplest Module which serves as a template for other modules.

  ![Preview Date](./screenshots/date.png)

* BatteryModule - For Battery Status. Uses Fontawesome icons and requires:
  `github.com/distatus/battery` - Easily Installed via `go get -u github.com/distatus/battery`
  
  ![Preview Battery](./screenshots/battery.png)

* BrightnessModule - For Brightness Value. Uses `/sys/class/backlight/intel-backlight`

  ![Preview Brightness](./screenshots/brightness.png)

* IpModule - Finds currently active network interface using `net` package

  ![Preview Network](./screenshots/ip.png)

* WorkspaceModule - For communicating with i3 and maintaining currentWorkspace. Requires:
  `github.com/proxypoke/i3ipc` - Easily Installed via `go get -u github.com/proxypoke/i3ipc`

  ![Preview Network](./screenshots/wsp.png)

* ShellModule - For those of us who can't get over our shell scripts :stuck_out_tongue: 

  ![Preview Network](./screenshots/fortune.png)

# TODO:

* Read all configuration from config file instead of hardcoding into code.
