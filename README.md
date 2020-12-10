# OctoScreen [![GitHub release](https://img.shields.io/github/release/Z-Bolt/OctoScreen.svg)](https://github.com/Z-Bolt/OctoScreen/releases) [![license](https://img.shields.io/github/license/Z-Bolt/OctoScreen.svg)]()

_OctoScreen_ is a LCD touch interface for your OctoPrint server.  It is based on GTK+3 and allows you to control your 3D Printer using a [LCD touch screen](https://amzn.to/2L8cRkR), a [Raspberry Pi](https://amzn.to/39LPvvF), and [OctoPrint](https://octoprint.org/).  It's an _X application_ that's executed directly in the X Server without a window manager or browser, and opreates as a frontend for OctoPrint.

<img width="480" src="https://user-images.githubusercontent.com/390214/60487814-ef9d1a00-9ca8-11e9-9c48-31bf54a5488d.png" />
<img width="240" src="https://user-images.githubusercontent.com/390214/60277300-f4a74580-9905-11e9-8b88-f6cc35533c2a.png" />
<img width="240" src="https://user-images.githubusercontent.com/390214/60277572-84e58a80-9906-11e9-8334-202544f0191d.png" />

Some of the functionality of OctoScreen incude:

- Print jobs monitoring.
- Temperature and Filament management.
- Jogging operations.
- Wifi connection management.
- Toolchanger management tools.




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

For example, to install on Raspbian or OctoPi:
```sh
wget https://github.com/Z-Bolt/OctoScreen/releases/download/v2.5.1/octoscreen_2.5-1_armhf.deb
sudo dpkg -i octoscreen_2.5-1_armhf.deb
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

- `OCTOPRINT_HOST` - The OctoPrint HTTP address.  Example: `http://octopi.local` or `http://localhost:5000`.  Note: the protocol (`http://` or `https://`) is required.  If the setting for OCTOPRINT_HOST does not contain the protocol, an error will be displayed when OctoScreen starts.

- `OCTOPRINT_APIKEY` - OctoScreen expects an [API key]( http://docs.octoprint.org/en/master/api/general.html) to be supplied. This API key can be either the globally configured key, or a user-specific one if “Access Control” is enabled.

- `OCTOSCREEN_STYLE_PATH` - Several themes are supported, and style configurations can be done through CSS.  This variable defines the location of the application theme.


#### Optional Configuration Settings

- `OCTOPRINT_CONFIG_FILE` - The location of the OctoPrint's config.yaml file.  If empty, the file path used will be the `pi` home folder of the current user. Only used for locally installed OctoPrint servers.

- `OCTOSCREEN_LOG_FILE_PATH` - The file path to where the log file will be saved.  The file path must be in the location of where the app is run, and can not use external paths (eq, ~/ or / are out).  The file path should not only include the path to the log file, but the name of the log file as well (eq logs/logfile.txt).  The log file is appended to, and never automatically truncated and will grow over time.  If you turn log file logging on (by specifying a path), be sure to turn it off (by setting the value to "").

- `OCTOSCREEN_LOG_LEVEL` - Controls the level of logging.  Accepted values are (with increasing levels): debug, info, warn, and error.  If no value is provided, the log level will default to warn.

- `OCTOSCREEN_RESOLUTION` - Resolution of the application, and should be configured to the resolution of your screen.  Optimal resolution for OctoScreen is no less than 800x480, so if the physical resolution of your screen is 480x320, it's recommended to set the software resolution 800x533.  If you are using Raspbian you can do it by changing [`hdmi_cvt`](https://www.raspberrypi.org/documentation/configuration/config-txt/video.md) param in `/boot/config.txt` file.  Please see [Setting Up OctoScreen and Your Display](https://github.com/Z-Bolt/OctoScreen/wiki/Setting-Up-OctoScreen-and-Your-Display) and [Installing OctoScreen with a 3.5" 480x320 TFT screen](https://github.com/Z-Bolt/OctoScreen/wiki/Installing-OctoScreen-with-a-3.5%22-480x320-TFT-screen) for more information.




------------
## Menu Configuration

### Custom Controls and Commands

Custom [controls](http://docs.octoprint.org/en/master/configuration/config_yaml.html#controls) to execute GCODE instructions and [commands](http://docs.octoprint.org/en/master/configuration/config_yaml.html#system) to execute shell commands can be defined in the `config.yaml` file.

The controls are limit to static controls without `inputs`.





------------
## [Roadmap](https://github.com/Z-Bolt/OctoScreen/projects/2)





------------
## License

GNU Affero General Public License v3.0, see [LICENSE](LICENSE)

This project is a hard fork from [Octiprint-TFT](https://github.com/mcuadros/OctoPrint-TFT) created by [@mcuadros](https://github.com/mcuadros/OctoPrint-TFT)
