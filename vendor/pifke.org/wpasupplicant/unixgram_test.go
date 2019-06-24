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

package wpasupplicant

import (
	"bytes"
	"net"
	"testing"
)

var parseScanResultTests = []struct {
	input  string
	expect []*scanResult
}{
	{
		// actual output from wpa_supplicant 2.4-0ubuntu6
		input: "bssid / frequency / signal level / flags / ssid\n" +
			"8a:15:14:8a:46:51\t5560\t-58\t[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][ESS]\tWIP-Backoffice\n" +
			"8a:15:14:8a:46:50\t5560\t-58\t[WPA-PSK-CCMP+TKIP][WPA2-PSK-CCMP+TKIP][ESS]\tWorkInProgressMember\n",
		expect: []*scanResult{
			&scanResult{
				bssid:     net.HardwareAddr{0x8a, 0x15, 0x14, 0x8a, 0x46, 0x51},
				frequency: 5560,
				rssi:      -58,
				flags:     []string{"WPA-PSK-CCMP+TKIP", "WPA2-PSK-CCMP+TKIP", "ESS"},
				ssid:      "WIP-Backoffice",
			},
			&scanResult{
				bssid:     net.HardwareAddr{0x8a, 0x15, 0x14, 0x8a, 0x46, 0x50},
				frequency: 5560,
				rssi:      -58,
				flags:     []string{"WPA-PSK-CCMP+TKIP", "WPA2-PSK-CCMP+TKIP", "ESS"},
				ssid:      "WorkInProgressMember",
			},
		},
	}, {
		// reordered/added/missing columns from some theoretical
		// future version of wpa_supplicant
		input: "frequency / bssid / foobar / ssid\n" +
			"5560\t8a:15:14:8a:46:51\thello\tWIP-Backoffice\n" +
			"5560\t8a:15:14:8a:46:50\tgoodbye\tWorkInProgressMember\n",
		expect: []*scanResult{
			&scanResult{
				bssid:     net.HardwareAddr{0x8a, 0x15, 0x14, 0x8a, 0x46, 0x51},
				frequency: 5560,
				ssid:      "WIP-Backoffice",
			},
			&scanResult{
				bssid:     net.HardwareAddr{0x8a, 0x15, 0x14, 0x8a, 0x46, 0x50},
				frequency: 5560,
				ssid:      "WorkInProgressMember",
			},
		},
	},
}

func TestParseScanResults(t *testing.T) {
	for _, test := range parseScanResultTests {
		output, errs := parseScanResults(bytes.NewBufferString(test.input))
		if len(errs) > 0 {
			t.Error("errors parsing scan results")
		}

		if len(output) != len(test.expect) {
			t.Errorf("wrong number of results (got %d, expect %d)", len(output), len(test.expect))
		}

		for i := range output {
			if test.expect[i].bssid != nil {
				if bytes.Compare(output[i].BSSID(), test.expect[i].bssid) != 0 {
					t.Errorf("wrong bssid (got %q, expect %q)", output[i].BSSID(), test.expect[i].bssid)
				}
			}
			if test.expect[i].frequency != 0 {
				if output[i].Frequency() != test.expect[i].frequency {
					t.Errorf("wrong frequency (got %d, expect %d)", output[i].Frequency(), test.expect[i].frequency)
				}
			}
			if test.expect[i].rssi != 0 {
				if output[i].RSSI() != test.expect[i].rssi {
					t.Errorf("wrong rssi (got %d, expect %d)", output[i].RSSI(), test.expect[i].rssi)
				}
			}
			if test.expect[i].ssid != "" {
				if output[i].SSID() != test.expect[i].ssid {
					t.Errorf("wrong rssi (got %s, expect %s)", output[i].SSID(), test.expect[i].ssid)
				}
			}

			flags := output[i].Flags()
			flagsMatch := true
			if len(test.expect[i].flags) != len(flags) {
				flagsMatch = false
			} else {
				for j := range test.expect[i].flags {
					if flags[j] != test.expect[i].flags[j] {
						flagsMatch = false
						break
					}
				}
			}
			if !flagsMatch {
				t.Errorf("got flags %q, expected %q", flags, test.expect[i].flags)
			}
		}
	}
}

func TestParseStatusResults(t *testing.T) {
	testData := "bssid=02:00:01:02:03:04\n" +
		"ssid=test network\n" +
		"pairwise_cipher=CCMP\n" +
		"group_cipher=CCMP\n" +
		"key_mgmt=WPA-PSK\n" +
		"wpa_state=COMPLETED\n" +
		"ip_address=192.168.1.21\n" +
		"Supplicant PAE state=AUTHENTICATED\n" +
		"suppPortStatus=Authorized\n" +
		"EAP state=SUCCESS"

	res, err := parseStatusResults(bytes.NewBufferString(testData))
	if err != nil {
		t.Errorf("Error parsing status result %t", err)
	}

	if res.WPAState() != "COMPLETED" {
		t.Errorf("WPAState was not COMPLETED. Was %s", res.WPAState())
	}

	if res.IPAddr() != "192.168.1.21" {
		t.Errorf("IPAddr was not 192.168.1.21. Was %s", res.IPAddr())
	}

	if res.KeyMgmt() != "WPA-PSK" {
		t.Errorf("KeyMgmt was not WPA-PSK. Was %s", res.KeyMgmt())
	}

	if res.Address() != "" {
		t.Errorf("Address should be empty. Was %s", res.Address())
	}
}
