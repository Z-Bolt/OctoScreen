OctoPrint-TFT
=============

_OctoPrint-TFT_, a touch interface for TFT touch modules based on GTK+3.

Is a _X application_ to be executed directly in the X Server without any windows
manager, as _frontend of a [OctoPrint](http://octoprint.org) server_ in a Raspberry Pi
equipped with any [TFT Touch module](https://www.waveshare.com/wiki/3.5inch_RPi_LCD_(A)).

Allows you to control your 3D Printer, like you can do with any [TFT/LCD panel](http://reprap.org/wiki/RepRapTouch), but using _OctoPrint_ and a Raspberry Pi.

<img width="480" src="https://user-images.githubusercontent.com/1573114/33559609-a73a969e-d90d-11e7-9cf2-cf212412aaa5.png" />


### How this is different from TouchUI?

[TouchUI](http://plugins.octoprint.org/plugins/touchui/), is an amazing plugin
for Octoprint, was created as a responsive design for access to OctoPrint,
from low resolution devices, such as smartphones, tablets, etc.

Executing TouchUI under a RPi w/TFT modules, presents two big problems,
first isn't optimized to be used with resistive touch screens with low resolutions
like 480x320 and second requires a browser to be access, consuming a lot of
resources.

This is the main reason because I develop this X application to be executed
in my 3d printer.

License
-------

MIT, see [LICENSE](LICENSE)
