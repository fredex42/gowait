package main

import (
	"fmt"
	"github.com/fredex42/gowait/config"
	"github.com/fredex42/gowait/watcher"
	_ "github.com/fredex42/gowait/watcher"
	"github.com/go-redis/redis"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You must specify a configuration file as the first positional argument")
	}
	configData, cfgErr := config.LoadConfig(os.Args[1])
	if cfgErr != nil {
		log.Fatal("Could not load configuration:", cfgErr)
	}

	//connect to Redis
	redisdb := redis.NewClient(&redis.Options{
		Addr:     configData.REDIS.REDISHOST,
		Password: configData.REDIS.REDISPASS,
		DB:       configData.REDIS.REDISDB,
	})

	//test redis connection
	_, err := redisdb.Ping().Result()
	if err != nil {
		fmt.Printf("ERROR: Couldn't connect to Redis: %s\n", err)
		panic("Couldn't connect to redis")
	}

	fmt.Printf("INFO: Connected to Redis at %s on db %d\n\n", configData.REDIS.REDISHOST, configData.REDIS.REDISDB)

	log.Println(configData.WATCHERS)

	quitChans := make([](chan struct{}), len(configData.WATCHERS))

	for w := range configData.WATCHERS {
		ch, tickerErr := watcher.SetupTicker(&configData.WATCHERS[w], 2, redisdb)
		if tickerErr != nil {
			log.Print("Could not set up watcher for ", configData.WATCHERS[w], err)
			quitChans[w] = nil
		} else {
			quitChans[w] = ch
		}
		log.Print("Set up watcher for ", configData.WATCHERS[w])
	}

	select {} //now block forever, waiting for the watchers to quit
}
