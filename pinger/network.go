package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-ping/ping"
	"log/slog"
	"net/http"
	"time"
)

func pingIP(ip string, timeout time.Duration) (float64, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return 0, fmt.Errorf("pinger.pingIP: failed to initialize pinger for IP %s: %w", ip, err)
	}
	pinger.SetPrivileged(true)

	pinger.Count = 3
	pinger.Timeout = timeout
	pinger.Interval = 50 * time.Millisecond

	if err = pinger.Run(); err != nil {
		return 0, fmt.Errorf("pinger.pingIP: failed to run pinger for IP %s: %w", ip, err)
	}
	stats := pinger.Statistics()

	slog.Debug("pinger.pingIP: ping stats",
		slog.String("ip_address", ip),
		slog.Int("sent_packets", stats.PacketsSent),
		slog.Int("received_packets", stats.PacketsRecv),
		slog.Duration("avg_rtt", stats.AvgRtt),
	)

	if stats.PacketsRecv == 0 {
		slog.Debug("pinger.pingIP: No packets received",
			slog.String("ip_address", ip),
			slog.Int("sent_packets", stats.PacketsSent),
			slog.Int("received_packets", stats.PacketsRecv),
		)
		return 0, fmt.Errorf("pinger.pingIP: No received packets from %s", ip)
	}
	fmt.Println("packets: ", stats.PacketsRecv)
	return float64(stats.AvgRtt.Milliseconds()), nil
}

func sendResult(ctx context.Context, results *[]PingResult, config Config) error {
	jsonData, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("pinger.sendResult: Failed to parse data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.BackendURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("pinger.sendResult: Error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("pinger.sendResult: error sending HTTP request: %w\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pinger.sendResult: unexpected HTTP status: %d", resp.StatusCode)
	}

	return nil
}
