# Cloud-native Weather Service with Golang

This example implements a simple weather REST service using Golang, GORM and Gin.

![Weather Service Architecture](architecture.png)

## Build and run locally

```bash
$ tilt up
$ skaffold dev --no-prune=false --cache-artifacts=false
```

## Exercise the application

```bash
$ curl -X GET http://localhost:18080/api/weather\?city\=Rosenheim
{"city":"Rosenheim","weather":"Clear"}

$ curl -X GET http://localhost:18080/
```

## Maintainer

M.-Leander Reimer (@lreimer), <mario-leander.reimer@qaware.de>

## License

This software is provided under the Apache v2.0 open source license, read the `LICENSE`
file for details.
