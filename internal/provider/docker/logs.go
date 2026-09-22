package docker

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"
)

func (p *Provider) Logs(ctx context.Context, target string, limit int) ([]string, error) {
	if limit <= 0 {
		return []string{}, nil
	}
	logs, err := p.client.ContainerLogs(ctx, target, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       strconv.Itoa(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("get logs for container %q: %w", target, err)
	}
	defer logs.Close()

	inspectResult, err := p.client.ContainerInspect(ctx, target, client.ContainerInspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("inspect container %q for log configuration: %w", target, err)
	}

	tty := false
	if inspectResult.Container.Config != nil {
		tty = inspectResult.Container.Config.Tty
	}

	var buf bytes.Buffer
	if tty {
		_, err := io.Copy(&buf, logs)
		if err != nil {
			return nil, fmt.Errorf("read tty logs %s failed: %w", target, err)
		}
	} else {
		_, err := stdcopy.StdCopy(&buf, &buf, logs)
		if err != nil {
			return nil, err
		}
	}

	scanner := bufio.NewScanner(&buf)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
