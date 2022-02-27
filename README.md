# OctoScreen [![GitHub release](https://img.shields.io/github/release/Z-Bolt/OctoScreen.svg)](https://github.com/Z-Bolt/OctoScreen/releases) [![license](https://img.shields.io/github/license/Z-Bolt/OctoScreen.svg)]()

_OctoScreen_ is a LCD touch interface for your OctoPrint server.  It is based on GTK+3 and allows you to control your 3D Printer using a [LCD touch screen](https://amzn.to/2L8cRkR), a [Raspberry Pi](https://amzn.to/39LPvvF), and [OctoPrint](https://octoprint.org/).  It's an _X application_ that's executed directly in the X Server without a window manager or browser, and operates as a frontend for OctoPrint.

<img width="480" src="https://user-images.githubusercontent.com/10328858/101729412-f66ab780-3a6c-11eb-8fd8-8bbf5c8c1dc7.png" />
Idle
<br />
<br />
<br />

<img width="480" src="https://user-images.githubusercontent.com/10328858/101729660-6ed17880-3a6d-11eb-80b0-b1170d1e59f9.png" />
Idle-Multiple Hotends
<br />
<br />
<br />

<img width="480" src="https://user-images.githubusercontent.com/10328858/101729636-5feac600-3a6d-11eb-9121-83808b3decb7.png" />
Printing
<br />
<br />
<br />

<img width="480" src="https://user-images.githubusercontent.com/10328858/101729487-14381c80-3a6d-11eb-8e24-298dc3b34fe4.png" />
Home
<br />
<br />
<br />

<img width="480" src="https://user-images.githubusercontent.com/10328858/101729522-24e89280-3a6d-11eb-85bb-ae0f1973b867.png" />
Filament
<br />
<br />
<br />

<img width="480" src="https://user-images.githubusercontent.com/10328858/101729592-4a759c00-3a6d-11eb-8086-daab19cd6ff5.png" />
Actions
<br />
<br />
<br />



Some of the functionality of OctoScreen included:

- Print jobs monitoring.
- Temperature and Filament management.
- Jogging operations.
- WiFi connection management.
- Tool changer management tools.




## How Is OctoScreen Different From TouchUI?

[TouchUI](http://plugins.octoprint.org/plugins/touchui/) is an amazing plugin for OctoPrint.  It was created as a responsive design application to access OctoPrint from low resolution devices, such as smartphones, tablets, etc.

Executing TouchUI under a RPi with a TFT display presents two big problems: First, it isn't optimized to be used with resistive touch screens with low resolutions like 480x320, and second, it requires a browser to run, which consumes a lot of resources.  This is the main reason OctoScreen was developed as an X application.




------------
## Installation

### Dependencies

OctoScreen is based on [GoLang](https://golang.org).  GoLang applications are usually dependency-less, but in this case [GTK+3](https://developer.gnome.org/gtk3/3.0/gtk.html) is used, and the GTK+3 libraries are required in order to run.  Be sure that GTK+3 is installed and is the only graphical environment that's been installed.

If you are using `Raspbian` or any other `Debian` based distribution, required packages can be installed using:

```sh
sudo apt-get install libgtk-3-0 xserver-xorg xinit x11-xserver-utils
```

You will also need to set up the video drivers for the display you are using.  Installation and configuration of the drivers is usually specific to the display you are using, and you will need to contact the manufacturer for instructions.  To help you set up your system and display, a setup guide ([Setting Up OctoScreen and Your Display](https://github.com/Z-Bolt/OctoScreen/wiki/Setting-Up-OctoScreen-and-Your-Display)) is available in the wiki.



### Install From a .deb Package

There are two ways to install OctoScreen: the recommended and supported way is to install from a .deb package, or you can choose to install from building the sources yourself.  The recommended way to install OctoScreen is use the `.deb` packages from the [Releases](https://github.com/Z-Bolt/OctoScreen/releases) page.  The packages are available for Debian Stretch based distributions such as Raspbian and OctoPi.

For example, to install on a new RaspberryPi with OctoPi:
```sh
wget https://github.com/Z-Bolt/OctoScreen/releases/download/v2.8.0/octoscreen_2.8.0_armhf.deb
sudo dpkg -i octoscreen_2.8.0_armhf.deb
```

Or to update an existing version of OctoScreen:
```sh
wget https://github.com/Z-Bolt/OctoScreen/releases/download/v2.8.0/octoscreen_2.8.0_armhf.deb
sudo dpkg -r octoscreen
sudo dpkg -i octoscreen_2.8.0_armhf.deb
sudo reboot now
```

### Install From Source

The compilation and packaging tasks are managed by the [`Makefile`](Makefile) and backed on [Docker](Dockerfile).  Docker is used to avoid installing any other dependencies since all the operations are done inside of the container.

If you need to install docker inside `Raspbian` or any other linux distribution just run:

```sh
curl -fsSL get.docker.com -o get-docker.sh
sh get-docker.sh
```

You can read more about this at [`docker-install`](https://github.com/docker/docker-install)

To compile the project (assuming that you already cloned this repository), just execute the `build` target.  This will generate all the binaries and debian packages in the `build` folder:

```sh
make build
ls -1 build/
```

The default build is for the STRETCH release of debian, but BUSTER and JESSIE are also possible.  To build one of these targets, you just have to specify the package during make.
Example for BUSTER:
```sh
make build DEBIAN_PACKAGES=BUSTER
ls -1 build/
```

If you are using `Raspbian` you can install any of the `.deb` generated packages.  If not, just use the compiled binary.




------------
## Configuration

### Basic Configuration

The basic configuration is handled via environment variables, if you are using the `.deb` package you can configure it at `/etc/octoscreen/config`.

#### Required Configuration Settings

- `OCTOPRINT_HOST` - The URL to the OctoPrint server.  Example: `http://octopi.local` or `http://localhost:5000`.  Note: the protocol (`http://` or `https://`) is required.  If the setting for `OCTOPRINT_HOST` does not contain the protocol, an error will be displayed when OctoScreen starts.

- `OCTOPRINT_APIKEY` - OctoScreen expects an [API key]( http://docs.octoprint.org/en/master/api/general.html) to be supplied. The API key can either be set in OctoScreen's config file, or in OctoPrint's config file (OCTOPRINT_CONFIG_FILE below)

- `OCTOSCREEN_STYLE_PATH` - Several themes are supported, and style configurations can be done through CSS.  This variable defines the location of the application theme.


#### Optional Configuration Settings

- `OCTOPRINT_CONFIG_FILE` - The location of OctoPrint's config.yaml file.  If empty, the file path used will be the `pi` home folder of the current user.  The OCTOPRINT_APIKEY is required, and if it isn't defined in OctoScreen's config file (see OCTOPRINT_APIKEY above) it needs to be defined in OctoPrint's config file.

- `OCTOSCREEN_LOG_FILE_PATH` - The file path to where the log file will be saved.  The file path should be a fully qualified path and not only include the path to the log file, but the name of the log file as well (eg `/home/pi/logs/logfile.txt`).  The log file is appended to and is never automatically truncated, and will grow over time.  If you turn log file logging on (by specifying a path), be sure to turn it off (by setting the value to "").

- `OCTOSCREEN_LOG_LEVEL` - Controls the level of logging.  Accepted values are (with increasing levels): `debug`, `info`, `warn`, and `error`.  If no value is provided, the log level will default to `warn`.

- `OCTOSCREEN_RESOLUTION` - Resolution of the application, and should be configured to the resolution of your screen.  Optimal resolution for OctoScreen is no less than 800x480, so if the physical resolution of your screen is 480x320, it's recommended to set the software resolution 800x533.  If you are using Raspbian you can do it by changing [`hdmi_cvt`](https://www.raspberrypi.org/documentation/configuration/config-txt/video.md) param in `/boot/config.txt` file.  Please see [Setting Up OctoScreen and Your Display](https://github.com/Z-Bolt/OctoScreen/wiki/Setting-Up-OctoScreen-and-Your-Display) and [Installing OctoScreen with a 3.5" 480x320 TFT screen](https://github.com/Z-Bolt/OctoScreen/wiki/Installing-OctoScreen-with-a-3.5%22-480x320-TFT-screen) for more information.

- `DISPLAY_CURSOR` - To display the cursor, add `DISPLAY_CURSOR=true` to your config file.  In order to display the cursor, you will also need to edit `/lib/systemd/system/octoscreen.service` and remove `-nocursor`



------------
## Menu Configuration

### Custom Controls and Commands

Custom [controls](http://docs.octoprint.org/en/master/configuration/config_yaml.html#controls) to execute GCODE instructions and [commands](http://docs.octoprint.org/en/master/configuration/config_yaml.html#system) to execute shell commands can be defined in the `config.yaml` file.

The controls are limit to static controls without `inputs`.




------------
## Wiki
For troubleshooting and general information about this project, be sure to check out the Wiki page, located at https://github.com/Z-Bolt/OctoScreen/wiki



------------
<!--
## [Roadmap](https://github.com/Z-Bolt/OctoScreen/projects/2)
-->
## Roadmap
https://github.com/Z-Bolt/OctoScreen/wiki/Project-Roadmap




------------
## License

GNU Affero General Public License v3.0, see [LICENSE](LICENSE)

This project is a hard fork from [Octoprint-TFT](https://github.com/mcuadros/OctoPrint-TFT) created by [@mcuadros](https://github.com/mcuadros/OctoPrint-TFT)
