package installer

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	cp "github.com/otiai10/copy"

	"go.uber.org/zap"
)

const (
	webURL       = "http://unix/web/"
	completeURL  = "http://unix/Startup/Complete"
	maxAttempts  = 20
	attemptDelay = 10 * time.Second
	logTailLines = 40
	fatalMarker  = "Unhandled Exception"
)

type Jellyfin struct {
	appDir, dataDir string
	client          *http.Client
	executor        *Executor
	logger          *zap.Logger
}

func NewJellyfin(appDir, dataDir string, executor *Executor, logger *zap.Logger) *Jellyfin {
	return &Jellyfin{
		appDir:  appDir,
		dataDir: dataDir,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", path.Join(dataDir, "socket"))
				},
			},
		},
		executor: executor,
		logger:   logger,
	}
}

func (j *Jellyfin) Complete() error {
	err := j.waitForWeb()
	if err != nil {
		return err
	}
	return j.completeStartup()
}

func (j *Jellyfin) waitForWeb() error {
	lastResponse := "not attempted"
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(attemptDelay)
		}

		resp, err := j.client.Get(webURL)
		if err != nil {
			lastResponse = err.Error()
		} else {
			status := resp.StatusCode
			resp.Body.Close()
			if status == http.StatusOK {
				return nil
			}
			lastResponse = fmt.Sprintf("HTTP %d", status)
		}

		serverLog := j.ServerLog()
		if strings.Contains(serverLog, fatalMarker) {
			return fmt.Errorf("server failed to start, last web response: %s\n%s", lastResponse, serverLog)
		}
	}

	return fmt.Errorf("web endpoint not available after %d attempts, last web response: %s\n%s",
		maxAttempts, lastResponse, j.ServerLog())
}

func (j *Jellyfin) completeStartup() error {
	lastResponse := "not attempted"
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(attemptDelay)
		}

		resp, err := j.client.Post(completeURL, "application/json", nil)
		if err != nil {
			lastResponse = err.Error()
			continue
		}
		status := resp.StatusCode
		resp.Body.Close()
		if status == http.StatusNoContent {
			return nil
		}
		lastResponse = fmt.Sprintf("HTTP %d", status)
	}

	return fmt.Errorf("failed to complete startup after %d attempts, last response: %s\n%s",
		maxAttempts, lastResponse, j.ServerLog())
}

func (j *Jellyfin) ServerLog() string {
	dir := path.Join(j.dataDir, "data", "log")
	newest, err := newestFile(dir)
	if err != nil {
		return fmt.Sprintf("no server log in %s: %v", dir, err)
	}

	content, err := os.ReadFile(newest)
	if err != nil {
		return fmt.Sprintf("cannot read %s: %v", newest, err)
	}

	return fmt.Sprintf("last %d lines of %s:\n%s", logTailLines, newest, tail(string(content), logTailLines))
}

func newestFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	newest := ""
	newestTime := time.Time{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if newest == "" || info.ModTime().After(newestTime) {
			newest = path.Join(dir, entry.Name())
			newestTime = info.ModTime()
		}
	}

	if newest == "" {
		return "", fmt.Errorf("no files found")
	}
	return newest, nil
}

func tail(content string, lines int) string {
	trimmed := strings.TrimRight(content, "\n")
	if trimmed == "" {
		return ""
	}
	split := strings.Split(trimmed, "\n")
	if len(split) > lines {
		split = split[len(split)-lines:]
	}
	return strings.Join(split, "\n")
}

func (j *Jellyfin) UpdateAuthPlugin() error {
	srcDir := path.Join(j.appDir, "app", "plugins", "LDAP-Auth")
	dstDir := path.Join(j.dataDir, "data", "plugins", "LDAP-Auth")

	err := os.RemoveAll(dstDir)
	if err != nil {
		return err
	}

	err = cp.Copy(srcDir, dstDir)
	if err != nil {
		return err
	}

	err = cp.Copy(
		path.Join(j.appDir, "config", "jellyfin", "plugins"),
		path.Join(j.dataDir, "data", "plugins"),
	)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jellyfin) LocalIPv4() string {
	output, err := j.executor.Run("/snap/platform/current/bin/cli", "ipv4")
	if err != nil {
		j.logger.Error("failed to get local ipv4", zap.Error(err))
		return "localhost"
	}
	return strings.TrimSpace(output)
}

func (j *Jellyfin) IPv6() string {
	output, err := j.executor.Run("/snap/platform/current/bin/cli", "ipv6")
	if err != nil {
		j.logger.Error("failed to get ipv6", zap.Error(err))
		return "localhost"
	}
	return strings.TrimSpace(output)
}
