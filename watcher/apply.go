package watcher

import (
	"github.com/fredex42/gowait/config"
	"github.com/fredex42/gowait/filescanner"
	"github.com/fredex42/gowait/kuberneteslink"
	"github.com/go-redis/redis"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

func ApplySinglek8s(k8clientSet *kubernetes.Clientset, redisClient *redis.Client, record *filescanner.WatchRecord, runConfig *config.K8RunConfig) error {
	job, err := kuberneteslink.CreateJobForEvent(k8clientSet, runConfig, record)
	if err != nil {
		log.Print("Could not create Kubernetes job: ", err)
	} else {
		strings.Builder{}
		redisClient.Set("job:")
	}
}

func CheckAndApply(records []*filescanner.WatchRecord, redisClient *redis.Client, config *config.Watcher, k8clientSet *kubernetes.Clientset) error {
	for i := range records {
		rec := records[i]

		if rec.StableIterations >= config.STABLE {
			log.Print(rec, " is stable for longer than ", config.STABLE, "; triggering")
			err := ApplySinglek8s(k8clientSet, redisClient, rec, &config.RUNCONFIG)
			if err != nil {
				log.Print("Could not run command for ", rec, ": ", err)
				return err //FIXME: should accumulate and return a custom error
			}
		}
	}
	return nil
}
