package main

import "io/ioutil"

/**
scans the given directory and returns a list of WatchRecord pointers, updated for
this pass.
*/
func scan_pass(dirname string) ([]WatchRecord, error){
  files, err := ioutil.ReadDir(dirname)

  if(err != nil) {
    return nil, err
  }

  out_record := make([]WatchRecord, len(files))

  for i, f := range(files) {
    //see https://golang.org/pkg/os/#FileInfo for description of "f"'s fields
    record, watchErr := get_watch_record_with_retry(dirname, f.Name())
    if(watchErr != nil){
      return nil, watchErr
    }
    out_record[i] = update_watch_record(record, f)
  }
  return out_record, nil
}
