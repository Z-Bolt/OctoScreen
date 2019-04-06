go-octoprint [![Build Status](https://travis-ci.org/mcuadros/go-octoprint.svg?branch=master)](https://travis-ci.org/mcuadros/go-octoprint) [![GoDoc](http://godoc.org/github.com/mcuadros/go-octoprint?status.svg)](http://godoc.org/github.com/mcuadros/go-octoprint)
==============================

Go library for accessing the [OctoPrint](http://octoprint.org/)'s [REST API](http://docs.octoprint.org/en/master/api/index.html).

Installation
------------

The recommended way to install go-octoprint

```
go get github.com/mcuadros/go-octoprint
```

Example
-------

### Retrieving the current connection state:

```go
client, _ := NewClient("<octoprint-url>", "<api-key>")

r := octoprint.ConnectionRequest{}
s, err := r.Do(client)
if err != nil {
  log.Error("error requesting connection state: %s", err)
}

fmt.Printf("Connection State: %q\n", s.Current.State)
```


### Retrieving current temperature for bed and extruders:

```go
r := octoprint.StateRequest{}
s, err := r.Do(c)
if err != nil {
	log.Error("error requesting state: %s", err)
}

fmt.Println("Current Temperatures:")
for tool, state := range s.Temperature.Current {
	fmt.Printf("- %s: %.1f°C / %.1f°C\n", tool, state.Actual, state.Target)
}
```

## Implemented Methods

### [Version Information](http://docs.octoprint.org/en/master/api/version.html)
- [x] GET `/api/version`

### [Apps](http://docs.octoprint.org/en/master/api/apps.html)
- [ ] GET `/apps/auth`
- [ ] POST `/apps/auth`

### [Connection Operations](http://docs.octoprint.org/en/master/api/connection.html)
- [x] GET `/api/connection`
- [x] POST `/api/connection`

### [File Operations](http://docs.octoprint.org/en/master/api/files.html)
- [x] GET `/api/files
- [x] GET `/api/files/<location>`
- [x] POST `/api/files/<location>`
- [x] GET `/api/files/<location>/<filename>`
- [x] POST `/api/files/<location>/<path>` (Only select command)
- [x] DELETE `/api/files/<location>/<path>`

### [Job Operations](http://docs.octoprint.org/en/master/api/job.html)
- [x] POST `/api/job`
- [x] GET `/api/job`

### [Languages](http://docs.octoprint.org/en/master/api/languages.html)
- [ ] GET `/api/languages`
- [ ] POST `/api/languages`
- [ ] DELETE `/api/languages/<locale>/<pack>`

### [Log file management](http://docs.octoprint.org/en/master/api/logs.html)
- [ ] GET `/api/logs`
- [ ] DELETE `/api/logs/<filename>`

### [Printer Operations](http://docs.octoprint.org/en/master/api/printer.html)
- [x] GET `/api/printer`
- [x] POST `/api/printer/printhead`
- [x] POST `/api/printer/tool`
- [x] GET `/api/printer/tool`
- [x] POST `/api/printer/bed`
- [x] GET `/api/printer/bed`
- [x] POST `/api/printer/sd`
- [x] GET `/api/printer/sd`
- [x] POST `/api/printer/command`
- [x] GET `/api/printer/command/custom` ([un-documented at REST API](https://github.com/foosel/OctoPrint/blob/7f5d03d0549bcbd26f40e7e4a3297ea5204fb1cc/src/octoprint/server/api/printer.py#L376))

### [Printer profile operations](http://docs.octoprint.org/en/master/api/printerprofiles.html)
- [ ] GET `/api/printerprofiles`
- [ ] POST `/api/printerprofiles`
- [ ] PATCH `/api/printerprofiles/<profile>`
- [ ] DELETE `/api/printerprofiles/<profile>`

### [Settings](http://docs.octoprint.org/en/master/api/settings.html)
- [x] GET `/api/settings`
- [ ] POST `/api/settings`
- [ ] POST `/api/settings/apikey`

### [Slicing](http://docs.octoprint.org/en/master/api/slicing.html)
- [ ] GET `/api/slicing`
- [ ] GET `/api/slicing/<slicer>/profiles`
- [ ] GET `/api/slicing/<slicer>/profiles/<key>`
- [ ] PUT `/api/slicing/<slicer>/profiles/<key>`
- [ ] DELETE `/api/slicing/<slicer>/profiles/<key>`

### [System](http://docs.octoprint.org/en/master/api/system.html)
- [x] GET `/api/system/commands`
- [ ] GET `/api/system/commands/<source>`
- [x] POST `/api/system/commands/<source>/<action>`

### [Timelapse](http://docs.octoprint.org/en/master/api/timelapse.html)
- [ ] GET `/api/timelapse`
- [ ] DELETE `/api/timelapse/<filename>`
- [ ] POST `/api/timelapse/unrendered/<name>`
- [ ] DELETE `/api/timelapse/unrendered/<name>`
- [ ] POST `/api/timelapse`

### [User](http://docs.octoprint.org/en/master/api/users.html)
- [ ] GET `/api/users`
- [ ] GET `/api/users/<username>`
- [ ] POST `/api/users`
- [ ] PUT `/api/users/<username>`
- [ ] DELETE `/api/users/<username>`
- [ ] PUT `/api/users/<username>/password`
- [ ] GET `/api/users/<username>/settings`
- [ ] PATCH `/api/users/<username>/settings`
- [ ] POST `/api/users/<username>/apikey`
- [ ] DELETE `/api/users/<username>/apikey`

### [Util](http://docs.octoprint.org/en/master/api/util.html)
- [ ] POST `/api/util/test`

### [Wizard](http://docs.octoprint.org/en/master/api/wizard.html)
- [ ] GET `/setup/wizard`
- [ ] POST `/setup/wizard`

License
-------

MIT, see [LICENSE](LICENSE)
