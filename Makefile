NAME = cloud-native-weather-golang

default: build

image:
	@docker build -t $(NAME) .

build:
	@go build 

clean:
	@rm -f $(NAME)