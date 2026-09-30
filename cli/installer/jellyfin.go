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
	serverUnit   = "snap.jellyfin.server.service"
	logTailLines = "40"
	fatalMarker  = "Unhandled Exception"
)

type CommandRunner interface {
	Run(app string, args ...string) (string, error)
}

type Jellyfin struct {
	appDir, dataDir string
	client          *http.Client
	executor        CommandRunner
	logger          *zap.Logger
}

func NewJellyfin(appDir, dataDir string, executor CommandRunner, logger *zap.Logger) *Jellyfin {
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

		output := j.serverLog()
		if strings.Contains(output, fatalMarker) {
			return fmt.Errorf("server failed to start, last web response: %s\n%s", lastResponse, output)
		}
	}

	return fmt.Errorf("web endpoint not available after %d attempts, last web response: %s\n%s",
		maxAttempts, lastResponse, j.serverLog())
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
		maxAttempts, lastResponse, j.serverLog())
}

func (j *Jellyfin) serverLog() string {
	output, err := j.executor.Run("journalctl", "-u", serverUnit, "-n", logTailLines, "--no-pager")
	if err != nil {
		return fmt.Sprintf("cannot read the journal of %s: %v\n%s", serverUnit, err, output)
	}
	return fmt.Sprintf("last %s journal lines of %s:\n%s", logTailLines, serverUnit, output)
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
