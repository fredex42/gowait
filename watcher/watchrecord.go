package main

import "encoding/gob"
import "bytes"
import "github.com/go-redis/redis"
import "time"
import "os"
import "strings"

type WatchRecord struct {
  Path string
  Filename string
  LastMtime time.Time
  StableIterations int
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
  var keyname strings.Builder

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

  val, readErr := client.Get(keyname).Result()
  if(readErr != nil){
    if(readErr == redis.Nil){ //no record existed, create a new one
      newRecord := WatchRecord{}
      newRecord.Path = path
      newRecord.Filename = filename
      newRecord.LastMtime = time.Unix(0,0)
      newRecord.StableIterations = 0
      return &newRecord, nil
    } else {                  //something went kaboom
      return nil, readErr
    }
  } else {                  //we got data
    bytesval := []byte(val)
    newRecord,err := gobstring_to_watchrecord(&bytesval)
    if(err!=nil){
      return nil, err
    } else {
      return newRecord, nil
    }
  }
}

func update_watch_record(old_record *WatchRecord, current_state *os.FileInfo, client *redis.Client) (*WatchRecord, error){
  new_record := WatchRecord{}

  new_record.Path = old_record.Path
  new_record.Filename = old_record.Filename
  new_record.LastMtime = (*current_state).ModTime()
  if((*current_state).ModTime()==old_record.LastMtime){
    new_record.StableIterations=old_record.StableIterations+1
  } else {
    new_record.StableIterations = 0
  }

  keyname := get_keyname(new_record.Path, new_record.Filename)
  bytesToWritePtr, gobErr := watchrecord_to_gobstring(&new_record)
  if(gobErr!=nil){
    return nil, gobErr
  }
  keepTime, timeErr := time.ParseDuration("5m")
  if(timeErr!=nil){
    return nil, timeErr
  }

  _, writeErr := client.Set(keyname, *bytesToWritePtr, keepTime).Result()
  if(writeErr!=nil){
    return nil, writeErr
  }
  return &new_record, nil
}
