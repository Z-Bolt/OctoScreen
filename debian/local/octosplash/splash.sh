#!/usr/bin/env bash

#OCTOSCREEN_DEFAULT_BOOT_SPLASH="/etc/octosplash/splash.png"
OCTOSCREEN_DEFAULT_BOOT_SPLASH="/home/pi/OctoScreen-2/debian/local/octosplash/splash.png"
FRAME_BUFFER=$1

if [[ -z "$FRAME_BUFFER" ]]; then
    echo "Unable to detect valid Frame Buffer for OctoScreen Boot Splash"
else
    if [[ -z "$OCTOSCREEN_BOOT_SPLASH" ]] || ! [[ -f "$OCTOSCREEN_BOOT_SPLASH" ]] || ! [[ $OCTOSCREEN_BOOT_SPLASH =~ .*\.png$ ]]; then
        OCTOSCREEN_BOOT_SPLASH="$OCTOSCREEN_DEFAULT_BOOT_SPLASH"
        echo "OCTOSCREEN_BOOT_SPLASH not set, not a file, or not valid (*.png), using default ( $OCTOSCREEN_BOOT_SPLASH )"
    fi
    if ! [[ -f "$OCTOSCREEN_BOOT_SPLASH" ]]; then
        echo "OCTOSCREEN_BOOT_SPLASH ( $OCTOSCREEN_BOOT_SPLASH ) not a file, exiting."
        exit 1
    fi
    echo "Using Frame Buffer ( $FRAME_BUFFER ) and image ( $OCTOSCREEN_BOOT_SPLASH )"
    /usr/bin/fbi --noverbose --autodown -d $FRAME_BUFFER "$OCTOSCREEN_BOOT_SPLASH"
fi
