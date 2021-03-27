#!/usr/bin/env bash

OCTOSCREEN_DEFAULT_BOOT_SPLASH="/opt/octoscreen/styles/z-bolt/images/boot-splash/"
FRAME_BUFFER=$1

if [[ -z "$FRAME_BUFFER" ]]; then
    echo "Unable to detect valid Frame Buffer for OctoScreen Boot Splash"
else
    # PARSE THE USER INPUT
    if [[ -z "$OCTOSCREEN_BOOT_SPLASH" ]]; then
        OCTOSCREEN_BOOT_SPLASH="$OCTOSCREEN_DEFAULT_BOOT_SPLASH"
        echo "OCTOSCREEN_BOOT_SPLASH not set, using default ( $OCTOSCREEN_BOOT_SPLASH )"
    elif [[ -f "$OCTOSCREEN_BOOT_SPLASH" ]]; then
        if [[ ! $OCTOSCREEN_BOOT_SPLASH =~ .*\.png$ ]]; then
            OCTOSCREEN_BOOT_SPLASH="$OCTOSCREEN_DEFAULT_BOOT_SPLASH"
            echo "OCTOSCREEN_BOOT_SPLASH is not a valid PNG file, using default ( $OCTOSCREEN_BOOT_SPLASH )"
        fi
    elif [[ -d "$OCTOSCREEN_BOOT_SPLASH" ]]; then
        if [[ "$( ls -A ${OCTOSCREEN_BOOT_SPLASH}/*.png )" ]]; then
            OCTOSCREEN_BOOT_SPLASH="${OCTOSCREEN_BOOT_SPLASH}/*.png"
        else
            OCTOSCREEN_BOOT_SPLASH="$OCTOSCREEN_DEFAULT_BOOT_SPLASH"
            echo "OCTOSCREEN_BOOT_SPLASH is a directory, but doesn't contain any PNG files, using default ( $OCTOSCREEN_BOOT_SPLASH )"
        fi
    else
        OCTOSCREEN_BOOT_SPLASH="$OCTOSCREEN_DEFAULT_BOOT_SPLASH"
        echo "OCTOSCREEN_BOOT_SPLASH neither a file or directory, using default ( $OCTOSCREEN_BOOT_SPLASH )"
    fi
    
    # FINAL VALIDATION OF PATH
    if [[ ! $OCTOSCREEN_BOOT_SPLASH =~ .*\.png$ ]]; then
        if [[ -d "$OCTOSCREEN_BOOT_SPLASH" ]] && [[ "$( ls -A ${OCTOSCREEN_BOOT_SPLASH}/*.png )" ]]; then
            OCTOSCREEN_BOOT_SPLASH="${OCTOSCREEN_BOOT_SPLASH}/*.png"
            echo "OCTOSCREEN_BOOT_SPLASH is a directory with PNGs, globbing ( $OCTOSCREEN_BOOT_SPLASH )"
        else
            echo "OCTOSCREEN_BOOT_SPLASH is neither a PNG file, or directory with PNG files. Exiting"
            exit 1
        fi
    fi
    
    if [[ ! "$( ls -A $OCTOSCREEN_BOOT_SPLASH )" ]]; then
        echo "Can not resolve any files from supplied path ( $OCTOSCREEN_BOOT_SPLASH ), exiting."
        exit 1;
    fi
    
    echo "Using Frame Buffer ( $FRAME_BUFFER ) and image/directory ( $OCTOSCREEN_BOOT_SPLASH )"
    /usr/bin/fbi --noverbose --autodown -t 1 -d $FRAME_BUFFER $OCTOSCREEN_BOOT_SPLASH
fi
