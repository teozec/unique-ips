package services

import (
	_ "github.com/joho/godotenv/autoload"
)

// Receives IP addresses logs and calculates the number of unique ones.
type UniqueIpCalculator struct {
	logs       chan string         // Channel which receives the logges IPs.
	unique_ips map[string]struct{} // Set containing the unique IPs.
	sync       chan struct{}       // Channel used for synchronizing access to unique_ips
}

func NewUniqueIpCalculator() *UniqueIpCalculator {

	unique_calculator := &UniqueIpCalculator{
		logs:       make(chan string),
		unique_ips: make(map[string]struct{}),
		sync:       make(chan struct{}),
	}

	go unique_calculator.receiveLogs()
	return unique_calculator
}

// Add a new IP address to the logs
func (u UniqueIpCalculator) LogIp(ip string) {
	u.logs <- ip
}

// Get the number of unique IP address in the logs
func (u UniqueIpCalculator) GetUniqueIpNumber() int {
	u.sync <- struct{}{} // Send a message on the channel to tell the handleLogs goroutine that it should wait
	l := len(u.unique_ips)
	u.sync <- struct{}{} // Send a second message to tell the goroutine that it can resume writing
	return l
}

// Add all IPs from the logs to the unique_ip set as soon as they arrive in the channel
func (u UniqueIpCalculator) receiveLogs() {
	for {
		select {
		case <-u.sync: // A message has been sent on the synchronization channel: we need to wait for a second one.
			<-u.sync
		case ip := <-u.logs: // Otherwise, add the logged IP address to the set.
			u.unique_ips[ip] = struct{}{}
		}
	}
}
