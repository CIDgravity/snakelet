# snakelet

Ultra minimal go config lib based on [viper](https://github.com/spf13/viper).

Goals:
- simplicity => correctly validate once and return a nice statically typed struct that contains all the config
- testability => dependency to config library only in main package. Everywhere else use pure functions and params

How?:
- define whichever complex type you want for your config and setup which params are required and which are not
- call `InitAndLoad` only once in main package and unmarshal to a struct that is statically typed against the config type, reusing nice viper default for precedence (env -> file -> default)
- crash fast if required fields hasn't been set or if config file cannot be parsed to the config types (ex: random string to int)
- if no error, forget about `viper` and `snakelet`, only use this unmarshaled struct everywhere

Instead of using the lib directly, the source code itself can be used for inspiration as to how to use `viper` internally to achieve the goals set above.

## Usage

TL;DR:
```go
func main() {
	// Get config struct and default variables
	conf := GetDefaultConfig()

	// this will mutate the `conf` variable:
	err := snakelet.InitAndLoad(conf)

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

See the [example directory](./example) for a full working example.

Also checkout [the test file](./config_test.go) and [the source code](./config.go).

## License

All this work is licensed under permissive open-source licensek.

Except written otherwise in particular files, this work is dual-licensed under both [Apache v2](./LICENSE-APACHE) and [MIT](./LICENSE-MIT).
