# pifke.org/wpasupplicant

[![GoDoc](https://godoc.org/pifke.org/wpasupplicant?status.svg)](https://godoc.org/pifke.org/wpasupplicant)
[![Build Status](https://api.travis-ci.org/dpifke/golang-wpasupplicant.svg)](https://travis-ci.org/dpifke/golang-wpasupplicant)
[![Test Coverage](https://coveralls.io/repos/github/dpifke/golang-wpasupplicant/badge.svg)](https://coveralls.io/github/dpifke/golang-wpasupplicant)

Golang interface for talking to wpa_supplicant.

At the moment, this simply provides an interface for fetching wifi scan
results.  More functionality (probably) coming soon.

## Example

```
import (
	"fmt"

	"pifke.org/wpasupplicant"
)

// Prints the BSSID (MAC address) and SSID of each access point in range:
w, err := wpasupplicant.Unixgram("wlan0")
if err != nil {
	panic(err)
}
for _, bss := range w.ScanResults() {
	fmt.Fprintf("%s\t%s\n", bss.BSSID(), bss.SSID())
}
```

## Downloading

If you use this library in your own code, please use the canonical URL in your
Go code, instead of Github:

```
go get pifke.org/wpasupplicant
```

Or (until I finish setting up the self-hosted repository):

```
# From the root of your project:
git submodule add https://github.com/dpifke/golang-wpasupplicant vendor/pifke.org/wpasupplicant
```

Then:

```
import (
        "pifke.org/wpasupplicant"
)
```

As opposed to the pifke.org URL, I make no guarantee this Github repository
will exist or be up-to-date in the future.

## Documentation

Available on [godoc.org](https://godoc.org/pifke.org/wpasupplicant).

## License

Three-clause BSD.  See LICENSE.txt.

Contact me if you want to use this code under different terms.

## Author

Dave Pifke.  My email address is my first name "at" my last name "dot org."

I'm [@dpifke](https://twitter.com/dpifke) on Twitter.  My PGP key
is available on [Keybase](https://keybase.io/dpifke).
