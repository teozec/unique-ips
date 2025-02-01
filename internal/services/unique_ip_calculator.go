package services

import (
	_ "github.com/joho/godotenv/autoload"
)

type UniqueIpCalculator struct {
	logs       chan string
	unique_ips map[string]struct{}
}

func NewUniqueIpCalculator() *UniqueIpCalculator {

	unique_calculator := &UniqueIpCalculator{
		logs:       make(chan string),
		unique_ips: make(map[string]struct{}),
	}

	go unique_calculator.handleLogs()
	return unique_calculator
}

// Add a new IP address to the logs
func (u UniqueIpCalculator) LogIp(ip string) {
	u.logs <- ip
}

// Get the number of unique IP address in the logs
func (u UniqueIpCalculator) GetUniqueIpNumber() int {
	return len(u.unique_ips)
}

// Add all IPs from the logs to the unique_ip set as soon as they arrive in the channel
func (u UniqueIpCalculator) handleLogs() {
	for ip := range u.logs {
		u.unique_ips[ip] = struct{}{}
	}
}
