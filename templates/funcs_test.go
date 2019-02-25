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
		result, err := ip4(tc.ip, tc.idx)
		if err != nil {
			t.Error(err)
		}

		if result != tc.expected {
			t.Errorf("expected to get '%s', instead got '%s'\n", tc.expected, result)
		}
	}
}

func TestIP4Mask(t *testing.T) {
	var ip4MaskTescCases = []struct {
		ip       string
		expected string
	}{
		{"10.0.0.0/24", "255.255.255.0"},
		{"10.0.0.0/32", "255.255.255.255"},
		{"10.0.0.0/26", "255.255.255.192"},
	}

	for _, tc := range ip4MaskTescCases {
		result, err := ip4mask(tc.ip)
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
		result := split(tc.str, tc.sep, tc.idx)
		if result != tc.expected {
			t.Errorf("expected to get '%s', instead got '%s'\n", tc.expected, result)
		}
	}
}
