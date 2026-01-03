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

	"go.uber.org/zap"
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

	// Wait for REST endpoint to be available
	webURL := "http://unix/web/"
	maxWebAttempts := 10
	for attempt := 0; attempt < maxWebAttempts; attempt++ {
		resp, err := j.client.Get(webURL)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		if attempt == maxWebAttempts-1 {
			return fmt.Errorf("web endpoint not available after %d attempts", maxWebAttempts)
		}
		time.Sleep(1 * time.Second)
	}

	// Complete startup wizard
	completeURL := "http://unix/Startup/Complete"
	maxAttempts := 20
	var lastError string
	for attempt := 0; attempt < maxAttempts; attempt++ {
		resp, err := j.client.Post(completeURL, "application/json", nil)
		if err != nil {
			lastError = fmt.Sprintf("error: %v", err)
		} else {
			lastError = fmt.Sprintf("%d: request failed", resp.StatusCode)
			if resp.StatusCode == 204 {
				resp.Body.Close()
				return nil
			}
			resp.Body.Close()
		}
		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf("failed to complete startup: %s", lastError)
}

func (j *Jellyfin) LinkAuthPlugin() error {
	srcDir := path.Join(j.appDir, "app", "plugins", "LDAP-Auth")
	dstDir := path.Join(j.dataDir, "data", "plugins", "LDAP-Auth")

	_, err := os.Lstat(dstDir)
	if err == nil {
		fileInfo, err := os.Lstat(dstDir)
		if err != nil {
			return err
		}
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			err = os.Remove(dstDir)
			if err != nil {
				return err
			}
		} else {
			err = os.RemoveAll(dstDir)
			if err != nil {
				return err
			}
		}
	}

	err = os.MkdirAll(path.Dir(dstDir), 0755)
	if err != nil {
		return err
	}

	return os.Symlink(srcDir, dstDir)
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
