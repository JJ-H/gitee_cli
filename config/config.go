package config

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var Conf Config

type Config struct {
	AccessToken              string `yaml:"access_token"`
	ApiPrefix                string `yaml:"api_prefix"`
	UserId                   int    `yaml:"user_id"`
	UserName                 string `yaml:"user_name"`
	DefaultEntPath           string `yaml:"default_ent_path"`
	DefaultPathWithNamespace string `yaml:"default_path_with_namespace"`
	PremiumBuildPrefix       string `yaml:"premium_build_prefix"`
	SaasBuildPrefix          string `yaml:"saas_build_prefix"`
	CookiesJar               string `yaml:"cookies_jar"`
	DefaultEditor            string `yaml:"default_editor"`
}

func Read(key string) (string, error) {
	switch key {
	case "access_token":
		return Conf.AccessToken, nil
	case "api_prefix":
		return Conf.ApiPrefix, nil
	case "user_id":
		return fmt.Sprintf("%v", Conf.UserId), nil
	case "user_name":
		return Conf.UserName, nil
	case "default_ent_path":
		return Conf.DefaultEntPath, nil
	case "default_path_with_namespace":
		return Conf.DefaultPathWithNamespace, nil
	case "premium_build_prefix":
		return Conf.PremiumBuildPrefix, nil
	case "saas_build_prefix":
		return Conf.SaasBuildPrefix, nil
	case "cookies_jar":
		return Conf.CookiesJar, nil
	case "default_editor":
		return Conf.DefaultEditor, nil
	default:
		return "", fmt.Errorf("Unknown config key: %s", key)
	}
}

// Update updates the configuration values from a provided map.
func Update(values map[string]interface{}) error {
	for key, value := range values {
		switch key {
		case "access_token":
			Conf.AccessToken = value.(string)
		case "api_prefix":
			Conf.ApiPrefix = value.(string)
		case "user_id":
			Conf.UserId, _ = strconv.Atoi(parseInput(value))
		case "user_name":
			Conf.UserName = value.(string)
		case "default_ent_path":
			Conf.DefaultEntPath = value.(string)
		case "default_path_with_namespace":
			Conf.DefaultPathWithNamespace = value.(string)
		case "premium_build_prefix":
			Conf.PremiumBuildPrefix = value.(string)
		case "saas_build_prefix":
			Conf.SaasBuildPrefix = value.(string)
		case "cookies_jar":
			Conf.CookiesJar = strings.TrimSpace(value.(string))
		case "default_editor":
			Conf.DefaultEditor = value.(string)
		default:
			return fmt.Errorf("Unknown configuration key: %s", key)
		}
	}

	// Save the updated configuration to the file
	config, err := yaml.Marshal(&Conf)
	if err != nil {
		return fmt.Errorf("Error marshalling configuration: %w", err)
	}

	homeDir, _ := os.UserHomeDir()
	configPath := path.Join(homeDir, ".gitee", "config.yml")

	err = os.WriteFile(configPath, config, 0644)

	if err != nil {
		return fmt.Errorf("Error overwriting configuration file: %w", err)
	}

	return nil
}

func parseInput(input interface{}) string {
	switch input.(type) {
	case string:
		return input.(string)
	case int:
		return fmt.Sprintf("%v", input.(int))
	default:
		return ""
	}
}

func init() {

	homeDir, _ := os.UserHomeDir()
	configPath := path.Join(homeDir, ".gitee", "config.yml")

	config, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("读取配置文件失败！请检查 %s 配置内容！\n", configPath)
		os.Exit(1)
	}

	err = yaml.Unmarshal(config, &Conf)

	if err != nil {
		fmt.Printf("初始化配置文件失败，请检查 %s 配置内容！\n", configPath)
		os.Exit(1)
	}

	// 兼容 bubbletea border 渲染问题
	// https://github.com/charmbracelet/lipgloss/issues/40
	os.Setenv("RUNEWIDTH_EASTASIAN", "0")
}
