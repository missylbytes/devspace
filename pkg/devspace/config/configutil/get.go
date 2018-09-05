package configutil

import (
	"os"
	"unsafe"

	"github.com/covexo/devspace/pkg/util/log"

	"github.com/covexo/devspace/pkg/devspace/config/v1"
)

//ConfigInterface defines the pattern of every config
type ConfigInterface interface{}

const configGitignore = `logs/
overwrite.yaml
`

const configPath = "/.devspace/config.yaml"
const overwriteConfigPath = "/.devspace/overwrite.yaml"

var config = makeConfig()
var configRaw = makeConfig()
var overwriteConfig = makeConfig()
var overwriteConfigRaw = makeConfig()
var workdir string

func init() {
	workdir, _ = os.Getwd()
}

//ConfigExists chacks whether the yaml file for the config exists
func ConfigExists() (bool, error) {
	_, configNotFound := os.Stat(workdir + configPath)

	return (configNotFound == nil), nil
}

func GetConfig(reload bool) *v1.Config {
	isLoaded := (config.Version != nil)

	if !isLoaded || reload {
		err := loadConfig(configRaw, configPath)

		if err != nil {
			log.Fatal("Unable to load config.")
		}
		GetOverwriteConfig()

		merge(config, configRaw, unsafe.Pointer(&config), unsafe.Pointer(configRaw))
		merge(config, overwriteConfig, unsafe.Pointer(&config), unsafe.Pointer(overwriteConfig))
	}

	if config.Cluster.UseKubeConfig != nil && *config.Cluster.UseKubeConfig {
		loadClusterConfig(config.Cluster, false)
	}
	return config
}

func GetOverwriteConfig() *v1.Config {
	isLoaded := (overwriteConfig.Version != nil)

	if !isLoaded {
		//ignore error as overwrite.yaml is optional
		loadConfig(overwriteConfigRaw, overwriteConfigPath)

		merge(overwriteConfig, overwriteConfigRaw, unsafe.Pointer(&overwriteConfig), unsafe.Pointer(overwriteConfigRaw))
	}
	return overwriteConfig
}

func GetConfigInstance() *v1.Config {
	return config
}
