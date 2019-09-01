package watcher

import "time"
import "fmt"
import "github.com/fredex42/gowait/config"
import "github.com/go-redis/redis"
import "github.com/fredex42/gowait/filescanner"

func dumprecord(rec *filescanner.WatchRecord) {
	fmt.Printf("Got record:\n")
	if rec == nil {
		return
	}
	fmt.Printf("\tPath: %s\n", rec.Path)
	fmt.Printf("\tFilename: %s\n", rec.Filename)
	fmt.Printf("\tmtime: %s\n", rec.LastMtime)
	fmt.Printf("\tstable iterations: %d\n", rec.StableIterations)
}

func DumpRecords(records []*filescanner.WatchRecord) {
	for _, rec := range records {
		dumprecord(rec)
	}
}

func TimerFunc(w *config.Watcher, ticker *time.Ticker, quit chan struct{}, redisClient *redis.Client) {
	for {
		select {
		case <-ticker.C:
			fmt.Print("tick\n")
			records, err := filescanner.ScanPass(w.PATH, redisClient)
			if err != nil {
				fmt.Printf("ERROR: could not perform scan pass: %s\n", err)
			} else {
				CheckAndApply(records)
				DumpRecords(records)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func SetupTicker(w *config.Watcher, duration int, redisClient *redis.Client) (chan struct{}, error) {
	ticker := time.NewTicker(time.Duration(duration) * time.Second)
	quit := make(chan struct{})

	go TimerFunc(w, ticker, quit, redisClient)
	return quit, nil
}
