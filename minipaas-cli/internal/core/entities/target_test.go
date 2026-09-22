package entities

import "testing"

func TestParseServiceTarget(t *testing.T) {
	cases := []struct {
		in      string
		service string
		port    string
	}{
		{"api", "api", "80"},
		{"api:8080", "api", "8080"},
		{"example.com:5000/api:8080", "example.com:5000/api", "8080"},
	}

	for _, c := range cases {
		service, port := ParseServiceTarget(c.in)
		if service != c.service || port != c.port {
			t.Fatalf("ParseServiceTarget(%q) = (%q, %q), want (%q, %q)", c.in, service, port, c.service, c.port)
		}
	}
}
