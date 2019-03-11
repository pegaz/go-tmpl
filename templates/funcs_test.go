// Copyright © 2019 Pawel Potrykus <pawel.potrykus@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package templates

import "testing"

func TestIP4(t *testing.T) {
	var ip4TestCases = []struct {
		ip       string
		idx      int
		expected string
	}{
		{"10.0.0.1/24", 0, "10.0.0.0"},
		{"10.0.0.1", 0, "10.0.0.1"},
		{"10.0.0.2", 1, "10.0.0.3"},
		{"10.0.0.0", 256, "10.0.1.0"},
		{"10.0.0.0/32", 2, "10.0.0.2"},
	}

	for _, tc := range ip4TestCases {
		result, err := IP4(tc.ip, tc.idx)
		if err != nil {
			t.Error(err)
		}

		if result != tc.expected {
			t.Errorf("expected to get '%s', instead got '%s'\n", tc.expected, result)
		}
	}
}

func TestIP4Mask(t *testing.T) {
	var IP4MaskTestCases = []struct {
		ip       string
		expected string
	}{
		{"10.0.0.0/24", "255.255.255.0"},
		{"10.0.0.0/32", "255.255.255.255"},
		{"10.0.0.0/26", "255.255.255.192"},
	}

	for _, tc := range IP4MaskTestCases {
		result, err := IP4Mask(tc.ip)
		if err != nil {
			t.Error(err)
		}

		if result != tc.expected {
			t.Errorf("expected to get '%s', instead got '%s'\n", tc.expected, result)
		}
	}
}

func TestSplit(t *testing.T) {
	var splitTestCases = []struct {
		str      string
		sep      string
		idx      int
		expected string
	}{
		{"10.0.0.0/24", "/", 0, "10.0.0.0"},
		{"10.0.0.0/24", "/", 1, "24"},
		{"10.0.0.0/24", "/", 2, ""},
	}

	for _, tc := range splitTestCases {
		result := Split(tc.str, tc.sep, tc.idx)
		if result != tc.expected {
			t.Errorf("expected to get '%s', instead got '%s'\n", tc.expected, result)
		}
	}
}

func TestIP4Cidr(t *testing.T) {
	var testCases = []struct {
		ip       string
		expected string
	}{
		{"192.168.0.0/24", "24"},
		{"10.0.0.0/8", "8"},
		{"10.0.1.0/29", "29"},
		{"10.0.2.0/32", "32"},
	}

	for _, tc := range testCases {
		cidr, err := IP4Cidr(tc.ip)
		if err != nil {
			t.Error(err)
		}

		if cidr != tc.expected {
			t.Errorf("expected to get CIDR %s, instead got %s", tc.expected, cidr)
		}
	}
}

func TestIP4CidrToMask(t *testing.T) {
	var testCases = []struct {
		cidr     string
		expected string
	}{
		{"24", "255.255.255.0"},
		{"28", "255.255.255.240"},
		{"32", "255.255.255.255"},
	}

	for _, tc := range testCases {
		mask, err := IP4CidrToMask(tc.cidr)
		if err != nil {
			t.Error(err)
		}

		if mask != tc.expected {
			t.Errorf("expected to get mask %s, instead got %s", tc.expected, mask)
		}
	}
}

func TestIP4MaskToCidr(t *testing.T) {
	var testCases = []struct {
		mask     string
		expected string
	}{
		{"255.255.255.0", "24"},
		{"255.255.255.240", "28"},
		{"255.255.255.255", "32"},
	}

	for _, tc := range testCases {
		cidr, err := IP4MaskToCidr(tc.mask)
		if err != nil {
			t.Error(err)
		}

		if cidr != tc.expected {
			t.Errorf("expected to get CIDR %s, instead got %s", tc.expected, cidr)
		}
	}
}
