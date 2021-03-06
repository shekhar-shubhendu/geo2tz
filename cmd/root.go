package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/noandrea/geo2tz/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug bool
var settings server.ConfigSchema

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "geo2tz",
	Short: "A rest API service to get the timezone from geo coordinates",
	Long: `Throwing around coordinates to online services seems like not 
  a great idea privacy wise.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	rootCmd.Version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/geo2tz/config.yaml)")
	// for debug logging
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if debug {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/geo2tz")
		viper.SetConfigName("config")
	}
	server.Defaults()
	viper.SetEnvPrefix("GEO2TZ")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in, else use the defaults
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&settings)
		server.Validate(&settings)
		log.Println("Using config file at ", viper.ConfigFileUsed())
	} else {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			viper.Unmarshal(&settings)
		}
	}
	// make the version available via settings
	settings.RuntimeVersion = rootCmd.Version
	log.Debug(fmt.Sprintf("config %#v", settings))
}
