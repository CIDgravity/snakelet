package snakelet

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
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

func InitAndLoad(configStruct interface{}) (error, bool, string) {
	return InitAndLoadWithParams(configStruct, "", "")
}

func InitAndLoadWithParams(configStruct interface{}, cfgFile string, prefix string) (error, bool, string) {
	return InitAndLoadWithParamsAndLogger(configStruct, cfgFile, prefix)
}

// InitAndLoadWithParamsAndLogger The `prefix` must be unique between each project. Env variables are only set if a prefix has been set.
// if `cfgFile` == "", it search for `config.yaml` in the directory where the executable is located,
// else `snakelet` will use cfgFile (full path required)
// will return an error if occured, a boolean if default config is used and the name of the config file used
func InitAndLoadWithParamsAndLogger(configStruct interface{}, cfgFile string, prefix string) (error, bool, string) {

	// preprare config file and env variables
	// only use env variables if prefix is set
	if prefix != "" {
		viper.SetEnvPrefix(prefix)
		replacer := strings.NewReplacer("-", "_")
		viper.SetEnvKeyReplacer(replacer) // replace `-` with `_` in environment variables
		viper.AutomaticEnv()
	}
	if cfgFile == "" {

		// See https://stackoverflow.com/a/18537419/11046178
		executable, err := os.Executable()
		if err != nil {
			return err, false, viper.ConfigFileUsed()
		}

		executablePath := filepath.Dir(executable)
		if err != nil {
			return err, false, viper.ConfigFileUsed()
		}

		viper.AddConfigPath(executablePath)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

	} else {
		viper.SetConfigFile(cfgFile)
	}

	// read from config file
	// if config file not found = will use default value = set specific variable
	useDefaultValue := false

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			useDefaultValue = true
		} else {
			return err, false, viper.ConfigFileUsed()
		}
	}

	// set default
	var decodedConfig map[string]interface{}
	err := mapstructure.Decode(configStruct, &decodedConfig)

	if err != nil {
		return err, useDefaultValue, viper.ConfigFileUsed()
	}

	for k, v := range decodedConfig {
		viper.SetDefault(k, v)
	}

	// unmarshal config file / env variables to exact struct, overriding potential defaults
	// fail if the config does not match the struct
	err = viper.UnmarshalExact(&configStruct)

	if err != nil {
		return err, useDefaultValue, viper.ConfigFileUsed()
	}

	validate := validator.New()
	if err := validate.Struct(configStruct); err != nil {
		return err, useDefaultValue, viper.ConfigFileUsed()
	}

	return nil, useDefaultValue, viper.ConfigFileUsed()
}

// Debug Only used for debug purpose during development (print to standard output)
// DO NOT USE IT IN PRODUCTION OR CRITICAL INFORMATION FROM CONF MAY BE STORED IN LOG SERVERS.
// THIS WILL PRINT THE ENTIRE CONFIGURATION!
func Debug() {
	viper.GetViper().Debug()
}
