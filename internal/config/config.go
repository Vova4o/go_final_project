package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var flags = pflag.NewFlagSet("flags", pflag.ExitOnError)

func init() {
	// Define the flags and bind them to viper
	flags.StringP("ServerAddress", "a", ":7540", "HTTP server network address")
	flags.StringP("DBPath", "d", "scheduler.db", "Path to the SQLite database file")
	flags.StringP("Password", "s", "", "Password for the app")

	// Parse the command-line flags
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Printf("Error parsing flags: %v", err)
	}

	// Bind the flags to viper
	bindFlagToViper("ServerAddress")
	bindFlagToViper("DBPath")
	bindFlagToViper("Password")

	// Set the environment variable names
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnvToViper("ServerAddress", "TODO_PORT")
	bindEnvToViper("DBPath", "TODO_DBFILE")
	bindEnvToViper("Password", "TODO_PASSWORD")

	// Read the environment variables
	viper.AutomaticEnv()
}

func bindFlagToViper(flagName string) {
	if err := viper.BindPFlag(flagName, flags.Lookup(flagName)); err != nil {
		log.Println(err)
	}
}

func bindEnvToViper(viperKey, envKey string) {
	if err := viper.BindEnv(viperKey, envKey); err != nil {
		log.Println(err)
	}
}

func Address() string {
	return viper.GetString("ServerAddress")
}

func DBPath() string {
	return viper.GetString("DBPath")
}

func Password() string {
	return viper.GetString("Password")
}
