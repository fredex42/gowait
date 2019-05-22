package kuberneteslink

import (
  "github.com/gowait/filescanner"
  "github.com/gowait/config"
  //appsv1 "k8s.io/api/apps/v1"
  v1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
  //"k8s.io/client-go/tools/clientcmd"
  //"k8s.io/client-go/rest"
)
/**
creates a kubernetes job in response to an event for the given watcher
*/
func CreateJobForEvent(k8clientSet *kubernetes.Clientset, runConfig *config.RunConfig, event *filescanner.WatchRecord) (*v1.Job, error) {
  jobsClient := k8clientSet.BatchV1().Jobs(runConfig.NAMESPACE)

  newJob := &v1.Job{
    ObjectMeta: metav1.ObjectMeta{
      GenerateName: "gowait-event-",
      Namespace: runConfig.NAMESPACE,
    },
    Spec: v1.JobSpec{
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
  }

  result, err := jobsClient.Create(newJob)
  if(err!=nil){
    return nil, err
  } else {
    return result, nil
  }
}
