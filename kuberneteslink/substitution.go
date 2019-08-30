package kuberneteslink

type TemplateSubs struct {
  Namespace string
  Cmd []string
  Environment map[string]string
  Image string
}
