package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log/slog"
	"net"
	"time"
)

type PingResult struct {
	Status       string    `json:"container_status"`
	Addr         string    `json:"addr"`
	PingDuration float64   `json:"p_duration"`
	PingedAt     time.Time `json:"pinged_at"`
}

func (pr *PingResult) Log() {
	slog.Info("Ping result",
		slog.String("module", "pinger.eachContainers"),
		slog.String("ip_address", pr.Addr),
		slog.String("status", pr.Status),
		slog.Float64("ping_time_ms", pr.PingDuration),
		slog.String("timestamp", pr.PingedAt.Format(time.RFC3339)),
	)
}

func fetchContainers(ctx context.Context, client *client.Client) ([]types.ContainerJSON, error) {
	c, err := client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}
	var containers = make([]types.ContainerJSON, 0, len(c))

	for _, cont := range c {
		contJSON, err := client.ContainerInspect(ctx, cont.ID)
		if err != nil {
			slog.Error("Failed to inspect container",
				slog.String("module", "pinger.fetchContainers"),
				slog.String("id", cont.ID),
				slog.String("error", err.Error()),
			)
			continue
		}
		containers = append(containers, contJSON)
	}
	return containers, nil
}

func eachContainers(containers []types.ContainerJSON, timeout time.Duration) *[]PingResult {
	results := make([]PingResult, 0, len(containers))
	for _, c := range containers {
		var result = PingResult{}

		slog.Debug("Check container... \n",
			slog.String("module", "pinger.eachContainers"),
			slog.String("id", c.ID),
			slog.String("status", c.State.Status),
			slog.String("image", c.Image),
			slog.String("created", c.Created))

		result.Status = c.State.Status

		for network, ip := range c.NetworkSettings.Networks {
			if net.ParseIP(ip.IPAddress) == nil {
				slog.Debug("The container does not have an IP address on the network",
					slog.String("module", "pinger.eachContainers"),
					slog.String("id", c.ID),
					slog.String("name", c.Name),
					slog.String("network", network))
				continue
			}

			result.Addr = ip.IPAddress
			rtt, err := pingIP(ip.IPAddress, timeout)
			if err != nil {
				slog.Error("Ping failed",
					slog.String("module", "pinger.eachContainers"),
					slog.String("id", c.ID),
					slog.String("name", c.Name),
					slog.String("network", network),
					slog.String("error", err.Error()))
				result.PingDuration = -1
				continue
			}
			result.PingDuration = rtt
			result.PingedAt = time.Now()
		}

		result.Log()

		results = append(results, result)
	}
	return &results
}
