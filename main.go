package main

import (
	"encoding/json"
	"flag"
	"os"

	"monopool/config"
	"monopool/pool"
	"monopool/utils"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("main")

const defaultConfigFileName = "config.json"

var (
	configFileName = flag.String("c", defaultConfigFileName, "configuration file for pool")
	logLevel       = flag.String("l", "debug", "log level")
)

func main() {
	flag.Parse()

	lvl, err := logging.LevelFromString(*logLevel)
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)

	var conf config.Options
	if !utils.FileExists(*configFileName) {
		log.Panic("the config file " + *configFileName + " does not exist")
	}

	f, err := os.Open(*configFileName)
	if err != nil {
		log.Panic(err)
	}

	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		log.Panic(err)
	}

	p := pool.NewPool(&conf)
	p.Init()
	for {
		select {}
	}
}
