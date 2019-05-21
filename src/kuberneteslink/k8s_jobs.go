package kuberneteslink

import (
  "config"
  "k8s.io/client-go/kubernetes"
  batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	scheme "k8s.io/client-go/kubernetes/scheme"
  rest "k8s.io/client-go/rest"
)
/**
creates a kubernetes job in response to an event for the given watcher
*/
func CreateJobForEvent(k8clientSet *kubernetes.ClientSet, runConfig *config.RunConfig, event *main.WatchRecord) (*v1.Job, error) {
  jobsClient := clientset.BatchV1().Jobs()

  newJob := &batchv1.Job{
    ObjectMeta: metav1.ObjectMeta{
      GenerateName: "gowait-event-",
      Namespace: runConfig.NAMESPACE,
    },
    Spec: batchv1.JobSpec{
      Template: apiv1.PodTemplateSpec{
        ObjectMeta: metav1.ObjectMeta{
          GenerateName: "gowait-event-",
        },
        Spec: apiv1.PodSpec{
          Containers: []apiv1.Container{
            {
              Name: "gowait-runner",
              Image: runConfig.IMAGE,
              Command: runConfig.COMMAND,
            },
          },
          RestartPolicy: apiv1.RestartPolicyOnFailure,
        },
      },
    },
  },

  result, err := jobsClient.CreateJob(newJob)
  if(err!=nil){
    return nil, err
  } else {
    return result
  }
}
