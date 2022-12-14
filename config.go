package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Logger interface {
	Printf(format string)
}

type DefaultLogger struct {
	logger Logger
}

func (p *DefaultLogger) Printf(format string) {
	fmt.Printf(format + "\n")
}

func InitAndLoad(configStruct interface{}) error {
	return InitAndLoadWithParams(configStruct, "", "")
}

func InitAndLoadWithParams(configStruct interface{}, cfgFile string, prefix string) error {
	return InitAndLoadWithParamsAndLogger(configStruct, cfgFile, prefix, &DefaultLogger{})
}

// Prefix must be unique in between each projects. Env variables are only set if a prefix has been set
// if cfgFile == "", it will used config.yaml in the directory where the executable is located.
// else it will use cfgFile (full path required)
func InitAndLoadWithParamsAndLogger(configStruct interface{}, cfgFile string, prefix string, logger Logger) error {

	// preprare config file and env variables
	if prefix != "" {
		// only use env variables if prefix if set
		viper.SetEnvPrefix(prefix)
		replacer := strings.NewReplacer("-", "_")
		viper.SetEnvKeyReplacer(replacer) // replace `-` with `_` in environment variables
		viper.AutomaticEnv()
	}
	if cfgFile == "" {
		// See https://stackoverflow.com/a/18537419/11046178
		executable, err := os.Executable()
		if err != nil {
			return err
		}
		executablePath := filepath.Dir(executable)
		if err != nil {
			return err
		}
		viper.AddConfigPath(executablePath)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	} else {
		viper.SetConfigFile(cfgFile)
	}

	// read from config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Printf("Config file not found, continuing with other values")
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("Config file found, but an error occured: %w", err)
		}
	} else {
		logger.Printf("Using config file: " + viper.ConfigFileUsed())
	}

	// set default
	var decodedConfig map[string]interface{}
	err := mapstructure.Decode(configStruct, &decodedConfig)
	if err != nil {
		return err
	}
	for k, v := range decodedConfig {
		viper.SetDefault(k, v)
	}

	// unmarshal config file / env variables to exact struct, overriding potential defaults
	// fail if the config does not match the struct
	err = viper.UnmarshalExact(&configStruct)
	if err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(configStruct); err != nil {
		return fmt.Errorf("Missing required config attributes %w\n", err)
	}

	return nil
}

// Only used for debug purpose during development (print to standard output)
// DO NOT USE IT IN PRODUCTION OR CRITICAL INFORMATION FROM CONF MAY BE STORED IN LOG SERVERS
// WARNING: this will print the entire configuration!
func Debug() {
	viper.GetViper().Debug()
}
