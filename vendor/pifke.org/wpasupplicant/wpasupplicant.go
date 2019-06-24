// Copyright (c) 2017 Dave Pifke.
//
// Redistribution and use in source and binary forms, with or without
// modification, is permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

// Package wpasupplicant provides an interface for talking to the
// wpa_supplicant daemon.
//
// At the moment, this simply provides an interface for fetching wifi scan
// results.  More functionality is (probably) coming soon.
package wpasupplicant

import (
	"net"
)

// Cipher is one of the WPA_CIPHER constants from the wpa_supplicant source.
type Cipher int

const (
	CIPHER_NONE Cipher = 1 << iota
	WEP40
	WEP104
	TKIP
	CCMP
	AES_128_CMAC
	GCMP
	SMS4
	GCMP_256
	CCMP_256
	_
	BIP_GMAC_128
	BIP_GMAC_256
	BIP_CMAC_256
	GTK_NOT_USED
)

// KeyMgmt is one of the WPA_KEY_MGMT constants from the wpa_supplicant
// source.
type KeyMgmt int

const (
	IEEE8021X KeyMgmt = 1 << iota
	PSK
	KEY_MGMT_NONE
	IEEE8021X_NO_WPA
	WPA_NONE
	FT_IEEE8021X
	FT_PSK
	IEEE8021X_SHA256
	PSK_SHA256
	WPS
	SAE
	FT_SAE
	WAPI_PSK
	WAPI_CERT
	CCKM
	OSEN
	IEEE8021X_SUITE_B
	IEEE8021X_SUITE_B_192
)

type Algorithm int

// ScanResult is a scanned BSS.
type ScanResult interface {
	// BSSID is the MAC address of the BSS.
	BSSID() net.HardwareAddr

	// SSID is the SSID of the BSS.
	SSID() string

	// Frequency is the frequency, in Mhz, of the BSS.
	Frequency() int

	// RSSI is the received signal strength, in dB, of the BSS.
	RSSI() int

	// Flags is an array of flags, in string format, returned by the
	// wpa_supplicant SCAN_RESULTS command.  Future versions of this code
	// will parse these into something more meaningful.
	Flags() []string
}

// scanResult is a package-private implementation of ScanResult.
type scanResult struct {
	bssid     net.HardwareAddr
	ssid      string
	frequency int
	rssi      int
	flags     []string
}

func (r *scanResult) BSSID() net.HardwareAddr { return r.bssid }
func (r *scanResult) SSID() string            { return r.ssid }
func (r *scanResult) Frequency() int          { return r.frequency }
func (r *scanResult) RSSI() int               { return r.rssi }
func (r *scanResult) Flags() []string         { return r.flags }

// ConfiguredNetwork is a configured network (from LIST_NETWORKS)
type ConfiguredNetwork interface {
	NetworkID() string
	SSID() string
	BSSID() string
	Flags() []string
}

type configuredNetwork struct {
	networkID string
	ssid      string
	bssid     string // Since bssid can be any
	flags     []string
}

func (r *configuredNetwork) NetworkID() string { return r.networkID }
func (r *configuredNetwork) BSSID() string     { return r.bssid }
func (r *configuredNetwork) SSID() string      { return r.ssid }
func (r *configuredNetwork) Flags() []string   { return r.flags }

type StatusResult interface {
	WPAState() string
	KeyMgmt() string
	IPAddr() string
	SSID() string
	Address() string
}

type statusResult struct {
	wpaState string
	keyMgmt  string
	ipAddr   string
	ssid     string
	address  string
}

func (s *statusResult) WPAState() string { return s.wpaState }
func (s *statusResult) KeyMgmt() string  { return s.keyMgmt }
func (s *statusResult) IPAddr() string   { return s.ipAddr }
func (s *statusResult) SSID() string     { return s.ssid }
func (s *statusResult) Address() string  { return s.address }

type WPAEvent struct {
	Event     string
	Arguments map[string]string
	Line      string
}

// Conn is a connection to wpa_supplicant over one of its communication
// channels.
type Conn interface {
	// Close closes the unixgram connection
	Close() error

	// Ping tests the connection.  It returns nil if wpa_supplicant is
	// responding.
	Ping() error

	// AddNetwork creates an empty network configuration. Returns the network
	// ID.
	AddNetwork() (int, error)

	// SetNetwork configures a network property. Returns error if the property
	// configuration failed.
	SetNetwork(int, string, string) error

	// EnableNetwork enables a network. Returns error if the command fails.
	EnableNetwork(int) error

	// EnableAllNetworks enables all configured networks. Returns error if the command fails.
	EnableAllNetworks() error

	// SelectNetwork selects a network (and disables the others).
	SelectNetwork(int) error

	// DisableNetwork disables a network.
	DisableNetwork(int) error

	// RemoveNetwork removes a network from the configuration.
	RemoveNetwork(int) error

	// RemoveAllNetworks removes all networks (basically running `REMOVE_NETWORK all`).
	// Returns error if command fails.
	RemoveAllNetworks() error

	// SaveConfig stores the current network configuration to disk.
	SaveConfig() error

	// Reconfigure sends a RECONFIGURE command to the wpa_supplicant. Returns error when
	// command fails.
	Reconfigure() error

	// Reassociate sends a REASSOCIATE command to the wpa_supplicant. Returns error when
	// command fails.
	Reassociate() error

	// Reconnect sends a RECONNECT command to the wpa_supplicant. Returns error when
	// command fails.
	Reconnect() error

	// ListNetworks returns the currently configured networks.
	ListNetworks() ([]ConfiguredNetwork, error)

	// Status returns current wpa_supplicant status
	Status() (StatusResult, error)

	// Scan triggers a new scan. Returns error if the wpa_supplicant does not
	// return OK.
	Scan() error

	// ScanResult returns the latest scanning results.  It returns a slice
	// of scanned BSSs, and/or a slice of errors representing problems
	// communicating with wpa_supplicant or parsing its output.
	ScanResults() ([]ScanResult, []error)

	EventQueue() chan WPAEvent
}
