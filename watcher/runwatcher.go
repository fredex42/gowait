package main

import "time"
import "fmt"
import "github.com/gowait/config"
import "github.com/go-redis/redis"

func dump_record(rec *WatchRecord) () {
  fmt.Printf("Got record:\n")
  if(rec==nil){ return; }
  fmt.Printf("\tPath: %s\n", rec.Path)
  fmt.Printf("\tFilename: %s\n", rec.Filename)
  fmt.Printf("\tmtime: %s\n", rec.LastMtime)
  fmt.Printf("\tstable iterations: %d\n", rec.StableIterations)
}

func dump_records(records []*WatchRecord) (){
  for _, rec := range records {
    dump_record(rec)
  }
}

func timer_func(w *config.Watcher, ticker *time.Ticker, quit chan struct{},redisClient *redis.Client) (){
  for {
    select {
    case <- ticker.C:
      fmt.Print("tick\n")
      records, err := scan_pass(w.PATH,redisClient)
      if(err!=nil){
        fmt.Printf("ERROR: could not perform scan pass: %s\n", err)
      } else {
        //check_and_apply(records)
        dump_records(records)
      }
    case <- quit:
      ticker.Stop()
      return
    }
  }
}

func setup_ticker(w *config.Watcher,duration int, redisClient *redis.Client) (chan struct{}, error){
  ticker := time.NewTicker(time.Duration(duration) * time.Second)
  quit := make(chan struct{})

  go timer_func(w, ticker, quit, redisClient)
  return quit, nil
}
