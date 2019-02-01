package constrictor

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//type App struct {
//cobra.Command
//}

var (
	app         = &cobra.Command{}
	programName string
)

type runFunc func()

func App(name string, shortDesc string, longDesc string, run runFunc) *cobra.Command {
	app.Use = name
	app.Short = shortDesc
	app.Long = longDesc
	app.Run = func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Printf("how??\n")
		run()
	}
	programName = name

	cobra.OnInitialize(readConfig)

	return app
}

func Launch() {
	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readConfig() {
	//cfg := new(config)

	viper.SetConfigName(programName)
	//viper.AddConfigPath(fmt.Sprintf("/etc/"))
	viper.AddConfigPath(fmt.Sprintf("."))

	viper.SetEnvPrefix(programName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Printf("error reading config file: %s", err)
	}
}

func StringVar(name string, shortName string, defaultVal string, desc string) func() string {
	app.PersistentFlags().StringP(name, shortName, defaultVal, desc)
	viper.BindPFlag(name, app.PersistentFlags().Lookup(name))

	return func() string {
		return viper.GetString(name)
	}
}

func AddressPortVar(name string, shortName string, defaultVal string, desc string) func() string {
	app.PersistentFlags().StringP(name, shortName, defaultVal, desc)
	viper.BindPFlag(name, app.PersistentFlags().Lookup(name))

	return func() string {
		val := viper.GetString(name)
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return fmt.Sprintf(":%d", i)
		}
		return val
	}
}

func TimeDurationVar(name string, shortName string, defaultVal string, desc string) func() time.Duration {
	app.PersistentFlags().StringP(name, shortName, defaultVal, desc)
	viper.BindPFlag(name, app.PersistentFlags().Lookup(name))

	return func() time.Duration {
		if delay, ok := viper.Get(name).(int); ok {
			return time.Duration(time.Duration(delay) * time.Second)
		}
		return viper.GetDuration(name)
	}
}
