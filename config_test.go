package main

import "testing"

func TestLoad(t *testing.T) {
  config, err := LoadConfig("sampleconfig.yaml")
  if(err != nil){
    t.Errorf("load_config returned error: %s", err)
  }

  if(len(config.WATCHERS)!=2){
    t.Errorf("Expected 2 watchers, got %d", len(config.WATCHERS))
  }

  watcher := config.WATCHERS[0]

  if(watcher.PATH!="/Users/localhome/test1") {
    t.Errorf("Watcher path did not match expected")
  }

  if(config.REDIS.REDISHOST!="localhost:6379"){
    t.Errorf("Unexpected redis host")
  }
  if(config.REDIS.REDISDB!=0){
    t.Errorf("Unexpected redis db index")
  }
}

func TestWatcherFor(t *testing.T) {
  config, err := LoadConfig("sampleconfig.yaml")
  if(err != nil){
    t.Errorf("load_config returned error: %s", err)
  }

  watcher, err := WatcherFor("/Users/localhome/test2", config)
  if(watcher==nil){
    t.Errorf("Could not locate expected watcher path")
  }
  if(watcher.TIMEOUT!=4 || watcher.STABLE!=2){
    t.Errorf("Unexpected timeout or stable time in located watcher")
  }
}
