package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/config"
	"server/handlers"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var port int
	var configFile string
	flag.IntVar(&port, "port", 80, "the port to listen on")
	flag.StringVar(&configFile, "config_file", "config.yaml", "the configuration file to be used")

	if _, err := os.Stat(configFile); err != nil {
		log.Fatalln("Config file does not exist:", configFile)
	}

	fmt.Println("Configuration file:", configFile)
	fmt.Println("Server port:", port)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	flag.Parse()

	config.ReadConfig()
	if err := config.WriteConfig(); err != nil {
		log.Fatalln("Cannot write config", err)
	}

	http.Handle("/api/v2/ping", http.HandlerFunc(handlers.Ping))

	log.Printf("Server listening on %d\n", port)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
