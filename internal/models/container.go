package models

import (
	"time"
)

type Container struct {
	ContainerID  int       `db:"container_id"`
	Status       string    `db:"container_status" json:"container_status"`
	Addr         string    `db:"addr" json:"addr"`
	PingDuration float64   `db:"p_duration" json:"p_duration"`
	PingedAt     time.Time `db:"pinged_at" json:"pinged_at"`
}

type ContainersList struct {
	TotalCount int          `json:"total_count"`
	TotalPages int          `json:"total_pages"`
	Page       int          `json:"page"`
	Size       int          `json:"size"`
	HasMore    bool         `json:"has_more"`
	Containers []*Container `json:"containers"`
}
