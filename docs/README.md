# Cloud-native Experience Lab Workshop

## Prerequisites

## Project setup

The focus of this lab is not on the actual implementation of the service itself. However, kicking off
a cloud-native project in Go is pretty straight forward.

```bash
go mod init github.com/qaware/cloud-native-weather-golang 
touch main.go

go get -u gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

With this, you can now start to implement the required business logic of the microservice application.

## Crosscutting Concerns

## Containerization

In this step we now need to containerize the application. With Go, we can leverage a multi-stage approach:
the final runtime artifact is build during the image build stage and then used and copied in the final runtime
stage (very similar to Cloud-native Build Packs approach).

**Lab Instructions**
1. Create a `Dockerfile` and add the following two stages
    - Build the Go binary using the correct `golang` base
    - Assemble the final runtime image using `gcr.io/distroless/base-debian11` as base

<details>
  <summary markdown="span">Click to expand solution ...</summary>

```
FROM golang:1.17-bullseye as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN go build -o /go/bin/weather-service

FROM gcr.io/distroless/base-debian11

ENV GIN_MODE=release
ENV PORT=8080

COPY --from=build /go/src/app/templates /templates
COPY --from=build /go/src/app/favicon.ico /
COPY --from=build /go/bin/weather-service /

CMD ["/weather-service"]
```
</details>

## K8s Deployment

### Kustomize



<details>
  <summary markdown="span">Click to expand solution ...</summary>

</details>

### Helm Chart

_TODO_

<details>
  <summary markdown="span">Click to expand solution ...</summary>

```bash
# prepare the gh-pages branch to serve the Helm chart
git checkout --orphan gh-pages
git reset --hard
git commit --allow-empty -m "fresh and empty gh-pages branch"
git push origin gh-pages
```
</details>

## Continuous Development

A good and efficient developer experience (DevEx) is of utmost importance for any cloud-native software
engineer. Rule 1: stay local as long as possible. Rule 2: automate all required steps: compile and package the source code, containerize the artifact and deploy the K8s resources locally. Continuously. The tools _Tilt_ and _Skaffold_ can both be used to establish this continuous dev-loop. For this section, you
only need to choose and use one of them.

### Option a) Tilt

In this step we are going to use [Tilt](https://tild.dev) to build, containerize and deploy the application
continuously to a local Kubernetes environment.

**Lab Instructions**
1. Make sure Tilt is installed locally on the development machine
2. Write a `Tiltfile` that performs the following steps
    - Build the required Docker image on every change to the source code
    - Apply the DEV overlay resources using Kustomize
    - Create a local port forward to the weather service HTTP port

<details>
  <summary markdown="span">Click to expand solution ...</summary>

Depending on your local K8s environment, the final `Tiltfile` might look slighty different.
```python
# -*- mode: Python -*-
# allow_k8s_contexts('rancher-desktop')

# to disable push with rancher desktop we need to use custom_build instead of docker_build
# docker_build('cloud-native-weather-golang', '.', dockerfile='Dockerfile')
custom_build('cloud-native-weather-golang', 'docker build -t $EXPECTED_REF .', ['./'], disable_push=True)

k8s_yaml(kustomize('./k8s/overlays/dev/'))
k8s_resource(workload='weather-service', port_forwards=[port_forward(18080, 8080, 'HTTP API')], labels=['Golang'])
```

To see of everything is working as expected issue the following command: `tilt up`
</details>

### Option b) Skaffold

In this step we are going to use [Skaffold](https://skaffold.dev) to build, containerize and deploy the application
continuously to a local Kubernetes environment.

**Lab Instructions**
1. Make sure Skaffold is installed locally on the development machine
2. Write a `skaffold.yaml` that performs the following steps
    - Build the required Docker image on every change to the source code
    - Apply the DEV overlay resources using Kustomize
    - Create a local port forward to the weather service HTTP port

<details>
  <summary markdown="span">Click to expand solution ...</summary>

The 3 steps of building, deployment and port-forwarding can all be codified in the
`skaffold.yaml` descriptor file.

```yaml
apiVersion: skaffold/v2beta24
kind: Config
metadata:
  name: weather-service-golang

# required for building the image
build:
  tagPolicy:
    gitCommit: {}
  artifacts:
    - image: cloud-native-weather-golang
      docker:
        dockerfile: Dockerfile
  local:
    push: false
    useBuildkit: true
    useDockerCLI: false

# required to deplo DEV overlay to default namespace
deploy:
  kustomize:
    defaultNamespace: default
    paths: ["k8s/overlays/dev"]

# create a local port-forward
portForward:
  - resourceName: weather-service
    resourceType: service
    namespace: default
    port: 8080
    localPort: 18080
```

To see of everything is working as expected issue the following command: `skaffold dev --no-prune=false --cache-artifacts=false`

</details>

## Continuous Integration

For any software project there must be a CI tool that takes care of continuously building and testing the produced software artifacts on every change.

### Github Actions

In this step we are going to use Github actions to build and test the application on every change. Also we are going to
leverage Github actions to perform 3rd party dependency checks as well as building and pushing the Docker image.

**Lab Instructions**
1. Create a Github action for each of the following tasks
    - Build the project on every change on main branch and every pull request
    - Build and push the Docker image to Github packages main branch, every pull request and tags
    - (_optional_) Perform CodeQL scans on main branch and every pull request
    - (_optional_) Perform a dependency review on every pull request

<details>
  <summary markdown="span">Click to expand solution ...</summary>

For each of the tasks, open the Github actions tab for the repository in your browser. Choose 'New workflow'. 

In the list of predefined actions, choose the **Go - Build a Go project** action. Adjust the suggested YAML
file content and commit.
```yaml
name: 'Go Build'

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.18

    - name: Build
      run: go build -v ./...

    - name: Test
      run: go test -v ./...
```

Next, choose the **Publish Docker Container** action from the Continuous integration section. Adjust the suggested YAML file content and commit.
```yaml
name: 'Docker Publish'

on:
  push:
    branches: [ "main" ]
    tags: [ 'v*.*.*' ]
  pull_request:
    branches: [ "main" ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:

    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
      id-token: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
        
      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.18

      - name: Build
        run: go build -v ./...

      - name: Test
        run: go test -v ./...

      # Install the cosign tool except on PR
      # https://github.com/sigstore/cosign-installer
      - name: Install cosign
        if: github.event_name != 'pull_request'
        uses: sigstore/cosign-installer@main
        with:
          cosign-release: 'v1.9.0'

      # Workaround: https://github.com/docker/build-push-action/issues/461
      - name: Setup Docker buildx
        uses: docker/setup-buildx-action@v2

      # Login against a Docker registry except on PR
      # https://github.com/docker/login-action
      - name: Log into registry ${{ env.REGISTRY }}
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v2
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # Extract metadata (tags, labels) for Docker
      # https://github.com/docker/metadata-action
      - name: Extract Docker metadata
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=semver,pattern={{major}}
            type=ref,event=branch
            type=raw,value=latest,enable={{is_default_branch}}

      # Build and push Docker image with Buildx (don't push on PR)
      # https://github.com/docker/build-push-action
      - name: Build and push Docker image
        id: build-and-push
        uses: docker/build-push-action@v3
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
```

Now repeat this process for the remaining two optional CI tasks of this lab.
</details>

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