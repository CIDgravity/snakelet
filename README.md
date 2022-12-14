# snakelet

Ultra minimal go config lib based on [viper](https://github.com/spf13/viper).

Alternatively, the source code itself can be used for inspiration as to how to use `viper` internally.

## Usage

```go
import (
	snakelet "github.com/CIDgravity/snakelet"
	"fmt"
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


// Get config struct and default variables
conf := GetDefaultConfig()
// this will mutate the `conf` variable:
err := snakelet.InitAndLoadWithParams(conf, "/opt/company/project.yaml", "company-project")
// or use default config, see source code for details:
// err := snakelet.InitAndLoad(conf)
if err != nil {
	return fmt.Errorf("Unable to init and load config: %w", err)
}
// Pass `conf` or a subset of this variable as param throughout the rest of the code so you don't have any dependency on snakelet or viper whatsoever.
// This also contributes to improving testability of your code (you can code with pure funcitons):
// someServerComponent := server.NewServer(conf.Server)
// someDbComponent := database.NewDatabase(conf.Database)
```

Also checkout `./config_test.go` and [InitAndLoadWithLogger](./config.go).

## License

All this work is licensed under permissive open-source license.

Except written otherwise in particular files, this work is dual-licensed under both [Apache v2](./LICENSE-APACHE) and [MIT](./LICENSE-MIT).
