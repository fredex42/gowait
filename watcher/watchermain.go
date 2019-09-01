package watcher

import "github.com/fredex42/gowait/config"
import "flag"
import "fmt"
import "time"
import "github.com/go-redis/redis"

func main() {
	configFilenamePtr := flag.String("configfile", "/etc/gowait/config.yaml", "Path to configuration YAML file")
	pathToWatchPtr := flag.String("path", "", "Specific path to watch. Must be listed in the config file.")

	flag.Parse()

	if configFilenamePtr == nil {
		panic("You need to specify a configuration file")
	}

	configData, configLoadErr := config.LoadConfig(*configFilenamePtr)
	if configLoadErr != nil {
		fmt.Printf("Could not parse configuration: %s\n", configLoadErr)
		panic("Need a readable configuration")
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
	//fmt.Print(configData)
	targetWatcher, findWatcherErr := config.WatcherFor(*pathToWatchPtr, configData)

	if targetWatcher == nil && findWatcherErr == nil {
		fmt.Printf("Could not find watcher %s in the config\n", *pathToWatchPtr)
		panic("Need a defined watcher")
	}
	fmt.Printf("Using watcher for %s\n", targetWatcher.PATH)

	SetupTicker(targetWatcher, targetWatcher.TIMEOUT, redisdb)

	for {
		time.Sleep(time.Duration(3600) * time.Hour)
	}
}
