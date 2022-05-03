package handler

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
)

var homeDir string

func init() {
	flag.StringVar(&homeDir, "tm-home", "", "Path to the tendermint config directory (if empty, uses $HOME/.tendermint)")
}

func ConfigureChain() {
	flag.Parse()
	if homeDir == "" {
		homeDir = os.ExpandEnv("$HOME/.tendermint")
	}
	config := cfg.DefaultConfig()

	config.SetRoot(homeDir)

	viper.SetConfigFile(fmt.Sprintf("%s/%s", homeDir, "config/config.toml"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Reading config: %v", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Decoding config: %v", err)
	}
	if err := config.ValidateBasic(); err != nil {
		log.Fatalf("Invalid configuration data: %v", err)
	}

}
