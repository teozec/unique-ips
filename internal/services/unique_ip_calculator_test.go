package services

import (
	"testing"
)

func TestUniqueIpCalculator(t *testing.T) {
	u := NewUniqueIpCalculator()
	var count int

	count = u.GetUniqueIpNumber()
	if count != 0 {
		t.Errorf("expected 0 unique ips, got %d", count)
	}

	u.LogIp("192.168.1.1")
	u.LogIp("192.168.1.2")
	u.LogIp("192.168.1.3")

	count = u.GetUniqueIpNumber()
	if count != 3 {
		t.Errorf("expected 3 unique ips, got %d", count)
	}

	u.LogIp("192.168.1.1")
	u.LogIp("192.168.1.2")
	u.LogIp("192.168.1.4")
	u.LogIp("192.168.1.5")

	count = u.GetUniqueIpNumber()
	if count != 5 {
		t.Errorf("expected 5 unique ips, got %d", count)
	}
}
