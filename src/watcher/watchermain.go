package main

import . "config"
import "flag"
import "fmt"
import "time"

func main() {
  configFilenamePtr := flag.String("configfile","/etc/gowait/config.yaml","Path to configuration YAML file")
  pathToWatchPtr := flag.String("path","","Specific path to watch. Must be listed in the config file.")

  flag.Parse()

  if(configFilenamePtr==nil){
    panic("You need to specify a configuration file")
  }

  configData, configLoadErr := LoadConfig(*configFilenamePtr)
  if(configLoadErr!=nil){
    fmt.Printf("Could not parse configuration: %s\n", configLoadErr)
    panic("Need a readable configuration")
  }

  //fmt.Print(configData)
  targetWatcher, findWatcherErr := WatcherFor(*pathToWatchPtr, configData)

  if(targetWatcher==nil && findWatcherErr==nil){
    fmt.Printf("Could not find watcher %s in the config\n", *pathToWatchPtr)
    panic("Need a defined watcher")
  }
  fmt.Printf("Using watcher for %s\n", targetWatcher.PATH)

  setup_ticker(targetWatcher, targetWatcher.TIMEOUT)

  for {
    time.Sleep(time.Duration(3600)*time.Hour)
  }
}
