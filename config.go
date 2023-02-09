package snakelet

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func InitAndLoad(configStruct interface{}, cfgFile string) (*viper.Viper, error) {
	return InitAndLoadWithParams(configStruct, cfgFile, validator.New())
}

// InitAndLoad
// if `cfgFile` == "", it will use only the default values otherwise, it will load the specified config file
// will return an error if occurred, and an instance of viper configuration
func InitAndLoadWithParams(configStruct interface{}, cfgFile string, validate *validator.Validate) (*viper.Viper, error) {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// read from config file
	if err := viper.ReadInConfig(); err != nil {
		if _, isConfigFileNotFoundError := err.(viper.ConfigFileNotFoundError); !isConfigFileNotFoundError {
			return viper.GetViper(), err
		}
	}

	// set default
	var decodedConfig map[string]interface{}
	err := mapstructure.Decode(configStruct, &decodedConfig)

	if err != nil {
		return viper.GetViper(), err
	}

	for k, v := range decodedConfig {
		viper.SetDefault(k, v)
	}

	// unmarshal config file / env variables to exact struct, overriding potential defaults
	// fail if the config does not match the struct
	err = viper.UnmarshalExact(&configStruct)

	if err != nil {
		return viper.GetViper(), err
	}

	if err := validate.Struct(configStruct); err != nil {
		return viper.GetViper(), fmt.Errorf("Missing required config attributes %w\n", err)
	}

	return viper.GetViper(), err
}
