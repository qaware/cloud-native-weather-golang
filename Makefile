NAME = cloud-native-weather-golang

default: build

image:
	@docker build -t $(NAME) .

build:
	@go build 

test:
	@go test -v -cover

clean:
	@rm -f $(NAME)