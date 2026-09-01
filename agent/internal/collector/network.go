package collector

import (
	"time"

	gnet "github.com/shirou/gopsutil/v4/net"
)

var previousSent uint64
var previousRecv uint64
var previousTime time.Time

func GetNetworkSpeed() (float64, float64, error) {

	stats, err := gnet.IOCounters(false)
	if err != nil {
		return 0, 0, err
	}

	now := time.Now()

	if previousTime.IsZero() {
		previousSent = stats[0].BytesSent
		previousRecv = stats[0].BytesRecv
		previousTime = now
		return 0, 0, nil
	}

	elapsed := now.Sub(previousTime).Seconds()
	if elapsed <= 0 {
		return 0, 0, nil
	}

	upload := float64(stats[0].BytesSent-previousSent) * 8 / elapsed / 1000000
	download := float64(stats[0].BytesRecv-previousRecv) * 8 / elapsed / 1000000

	previousSent = stats[0].BytesSent
	previousRecv = stats[0].BytesRecv
	previousTime = now

	return upload, download, nil
}
