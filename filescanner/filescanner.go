package filescanner

import "io/ioutil"
import "github.com/go-redis/redis"

/**
scans the given directory and returns a list of WatchRecord pointers, updated for
this pass.
*/
func ScanPass(dirname string, redisClient *redis.Client) ([]*WatchRecord, error) {
	files, err := ioutil.ReadDir(dirname)

	if err != nil {
		return nil, err
	}

	out_records := make([]*WatchRecord, len(files))

	for i, f := range files {
		//see https://golang.org/pkg/os/#FileInfo for description of "f"'s fields
		record, watchErr := get_watch_record_with_retry(dirname, f.Name(), redisClient)
		if watchErr != nil {
			return nil, watchErr
		}
		new_record, updateErr := update_watch_record(record, &f, redisClient)
		out_records[i] = new_record
		if updateErr != nil {
			return nil, updateErr
		}
	}
	return out_records, nil
}
