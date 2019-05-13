type WatchRecord struct {
  path String
  filename String
  last_mtime time.Time
  stable_iterations int

}

/**
look up the path in Redis. Return nil if we have no record
*/
func get_watch_record_with_retry(path string, filename string) (*WatchRecord, error){

}

func update_watch_record(old_record *WatchRecord, current_state *os.FileInfo) (*WatchRecord, error){

}
