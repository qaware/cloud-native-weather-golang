# Cloud-native Experience Lab Workshop

## Project setup

```bash
$ go mod init github.com/qaware/cloud-native-weather-golang 
$ go get -u gin-gonic/gin
$ go get -u gorm.io/gorm
$ go get -u gorm.io/driver/postgres
```

## Crosscutting Concerns

## Containerization

## K8s Deployment

## Continuous Development

### Tilt

### Skaffold

## Continuous Integration

### Github Actions

## Continuous Deployment

This application is part of the Cloud-native Experience Lab (for Software Engineers). It has been implemented 
in several different languages and frameworks (Go, JavaEE, .NET, NodeJS). The main purpose of this app is to 
showcase important cloud-native design and development principles.

**Lab Instructions**
1. Read the installation instructions
    - Weather Service: https://github.com/qaware/cloud-native-weather-golang
    - Weather UI: https://github.com/qaware/cloud-native-weather-vue3
2. Create a dedicated namespace via GitOps
2. Install the weather service into dedicated namespace using Kustomize
    - Patch the deployment and set `replicas: 2`
3. Install the weather UI into dedicated namespace using Kustomize
    - Patch the deployment and set `replicas: 2`
    - Patch the service and set `type: LoadBalancer`
3. (_optional_) Setup the image update automation workflow with suitable image repository and policy for the service and the UI

<details>
  <summary markdown="span">Click to expand solution ...</summary>

First, create a `weather-namespace.yaml` file with the following content and add it to
the applications GitOps directory like `applications/bare/microk8s-cloudkoffer/weather-service-golang/`.
```yaml
kind: Namespace
apiVersion: v1
metadata:
    name: weather-golang
```

Next, create the relevant Flux2 resources, such as `GitRepository` and `Kustomization` for the application.
```bash
cd applications/bare/microk8s-cloudkoffer

flux create source git cloud-native-weather-golang \
    --url=https://github.com/qaware/cloud-native-weather-golang \
    --branch=main \
    --interval=5m0s \
    --export > weather-service-golang/weather-source.yaml

flux create kustomization cloud-native-weather-golang \
    --source=GitRepository/cloud-native-weather-golang \
    --path="./k8s/overlays/dev" \
    --prune=true \
    --interval=5m0s \
    --target-namespace=weather-golang \
    --export > weather-service-golang/weather-kustomization.yaml
```

The Kustomize patches need to be added manually to the `weather-kustomization.yaml`.
```yaml
  images:
    - name: cloud-native-weather-golang
      newName: ghcr.io/qaware/cloud-native-weather-golang # {"$imagepolicy": "flux-system:cloud-native-weather-golang:name"}
      newTag: 1.2.0 # {"$imagepolicy": "flux-system:cloud-native-weather-golang:tag"}
  patchesStrategicMerge:
    - apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: weather-service
      spec:
        replicas: 2
```

Finally, add and configure image repository and policy for the image update automation to work.
```bash
flux create image repository cloud-native-weather-golang \
    --image=ghcr.io/qaware/cloud-native-weather-golang \
    --interval 1m0s \
    --export > weather-service-golang/weather-registry.yaml

flux create image policy cloud-native-weather-golang \
    --image-ref=cloud-native-weather-golang \
    --select-semver="1.2.x" \
    --export > weather-service-golang/weather-policy.yaml
```

</details>