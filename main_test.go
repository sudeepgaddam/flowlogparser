package main

import (
	"testing"
)

// TestParseFlowLog tests the parsing of flow logs and tagging logic.
func TestParseFlowLog(t *testing.T) {
	// Define a mock lookup table
	lookupTable := LookupTable{
		"25_tcp":  "sv_P1",
		"443_tcp": "sv_P2",
		"23_tcp":  "sv_P1",
		"53_udp":  "dns",
		"110_tcp": "email",
		"993_tcp": "email",
		"80_tcp":  "web",
		"22_tcp":  "ssh",
	}

	// Define test cases
	tests := []struct {
		line     string
		expected string
		error    bool
	}{
		// Test case 1: Exact match (port 443, protocol TCP)
		{"2 123456789012 eni-0a1b2c3d 10.0.1.201 198.51.100.2 443 49153 6 25 20000 1620140761 1620140821 ACCEPT OK", "sv_P2", false},
		// Test case 2: Match with UDP (port 53, protocol UDP)
		{"2 123456789012 eni-5e6f7g8h 192.168.1.101 198.51.100.3 53 49155 17 10 8000 1620140761 1620140821 ACCEPT OK", "dns", false},
		// Test case 3: No match (Untagged)
		{"2 123456789012 eni-4j5k6l7m 10.0.0.1 192.0.2.5 9999 12345 6 5 10000 1620140761 1620140821 ACCEPT OK", "Untagged", false},
		// Test case 4: SSH traffic (port 22, protocol TCP)
		{"2 123456789012 eni-1a2b3c4d 203.0.113.12 192.168.0.1 22 49158 6 12 14000 1620140761 1620140821 ACCEPT OK", "ssh", false},
		// Test case 5: Web traffic (port 80, protocol TCP)
		{"2 123456789012 eni-6m7n8o9p 10.0.2.200 198.51.100.4 80 49158 6 18 14000 1620140761 1620140821 ACCEPT OK", "web", false},
		// Test case 6: Email traffic (IMAP, port 993, protocol TCP)
		{"2 123456789012 eni-7i8j9k0l 172.16.0.101 192.0.2.203 993 49157 6 8 5000 1620140761 1620140821 ACCEPT OK", "email", false},
		// Test case 7: Invalid protocol (unknown protocol number 99)
		{"2 123456789012 eni-1b2c3d4e 10.0.3.200 198.51.100.5 53 49153 99 18 14000 1620140761 1620140821 ACCEPT OK", "Untagged", false},
		// Test case 8: Invalid port abc given
		{"2 123456789012 eni-1b2c3d4e 10.0.3.200 198.51.100.5 53 49153 abc 18 14000 1620140761 1620140821 ACCEPT OK", "", true},
	}

	for _, test := range tests {
		// Parse the flow log line
		flowLog, err := ParseFlowLog(test.line, lookupTable)
		if test.error {
			if err == nil {
				t.Errorf("For flow log: %s, expected error but got no error", test.line)
			}
		} else if flowLog.Tag != test.expected {
			// Compare the resulting tag with the expected tag
			t.Errorf("For flow log: %s, expected tag %s, got %s", test.line, test.expected, flowLog.Tag)
		}
	}
}

// TestCountTags tests the tag counting logic.
func TestCountTags(t *testing.T) {
	flowLogs := []FlowLog{
		{DstPort: "443", Protocol: "tcp", Tag: "sv_P2"},
		{DstPort: "53", Protocol: "udp", Tag: "dns"},
		{DstPort: "9999", Protocol: "tcp", Tag: "Untagged"},
		{DstPort: "22", Protocol: "tcp", Tag: "ssh"},
		{DstPort: "993", Protocol: "tcp", Tag: "email"},
	}

	expectedCounts := map[string]int{
		"sv_P2":    1,
		"dns":      1,
		"Untagged": 1,
		"ssh":      1,
		"email":    1,
	}

	tagCount := CountTags(flowLogs)

	for tag, expectedCount := range expectedCounts {
		if tagCount[tag] != expectedCount {
			t.Errorf("For tag %s, expected count %d, got %d", tag, expectedCount, tagCount[tag])
		}
	}
}

// TestCountPortProtocol tests the port/protocol counting logic.
func TestCountPortProtocol(t *testing.T) {
	flowLogs := []FlowLog{
		{DstPort: "443", Protocol: "tcp", Tag: "sv_P2"},
		{DstPort: "53", Protocol: "udp", Tag: "dns"},
		{DstPort: "9999", Protocol: "tcp", Tag: "Untagged"},
		{DstPort: "22", Protocol: "tcp", Tag: "ssh"},
		{DstPort: "993", Protocol: "tcp", Tag: "email"},
	}

	expectedCounts := map[string]int{
		"443_tcp":  1,
		"53_udp":   1,
		"9999_tcp": 1,
		"22_tcp":   1,
		"993_tcp":  1,
	}

	portProtocolCount := CountPortProtocol(flowLogs)

	for key, expectedCount := range expectedCounts {
		if portProtocolCount[key] != expectedCount {
			t.Errorf("For port/protocol %s, expected count %d, got %d", key, expectedCount, portProtocolCount[key])
		}
	}
}
