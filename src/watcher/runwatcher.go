package main

import "time"
import "fmt"
import "config"

func timer_func(w *config.Watcher, ticker *time.Ticker, quit chan struct{}) (){
  for {
    select {
    case <- ticker.C:
      fmt.Print("tick\n")
    case <- quit:
      ticker.Stop()
      return
    }
  }
}

func setup_ticker(w *config.Watcher,duration int) (chan struct{}, error){
  ticker := time.NewTicker(time.Duration(duration) * time.Second)
  quit := make(chan struct{})

  go timer_func(w, ticker, quit)
  return quit, nil
}
