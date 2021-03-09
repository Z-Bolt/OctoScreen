#!/usr/bin/env bash

REPO="thebeline/OctoScreen"
MAIN="IMP/Installer"
RELEASES="https://api.github.com/repos/$REPO/releases/latest"
WGET_RAW="https://github.com/$REPO/raw/$MAIN"
LIBRARIES=("inquirer.bash" "optparse.bash")

SOURCE="${BASH_SOURCE[0]}"; while [ -h "$SOURCE" ]; do DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"; SOURCE="$(readlink "$SOURCE")"; [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"; done
DIR="$( cd -P "$( dirname "$SOURCE" )" >/dev/null 2>&1 && pwd )"
PWD="$(pwd)"

echo "$DIR"
echo "$(pwd)"
echo "$(printf " '%q'" "${@}")"
eval set -- "$(printf " %q" "${@}")"
echo "$(printf " '%q'" "${@}")"
echo "${@}"

echo -e "Loading OctoScreen Manager dependencies, please wait...\n"; CL='\e[1A\e[K';
for LIBRARY in ${LIBRARIES[*]}; do
    [[ -f "$DIR/$LIBRARY" ]] && SOURCE_FILE="$DIR/$LIBRARY" || { [[ -f "$PWD/$LIBRARY" ]] && SOURCE_FILE="$PWD/$LIBRARY"; } || { [[ -f "$PWD/scripts/$LIBRARY" ]] && SOURCE_FILE="$PWD/scripts/$LIBRARY"; }
    if [[ -z "$SOURCE_FILE" ]]; then
        echo -e "${CL}Fetching '{$LIBRARY}'..."; SOURCE=$(wget -qO- "$WGET_RAW/scripts/$LIBRARY");
        if [[ "$?" != "0" ]]; then echo " ERROR Fetching dependency '$LIBRARY'!"; [ -v PS1 ] && return || exit 1; else echo ' SUCCESS'; fi
        SOURCE_FILE=<(echo "$SOURCE"); unset SOURCE;
    fi
    echo -en "${CL}Loading '{$LIBRARY}' from: $SOURCE_FILE..."; source "$SOURCE_FILE";
    if [[ "$?" != "0" ]]; then echo " ERROR Loading dependency '$LIBRARY'!"; [ -v PS1 ] && return || exit 1; else echo ' SUCCESS'; fi
    unset SOURCE_FILE
    sleep 1
done; echo -en "${CL}${CL}";

echo "DONE"
[ -v PS1 ] && return || exit 1

yes_no=( 'yes' 'no' )

arch=$(uname -m)
#if [[ $arch == x86_64 ]]; then
#    releaseURL=$(curl -s "$RELEASES" | grep "browser_download_url.*amd64.deb" | cut -d '"' -f 4)
#elif [[ $arch == aarch64 ]]; then
#    releaseURL=$(curl -s "$RELEASES" | grep "browser_download_url.*arm64.deb" | cut -d '"' -f 4)
if  [[ $arch == arm* ]]; then
    releaseURL=$(curl -s "$RELEASES" | grep "browser_download_url.*armf.deb" | cut -d '"' -f 4)
fi
dependencies="libgtk-3-0 xserver-xorg xinit x11-xserver-utils"
IFS='/' read -ra version <<< "$releaseURL"

echo "Installing OctoScreen "${version[7]}, $arch""

echo "Installing Dependencies ..."
sudo apt -qq update
sudo apt -qq install $dependencies -y

if [ -d "/home/pi/OctoPrint/venv" ]; then
    DIRECTORY="/home/pi/OctoPrint/venv"
elif [ -d "/home/pi/oprint" ]; then
    DIRECTORY="/home/pi/oprint"
else
    echo "Neither /home/pi/OctoPrint/venv nor /home/pi/oprint can be found."
    echo "If your OctoPrint instance is running on a different machine just type - in the following prompt."
    text_input "Please specify OctoPrints full virtualenv path manually (no trailing slash)." DIRECTORY
fi;

if [ $DIRECTORY == "-" ]; then
    echo "Not installing any plugins for remote installation. Please make sure to have Display Layer Progress installed."
elif [ ! -d $DIRECTORY ]; then
    echo "Can't find OctoPrint Installation, please run the script again!"
    exit 1
fi;

#if [ $DIRECTORY != "-" ]; then
#  plugins=( 'Display Layer Progress (mandatory)' 'Filament Manager' 'Preheat Button' 'Enclosure' 'Print Time Genius' 'Ultimaker Format Package' 'PrusaSlicer Thumbnails' )
#  checkbox_input "Which plugins should I install (you can also install them via the Octoprint UI)?" plugins selected_plugins
#  echo "Installing Plugins..."
#
#  if [[ " ${selected_plugins[@]} " =~ "Display Layer Progress (mandatory)" ]]; then
#      "$DIRECTORY"/bin/pip install -q --disable-pip-version-check "https://github.com/OllisGit/OctoPrint-DisplayLayerProgress/releases/latest/download/master.zip"
#  fi;
#  if [[ " ${selected_plugins[@]} " =~ "Filament Manager" ]]; then
#      "$DIRECTORY"/bin/pip install -q --disable-pip-version-check "https://github.com/OllisGit/OctoPrint-FilamentManager/releases/latest/download/master.zip"
#  fi;
#  if [[ " ${selected_plugins[@]} " =~ "Preheat Button" ]]; then
#      "$DIRECTORY"/bin/pip install -q --disable-pip-version-check "https://github.com/marian42/octoprint-preheat/archive/master.zip"
#  fi;
#  if [[ " ${selected_plugins[@]} " =~ "Enclosure" ]]; then
#      "$DIRECTORY"/bin/pip install -q --disable-pip-version-check "https://github.com/vitormhenrique/OctoPrint-Enclosure/archive/master.zip"
#  fi;
#  if [[ " ${selected_plugins[@]} " =~ "Print Time Genius" ]]; then
#      "$DIRECTORY"/bin/pip install -q --disable-pip-version-check "https://github.com/eyal0/OctoPrint-PrintTimeGenius/archive/master.zip"
#  fi;
#  if [[ " ${selected_plugins[@]} " =~ "Ultimaker Format Package" ]]; then
#      "$DIRECTORY"/bin/pip install -q --disable-pip-version-check "https://github.com/jneilliii/OctoPrint-UltimakerFormatPackage/archive/master.zip"
#  fi;
#  if [[ " ${selected_plugins[@]} " =~ "PrusaSlicer Thumbnails" ]]; then
#      "$DIRECTORY"/bin/pip install -q --disable-pip-version-check "https://github.com/jneilliii/OctoPrint-PrusaSlicerThumbnails/archive/master.zip"
#  fi;
#fi;

echo "Installing OctoScreen "${version[7]}, $arch" ..."
cd ~
wget -O octoscreen.deb $releaseURL -q --show-progress

sudo dpkg -i octoscreen.deb

rm octoscreen.deb


list_input "Shall I reboot your Pi now?" yes_no reboot
echo "OctoScreen has been successfully installed! :)"
if [ $reboot == 'yes' ]; then
    sudo reboot
fi