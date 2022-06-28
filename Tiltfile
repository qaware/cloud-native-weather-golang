# -*- mode: Python -*-
# allow_k8s_contexts('rancher-desktop')

# to disable push with rancher desktop we need to use custom_build instead of docker_build
# docker_build('cloud-native-weather-golang', '.', dockerfile='Dockerfile')
custom_build('cloud-native-weather-golang', 'docker build -t $EXPECTED_REF .', ['./'], disable_push=True)

k8s_yaml(kustomize('./k8s/overlays/dev/'))
k8s_resource(workload='weather-service', port_forwards=[port_forward(18080, 8080, 'HTTP API')], labels=['Golang'])
