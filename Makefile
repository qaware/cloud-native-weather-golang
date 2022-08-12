NAME = cloud-native-weather-golang

default: build

image:
	@docker build -t $(NAME) .

build:
	@go build 

test: build
	@go test -v -cover -coverprofile=cov.out

clean:
	@rm -f $(NAME)