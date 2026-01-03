package installer

import (
	"fmt"
	"os"
	"path"

	cp "github.com/otiai10/copy"
	"github.com/syncloud/golib/config"
	"github.com/syncloud/golib/linux"
	"github.com/syncloud/golib/platform"
	"go.uber.org/zap"
)

const App = "jellyfin"

type Variables struct {
	App              string
	AppDir           string
	DataDir          string
	CommonDir        string
	AppKey           string
	AppUrl           string
	Domain           string
	AuthUrl          string
	AuthClientId     string
	AuthClientSecret string
	AuthRedirectUri  string
	AppDomain        string
	IPv6             string
	LocalIPv4        string
}

type Installer struct {
	newVersionFile     string
	currentVersionFile string
	configDir          string
	platformClient     *platform.Client
	installFile        string
	appDir             string
	dataDir            string
	commonDir          string
	artisanPath        string
	jellyfin           *Jellyfin
	logger             *zap.Logger
}

func New(logger *zap.Logger) *Installer {
	appDir := fmt.Sprintf("/snap/%s/current", App)
	dataDir := fmt.Sprintf("/var/snap/%s/current", App)
	commonDir := fmt.Sprintf("/var/snap/%s/common", App)
	configDir := path.Join(dataDir, "config")
	executor := NewExecutor(logger)
	return &Installer{
		newVersionFile:     path.Join(appDir, "version"),
		currentVersionFile: path.Join(dataDir, "version"),
		configDir:          configDir,
		platformClient:     platform.New(),
		installFile:        path.Join(dataDir, "installed"),
		appDir:             appDir,
		dataDir:            dataDir,
		commonDir:          commonDir,
		jellyfin:           NewJellyfin(appDir, dataDir, executor, logger),
		logger:             logger,
	}
}

func (i *Installer) Install() error {
	err := CreateUser(App)
	if err != nil {
		return err
	}

	err = i.UpdateConfigs()
	if err != nil {
		return err
	}

	err = i.FixPermissions()
	if err != nil {
		return err
	}

	err = i.StorageChange()
	if err != nil {
		return err
	}
	return nil
}

func (i *Installer) Configure() error {

	if i.IsInstalled() {
		err := i.Upgrade()
		if err != nil {
			return err
		}
	} else {
		err := i.Initialize()
		if err != nil {
			return err
		}
	}

	return i.UpdateVersion()
}

func (i *Installer) DomainChange() error {
	err := i.UpdateConfigs()
	if err != nil {
		return err
	}
	err = i.FixPermissions()
	if err != nil {
		return err
	}
	return nil
}

func (i *Installer) Initialize() error {
	err := i.StorageChange()
	if err != nil {
		return err
	}

	err = i.jellyfin.Complete()
	if err != nil {
		return err
	}

	err = os.WriteFile(i.installFile, []byte("installed"), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (i *Installer) Upgrade() error {

	err := i.StorageChange()
	if err != nil {
		return err
	}

	return nil
}

func (i *Installer) IsInstalled() bool {
	// migrate from common status, remove after the next release
	old := path.Join(i.commonDir, "installed")
	_, err := os.Stat(old)
	if err == nil {
		i.logger.Info("migrating old installed status")

		err = os.WriteFile(i.installFile, []byte("installed"), 0644)
		if err != nil {
			i.logger.Error("cannot migrate installed status", zap.Error(err))
			return true
		}
		return true
	}
	// migrate end

	_, err = os.Stat(i.installFile)
	return err == nil
}

func (i *Installer) PreRefresh() error {
	return nil
}

func (i *Installer) PostRefresh() error {

	err := i.UpdateConfigs()
	if err != nil {
		return err
	}

	err = i.ClearVersion()
	if err != nil {
		return err
	}

	err = i.FixPermissions()
	if err != nil {
		return err
	}
	return nil
}

func (i *Installer) StorageChange() error {
	storageDir, err := i.platformClient.InitStorage(App, App)
	if err != nil {
		return err
	}

	err = Chown(storageDir, App)
	if err != nil {
		return err
	}
	return nil
}

func (i *Installer) ClearVersion() error {
	return os.RemoveAll(i.currentVersionFile)
}

func (i *Installer) UpdateVersion() error {
	return cp.Copy(i.newVersionFile, i.currentVersionFile)
}

func (i *Installer) UpdateConfigs() error {
	err := linux.CreateMissingDirs(
		path.Join(i.dataDir, "nginx"),
		path.Join(i.dataDir, "data", "plugins"),
		path.Join(i.dataDir, "cache"),
	)
	if err != nil {
		return err
	}

	err = i.jellyfin.LinkAuthPlugin()
	if err != nil {
		return err
	}

	appDomain, err := i.platformClient.GetAppDomainName(App)
	if err != nil {
		return err
	}
	variables := Variables{
		App:       App,
		AppDir:    i.appDir,
		DataDir:   i.dataDir,
		CommonDir: i.commonDir,
		AppDomain: appDomain,
		LocalIPv4: i.jellyfin.LocalIPv4(),
		IPv6:      i.jellyfin.IPv6(),
	}

	err = config.Generate(
		path.Join(i.appDir, "config"),
		path.Join(i.dataDir, "config"),
		variables,
	)
	if err != nil {
		return err
	}

	return nil
}

func (i *Installer) BackupPreStop() error {
	return i.PreRefresh()
}

func (i *Installer) RestorePreStart() error {
	return i.PostRefresh()
}

func (i *Installer) RestorePostStart() error {
	return i.Configure()
}

func (i *Installer) AccessChange() error {
	err := i.DomainChange()
	if err != nil {
		return err
	}
	return i.UpdateConfigs()
}

func (i *Installer) FixPermissions() error {
	err := Chown(i.dataDir, App)
	if err != nil {
		return err
	}
	err = Chown(i.commonDir, App)
	if err != nil {
		return err
	}
	return nil
}
