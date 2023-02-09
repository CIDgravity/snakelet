# snakelet

Ultra minimal go config lib based on [viper](https://github.com/spf13/viper) and [validator](https://github.com/go-playground/validator).

Goals:
- simplicity => correctly validate once and return a nice statically typed struct that contains all the config
- testability => dependency to config library only in main package. Everywhere else, just pass a subset of the config struct as parameter

How?:
- define whichever complex type you want for your config and setup which params are required and which are not
- call `InitAndLoad` only once in main package and unmarshal to a struct that is statically typed against the config type, reusing nice viper default for precedence (env -> file -> default)
- crash fast if required fields hasn't been set or if config file cannot be parsed to the config types (ex: random string to int)
- if no error, forget about `viper` and `snakelet`, only use this unmarshaled struct everywhere

Instead of using the lib directly, the source code itself can be used for inspiration as to how to use `viper` internally to achieve the goals set above.

## Install

```
go get github.com/CIDgravity/snakelet
```

## Usage

```go
package main

import (
	"fmt"

	snakelet "github.com/CIDgravity/snakelet"
)

type DatabaseConfig struct {
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required"`
	Name     string `mapstructure:"name" validate:"required"`
	SslMode  string `mapstructure:"sslMode"`
	MaxConns int    `mapstructure:"maxConns"`
}

type LogsConfig struct {
	LogLevel                string `mapstructure:"level"` // error | warn | info - case insensitive
	BackendLogsToJsonFormat bool   `mapstructure:"isJSON"`
	DatabaseSlowThreshold   string `mapstructure:"databaseSlowThreshold"` // Value under which a query is considered slow. "1ms", "1s", etc - anything that's parsable by time.ParseDuration(interval).
}

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Logs     LogsConfig     `mapstructure:"log"`
	Server   ServerConfig   `mapstructure:"server"`
}

type ServerConfig struct {
	Port int `mapstructure:"port" validate:"required"`
}

// only set default for non-required tag
func GetDefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			SslMode:  "disable",
			MaxConns: 200,
		},
		Logs: LogsConfig{
			LogLevel:                "debug",
			BackendLogsToJsonFormat: false,
			DatabaseSlowThreshold:   "1ms",
		},
	}

}

func main() {
	// Get config struct and default variables
	conf := GetDefaultConfig()

	// this will mutate the `conf` variable:
	_, err := snakelet.InitAndLoad(conf)

	// Example with params usage:
	// err := snakelet.InitAndLoadWithParams(conf, "/opt/company/project.yaml", "company-project")

	// Will throw an error because some param are required
	if err != nil {
		fmt.Println("Unable to init and load config: %w", err)
	}
	// Pass `conf` or a subset of this variable as param throughout the rest of the code so you don't have any dependency on snakelet or viper whatsoever.
	// This also contributes to improving testability of your code (you can code with pure funcitons):
	// someServerComponent := server.NewServer(conf.Server)
	// someDbComponent := database.NewDatabase(conf.Database)
}
```

See the [example directory](./example) for the full working example that you can run.

Also checkout [the test file](./config_test.go) and [the source code](./config.go).

## License

Copyright (c), 2022, CIDgravity.

All this work is licensed under permissive open-source license.

Except written otherwise in particular files, this work is dual-licensed under both [Apache v2](./LICENSE-APACHE) and [MIT](./LICENSE-MIT).
