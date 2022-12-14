package snakelet

import (
	"reflect"
	"strings"
	"testing"

	"github.com/CIDgravity/snakelet/testutil"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
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
	LogLevel      string `mapstructure:"level"` // error | warn | info - case insensitive
	IsJson        bool   `mapstructure:"isJSON"`
	SlowThreshold string `mapstructure:"databaseSlowThreshold"` // Value under which a query is considered slow. "1ms", "1s", etc - anything that's parsable by time.ParseDuration(interval).
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
			LogLevel:      "debug",
			IsJson:        false,
			SlowThreshold: "1ms",
		},
	}

}

// Test that config init and load correctly with no config file
func TestInitAndLoadNoConfig(t *testing.T) {
	conf := GetDefaultConfig()
	err := InitAndLoadWithParams(conf, "", "")
	if err == nil {
		t.Fatalf(`Conf should return an error mentioning that some values are required!`)
	}
	if !strings.Contains(err.Error(), "required") {
		t.Fatalf(`Conf should return an error mentioning that some values are required!: %e`, err)
	}
}

func TestInitAndLoadCorrectConfig(t *testing.T) {
	var yamlExample = []byte(`database:
  user: test_user
  password: "test_pwd"
  host: test_host
  name: test_name
  port: "90"
server:
  port: "90"
`)

	fs := afero.NewMemMapFs()

	fileDirectory := "/opt/viper"
	filePath := fileDirectory + "/config.yml"

	err := fs.Mkdir(testutil.AbsFilePath(t, fileDirectory), 0o777)
	require.NoError(t, err)

	_, err = fs.Create(testutil.AbsFilePath(t, filePath))
	require.NoError(t, err)

	err = afero.WriteFile(fs, testutil.AbsFilePath(t, filePath), yamlExample, 0644)
	require.NoError(t, err)

	// for testing purpose
	viper.SetFs(fs)

	conf := GetDefaultConfig()
	err = InitAndLoadWithParams(conf, filePath, "config-test")
	if err != nil {
		t.Fatalf(`Conf should not error: %e`, err)
	}

	// merge between default and config file's values
	expectedModifiedConfStruct := &Config{
		Database: DatabaseConfig{
			User:     "test_user",
			Password: "test_pwd",
			Host:     "test_host",
			Name:     "test_name",
			Port:     90,
			SslMode:  "disable",
			MaxConns: 200,
		},
		Logs: LogsConfig{
			LogLevel:      "debug",
			IsJson:        false,
			SlowThreshold: "1ms",
		},
		Server: ServerConfig{Port: 90},
	}

	if !reflect.DeepEqual(conf, expectedModifiedConfStruct) {
		t.Fatalf(`Conf didn't unmarshal content from the file properly: %v & %v`, *conf, *expectedModifiedConfStruct)
	}

}

func TestInitAndLoadWrongConfig(t *testing.T) {
	var yamlExample = []byte(`database:
  user: kdfj
  password: djfkdls
  host: dlkjf
  name: value
  port: should_be_a_number
server:
  port: "90"
`)

	fs := afero.NewMemMapFs()

	fileDirectory := "/opt/viper"
	filePath := fileDirectory + "/config.yml"

	err := fs.Mkdir(testutil.AbsFilePath(t, fileDirectory), 0o777)
	require.NoError(t, err)

	_, err = fs.Create(testutil.AbsFilePath(t, filePath))
	require.NoError(t, err)

	err = afero.WriteFile(fs, testutil.AbsFilePath(t, filePath), yamlExample, 0644)
	require.NoError(t, err)

	// for testing purpose
	viper.SetFs(fs)

	conf := GetDefaultConfig()
	err = InitAndLoadWithParams(conf, filePath, "config-test")
	if err == nil {
		t.Fatalf(`Conf should be in error: %e`, err)
	}
	if !strings.Contains(err.Error(), "cannot parse 'database.port' as int") {
		t.Fatalf(`Config should show specific error targeting database.port type: %e`, err)
	}

}
