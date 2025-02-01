package services

import (
	"testing"
)

func TestUniqueIpCalculator(t *testing.T) {
	u := NewUniqueIpCalculator()
	var count int

	count = u.GetUniqueIpNumber()
	if count != 0 {
		t.Fatalf("expected 0 unique ips, got %d", count)
	}

	u.LogIp("192.168.1.1")
	u.LogIp("192.168.1.2")
	u.LogIp("192.168.1.3")

	count = u.GetUniqueIpNumber()
	if count != 3 {
		t.Fatalf("expected 3 unique ips, got %d", count)
	}

	u.LogIp("192.168.1.1")
	u.LogIp("192.168.1.2")
	u.LogIp("192.168.1.4")

	count = u.GetUniqueIpNumber()
	if count != 4 {
		t.Fatalf("expected 4 unique ips, got %d", count)
	}
}
