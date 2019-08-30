package kuberneteslink

import (
  "github.com/gowait/filescanner"
  "github.com/gowait/config"
  //appsv1 "k8s.io/api/apps/v1"
  v1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
  "io/ioutil"
  "gopkg.in/yaml.v2"
  "strings"
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
              Environment: runConfig.ENVIRONMENT
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

func CreateJobFromTemplate(k8clientSet *kubernetes.Clientset, jobTemplate *v1.Job, substitutions TemplateSubs) (*v1.Job, err) {
  jobsClient := k8clientSet.BatchV1().Jobs(substitutions.Namespace)


}
/**
reads in a yaml template of a job spec
*/
func LoadTemplate(filename string) (*v1.Job, error) {
  rtn := v1.Job{}

  data, err := ioutil.ReadFile(filename)
  if(err!=nil){
    return nil, err
  }

  marshalErr := yaml.Unmarshal([]byte(data), &rtn)
  if(marshalErr!=nil){
    return nil, err
  }

  return &rtn, nil
}

func LoadTemplateWithSubs(filename string, substitutions TemplateSubs) (*v1.Job, error) {
  rtn := v1.Job{}

  data, err := ioutil.ReadFile(filename)
  if(err!=nil){
    return nil, err
  }

  replaced1 := strings.replaceAll(data, "{{ namespace }}", substitutions.Namespace)
  replaced2 := strings.replaceAll(replaced1, "{{ image }}", substitutions.Image)

  marshalErr := yaml.Unmarshal([]byte(replaced2), &rtn)
  if(marshalErr!=nil){
    return nil, err
  }

  rtn.Spec.Template.Spec.Containers[0].Command = substitutions.Command
  rtn.Spec.Template.Spec.Containers[0].Environment = substitutions.Environment

  return &rtn, nil
}
