package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configType            = "json"
	configEnvPrefix       = "NEW_RELIC_CLI"
	globalScopeIdentifier = "*"
)

type ConfigKey string
type CredentialKey string

const (
	LogLevel       ConfigKey = "loglevel"
	PluginDir      ConfigKey = "plugindir"
	PrereleaseMode ConfigKey = "prereleasefeatures"
	SendUsageData  ConfigKey = "sendusagedata"

	APIKey     CredentialKey = "apiKey"
	Region     CredentialKey = "region"
	AccountID  CredentialKey = "accountID"
	LicenseKey CredentialKey = "licenseKey"
)

var (
	configDir              string
	configFilename         = "config.json"
	credsFilename          = "credentials.json"
	defaultProfileFilename = "default-profile.json"
	defaultProfileValue    string
	viperConfig            *viper.Viper
	viperCreds             *viper.Viper
)

func init() {
	var err error
	configDir, err = getDefaultConfigDirectory()
	if err != nil {
		log.Error("could not get config directory")
	}
}

func GetConfigValue(key ConfigKey) interface{} {
	return viperConfig.Get(keyGlobalScope(string(key)))
}

func GetCredentialValue(key CredentialKey) interface{} {
	return viperCreds.Get(keyDefaultProfile(string(key)))
}

func GetDefaultProfileName() string {
	return defaultProfileValue
}

func SetLogLevel(logLevel string) error {
	return setConfigValue(LogLevel, logLevel)
}

func SetPluginDirectory(directory string) error {
	return setConfigValue(PluginDir, directory)
}

func SetPreleaseFeatures(prereleaseFeatures string) error {
	return setConfigValue(PrereleaseMode, prereleaseFeatures)
}

func SetSendUsageData(sendUsageData string) error {
	return setConfigValue(SendUsageData, sendUsageData)
}

func SetAPIKey(profileName string, apiKey string) error {
	return setCredentialValue(profileName, APIKey, apiKey)
}

func SetRegion(profileName string, region string) error {
	return setCredentialValue(profileName, Region, region)
}

func SetAccountID(profileName string, accountID string) error {
	return setCredentialValue(profileName, AccountID, accountID)
}

func SetLicenseKey(profileName string, licenseKey string) error {
	return setCredentialValue(profileName, LicenseKey, licenseKey)
}

func SetDefaultProfileName(profileName string) error {
	return saveDefaultProfileName(profileName)
}

func setConfigValue(key ConfigKey, value string) error {
	viperConfig.Set(keyGlobalScope(string(key)), value)

	if err := viperConfig.WriteConfigAs(path.Join(configDir, configFilename)); err != nil {
		return err
	}

	return nil
}

func setCredentialValue(profileName string, key CredentialKey, value string) error {
	keyPath := fmt.Sprintf("%s.%s", profileName, key)
	viperCreds.Set(keyPath, value)

	if err := viperCreds.WriteConfigAs(path.Join(configDir, credsFilename)); err != nil {
		return err
	}

	return nil
}

func load() error {
	if err := loadConfigFile(); err != nil {
		return err
	}

	if err := loadCredsFile(); err != nil {
		return err
	}

	if err := loadDefaultProfileFile(); err != nil {
		return err
	}

	return nil
}

func loadConfigFile() error {
	viperConfig = viper.New()
	viperConfig.SetEnvPrefix(configEnvPrefix)
	viperConfig.SetConfigName(configFilename)
	viperConfig.SetConfigType(configType)
	viperConfig.AddConfigPath(configDir)
	viperConfig.AutomaticEnv()

	if err := loadFile(viperConfig); err != nil {
		log.Debugf("config file not found: %s", path.Join(configDir, configFilename))
	}

	return nil
}

func loadCredsFile() error {
	viperCreds = viper.New()
	viperCreds.SetEnvPrefix(configEnvPrefix)
	viperCreds.SetConfigName(credsFilename)
	viperCreds.SetConfigType(configType)
	viperCreds.AddConfigPath(configDir)
	viperCreds.AutomaticEnv()

	if err := loadFile(viperCreds); err != nil {
		log.Debugf("credentials file not found: %s", path.Join(configDir, configFilename))
	}

	return nil
}

func loadDefaultProfileFile() error {
	defaultProfileFilePath := filepath.Join(configDir, defaultProfileFilename)
	defaultProfileBytes, err := ioutil.ReadFile(defaultProfileFilePath)
	if err != nil {
		return nil
	}

	defaultProfileValue = strings.Trim(string(defaultProfileBytes), "\"")

	return nil
}

func loadFile(v *viper.Viper) error {
	err := v.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Debug("file not found, using defaults")
	} else if e, ok := err.(viper.ConfigParseError); ok {
		return e
	}

	return nil
}

func saveDefaultProfileName(profileName string) error {
	defaultProfileFilePath := filepath.Join(configDir, defaultProfileFilename)
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	if err := ioutil.WriteFile(defaultProfileFilePath, []byte(profileName), os.FileMode(flags)); err != nil {
		return nil
	}

	defaultProfileValue = profileName

	return nil
}

func keyGlobalScope(key string) string {
	return fmt.Sprintf("%s.%s", globalScopeIdentifier, key)
}

func keyDefaultProfile(key string) string {
	return fmt.Sprintf("%s.%s", defaultProfile(), key)
}

func defaultProfile() string {
	return "default"
}

func getDefaultConfigDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.newrelic", home), nil
}