package src

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	Utils "github.com/jopemachine/mac-sync-config/src/utils"
	"github.com/keybase/go-keychain"
)

const (
	PREFERENCE_PATH = "~/Library/Preferences/Mac-sync-config"
	CACHE_PATH      = "~/Library/Caches/Mac-sync-config"
)

const (
	MAC_SYNC_CONFIGS_FILE = "mac-sync-configs.yaml"
)

var (
	ConfigFileLastChangedCachePath = strings.Join([]string{CACHE_PATH, "last-changed.json"}, "/")
	LocalPreferencePath            = strings.Join([]string{PREFERENCE_PATH, "local-preference.json"}, "/")
)

var (
	KeychainPreference = GetKeychainPreference()
)

var (
	Flag_OverWrite = false
)

func GetRemoteConfigFolderName() string {
	folderName := os.Getenv("Mac-sync-config-folder")
	if folderName != "" {
		return folderName
	}

	return ".mac-sync-configs"
}

type MacSyncConfigs struct {
	ConfigPathsToSync []string `yaml:"sync"`
}

type KeychainPreferenceType struct {
	GithubId                       string `json:"github_id"`
	GithubToken                    string `json:"github_token"`
	MacSyncConfigGitRepositoryName string `json:"mac_sync_config_git_repository_name"`
}

func scanKeyChainPreference(config *KeychainPreferenceType) {
	Logger.Info("Please enter some information for accessing your Github repository.")
	Logger.Info("This information will be stored in your keychain.")
	Logger.NewLine()

	Logger.Question("Enter your Github ID:")
	ghId := bufio.NewScanner(os.Stdin)
	ghId.Scan()
	config.GithubId = ghId.Text()

	Logger.Question("Enter your Github Access Token:")
	ghToken := bufio.NewScanner(os.Stdin)
	ghToken.Scan()
	config.GithubToken = ghToken.Text()

	Logger.Question("Enter Git repository name for saving the mac-sync-config's configuration files:")
	repoName := bufio.NewScanner(os.Stdin)
	repoName.Scan()
	config.MacSyncConfigGitRepositoryName = repoName.Text()
}

func GetKeychainPreference() KeychainPreferenceType {
	var config KeychainPreferenceType

	dat, err := keychain.GetGenericPassword("Mac-sync-config", "jopemachine", "Mac-sync-config", "org.jopemachine")

	// TODO: Improve below error handling logic.
	// If not exist, create new preference config file
	if len(dat) == 0 {
		scanKeyChainPreference(&config)
		bytesToWrite, err := json.Marshal(config)
		Utils.PanicIfErr(err)

		keyChainItem := keychain.NewItem()
		keyChainItem.SetSecClass(keychain.SecClassGenericPassword)
		keyChainItem.SetService("Mac-sync-config")
		keyChainItem.SetAccount("jopemachine")
		keyChainItem.SetLabel("Mac-sync-config")
		keyChainItem.SetAccessGroup("org.jopemachine")
		keyChainItem.SetData([]byte(bytesToWrite))
		keyChainItem.SetSynchronizable(keychain.SynchronizableNo)
		keyChainItem.SetAccessible(keychain.AccessibleWhenUnlocked)

		err = keychain.AddItem(keyChainItem)

		if err == keychain.ErrorDuplicateItem {
			err = keychain.DeleteItem(keyChainItem)
			Utils.PanicIfErr(err)
			err = keychain.AddItem(keyChainItem)
		}

		Utils.PanicIfErr(err)

		Logger.Success(fmt.Sprintf("mac-sync-config's configuration is saved successfully on the keychain.\n"))
	} else if err != nil {
		panic(err)
	} else {
		err = json.Unmarshal(dat, &config)
		Utils.PanicIfErrWithMsg("Json data malformed. Please delete and remake the keychain", err)
	}

	return config
}

func ReadLocalPreference() map[string]string {
	return ReadJSON(LocalPreferencePath)
}

func ReadLastChanged() map[string]string {
	return ReadJSON(ConfigFileLastChangedCachePath)
}

func WriteLocalPreference(localPreference map[string]string) {
	cacheDir := RelativePathToAbs(PREFERENCE_PATH)
	localPreferencePath := RelativePathToAbs(LocalPreferencePath)

	if _, err := os.Stat(cacheDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		Utils.PanicIfErr(err)
	} else if _, err := os.Stat(localPreferencePath); !errors.Is(err, os.ErrNotExist) {
		os.Remove(localPreferencePath)
	}

	WriteJSON(localPreferencePath, localPreference)
}

func WriteLastChangedConfigFile(lastChangedConfig map[string]string) {
	cacheDir := RelativePathToAbs(CACHE_PATH)
	configFileLastChangedCachePath := RelativePathToAbs(ConfigFileLastChangedCachePath)

	if _, err := os.Stat(cacheDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		Utils.PanicIfErr(err)
	} else if _, err := os.Stat(configFileLastChangedCachePath); !errors.Is(err, os.ErrNotExist) {
		os.Remove(configFileLastChangedCachePath)
	}

	WriteJSON(configFileLastChangedCachePath, lastChangedConfig)
}

func ClearCache() {
	configFileLastChangedCachePath := RelativePathToAbs(ConfigFileLastChangedCachePath)

	err := os.Remove(configFileLastChangedCachePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	Logger.Success("Cache file cleared.")
}

func PrintMacSyncConfigs() {
	Logger.Log(GetMacSyncConfigs())
}

func SwitchProfile(profileName string) {
	localPreference := ReadLocalPreference()
	localPreferencePath := RelativePathToAbs(LocalPreferencePath)

	localPreference["profile"] = profileName
	WriteJSON(localPreferencePath, localPreference)
}
