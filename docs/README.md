# Cloud-native Experience Lab Workshop

## Project setup

```bash
go mod init github.com/qaware/cloud-native-weather-golang 
touch main.go

go get -u gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

## Crosscutting Concerns

## Containerization

## K8s Deployment

### Kustomize

### Helm Chart

git checkout --orphan gh-pages
git reset --hard
git commit --allow-empty -m "fresh and empty gh-pages branch"
git push origin gh-pages

## Continuous Development

### Tilt

### Skaffold

## Continuous Integration

### Github Actions

## Continuous Deployment

## Flux2

In this step we are going to deploy the application using Flux2 onto the lab cluster. Flux will manage
the whole lifecycle, from initial deployment to automatic updates in case of changes and new versions.

**Lab Instructions**
1. Clone the Gitops repository for your cluster and create a dedicated apps directory
2. Create a dedicated K8s namespace resource
3. Install the weather service into the apps namespace using Kustomize
    - Patch the deployment and set `replicas: 2`
    - Patch the service and set `type: LoadBalancer`
4. (_optional_) Setup the image update automation workflow with suitable image repository and policy

<details>
  <summary markdown="span">Click to expand solution ...</summary>

First, we need to onboard and integrate the application with the Gitops workflow and repository.
```bash
# clone the experience lab Gitops repository
git clone https://github.com/qaware/cloud-native-explab.git
# create dedicated apps directory
take applications/bare/microk8s-cloudkoffer/weather-service-golang/
# initialize Kustomize descriptor
kustomize create
```

Create a `weather-namespace.yaml` file with the following content in the apps GitOps directory.
Do not forget to register the file resource in your `kustomization.yaml`.
```yaml
kind: Namespace
apiVersion: v1
metadata:
    name: weather-golang
```

Next, create the relevant Flux2 resources, such as `GitRepository` and `Kustomization` for the application.
```bash
flux create source git cloud-native-weather-golang \
    --url=https://github.com/qaware/cloud-native-weather-golang \
    --branch=main \
    --interval=5m0s \
    --export > weather-source.yaml

flux create kustomization cloud-native-weather-golang \
    --source=GitRepository/cloud-native-weather-golang \
    --path="./k8s/overlays/dev" \
    --prune=true \
    --interval=5m0s \
    --target-namespace=weather-golang \
    --export > weather-kustomization.yaml
```

The desired environment specific patches need to be added manually to the `weather-kustomization.yaml`, e.g.
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
    - apiVersion: v1
      kind: Service
      metadata:
        name: weather-service
      spec:
        type: LoadBalancer
```

Finally, add and configure image repository and policy for the image update automation to work.
```bash
flux create image repository cloud-native-weather-golang \
    --image=ghcr.io/qaware/cloud-native-weather-golang \
    --interval 1m0s \
    --export > weather-registry.yaml

flux create image policy cloud-native-weather-golang \
    --image-ref=cloud-native-weather-golang \
    --select-semver="1.2.x" \
    --export > weather-policy.yaml
```

Once all files have been created and modified, Git commit and push everything and watch the cluster
and Flux do the magic.

```bash
# to manually trigger the GitOps process use the following commands
flux reconcile source git flux-system
flux reconcile kustomization applications
flux get all
```
</details>