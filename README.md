# Helm Repository for Cloud-native Weather Service with Golang

## Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add <alias> https://qaware.github.io/cloud-native-weather-golang

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
<alias>` to see the charts.

To install the `cloud-native-weather-golang` chart:

    helm install my-cloud-native-weather-golang <alias>/cloud-native-weather-golang

To uninstall the chart:

    helm delete my-cloud-native-weather-golang