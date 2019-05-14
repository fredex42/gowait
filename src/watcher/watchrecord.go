package main

import "encoding/gob"
import "bytes"
import "github.com/go-redis/redis"
import "time"
import "os"

type WatchRecord struct {
  path string
  filename string
  last_mtime time.Time
  stable_iterations int
}

/**
serialize a watchrecord into a byte string
*/
func watchrecord_to_gobstring(record *WatchRecord) (*[]byte,error) {
  b := bytes.Buffer{}
  e := gob.NewEncoder(&b)
  err := e.Encode(record)
  if(err!=nil){ return nil, err }

  rtn := b.Bytes()
  return &rtn, nil
}

/**
deserialize from byte string to watchrecord
*/
func gobstring_to_watchrecord(bytestring *[]byte) (*WatchRecord, error) {
  record := WatchRecord{}
  buffer := bytes.NewBuffer(*bytestring)

  decoder := gob.NewDecoder(buffer)
  err := decoder.Decode(&record)
  if(err!=nil){
    return nil, err
  }

  return &record, nil
}

func get_keyname(path string, filename string) (string) {
  var keyname string.Builder

  keyname.WriteString(path)
  keyname.WriteString(":")
  keyname.WriteString(filename)

  return keyname.String()
}

/**
look up the path in Redis. Return nil if we have no record
*/
func get_watch_record_with_retry(path string, filename string, client *redis.Client) (*WatchRecord, error){
  keyname := get_keyname(path, filename)
  
  val, readErr := client.Get(keyname.String()).Result()
  if(readErr != nil){
    if(readErr == redis.Nil){ //no record existed, create a new one
      newRecord := WatchRecord{}
      newRecord.path = path
      newRecord.filename = filename
      newRecord.last_mtime = nil
      newRecord.stable_iterations = 0
      return &newRecord, 0
    } else {                  //something went kaboom
      return nil, readErr
    }
  } else {                  //we got data
    newRecord := gobstring_to_watchrecord(val)
    return &newRecord, 0
  }
}

func update_watch_record(old_record *WatchRecord, current_state *os.FileInfo, client *redis.Client) (*WatchRecord, error){
  return nil, nil
}
