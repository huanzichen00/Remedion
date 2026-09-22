package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/huanzichen00/remedion/internal/observe"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func (p *Provider) Stats(ctx context.Context, target string) (observe.Metrics, error) {
	result, err := p.client.ContainerStats(ctx, target, client.ContainerStatsOptions{
		Stream:                false,
		IncludePreviousSample: true,
	})
	if err != nil {
		return observe.Metrics{}, fmt.Errorf("get stats for container %q: %w", target, err)
	}
	defer result.Body.Close()

	var stats container.StatsResponse

	err = json.NewDecoder(result.Body).Decode(&stats)
	if err != nil {
		return observe.Metrics{}, fmt.Errorf("decode stats for contianer %q: %w", target, err)
	}

	memUsage := stats.MemoryStats.Usage
	memLimit := stats.MemoryStats.Limit
	var memPercent float64
	if memLimit > 0 {
		memPercent = float64(memUsage) / float64(memLimit) * 100
	}

	cpuPercent := calculateCPUPercent(stats)

	metrics := observe.Metrics{
		CPUPercent:    cpuPercent,
		MemoryUsage:   memUsage,
		MemoryLimit:   memLimit,
		MemoryPercent: memPercent,
		PIDs:          stats.PidsStats.Current,
	}

	return metrics, nil
}

func calculateCPUPercent(stats container.StatsResponse) float64 {
	currentCPU := stats.CPUStats.CPUUsage.TotalUsage
	previousCPU := stats.PreCPUStats.CPUUsage.TotalUsage

	currentSystem := stats.CPUStats.SystemUsage
	previousSystem := stats.PreCPUStats.SystemUsage

	if currentCPU <= previousCPU || currentSystem <= previousSystem {
		return 0
	}

	cpuDelta := currentCPU - previousCPU
	systemDelta := currentSystem - previousSystem

	cpus := stats.CPUStats.OnlineCPUs
	if cpus == 0 {
		cpus = uint32(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}

	if cpus == 0 {
		return 0
	}

	return float64(cpuDelta) /
		float64(systemDelta) *
		float64(cpus) *
		100
}
