package main

import "os"

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

func postgresHost() string {
	host := os.Getenv("POSTGRES_HOST")
	if len(host) == 0 {
		host = "localhost"
	}
	return host
}

func postgresPort() string {
	return "5432"
}

func postgresUser() string {
	user := os.Getenv("POSTGRES_USER")
	if len(user) == 0 {
		user = "postgres"
	}
	return user
}

func postgresPassword() string {
	pwd := os.Getenv("POSTGRES_PASSWORD")
	if len(pwd) == 0 {
		pwd = "postgres"
	}
	return pwd
}

func postgresDb() string {
	db := os.Getenv("POSTGRES_DB")
	if len(db) == 0 {
		db = "weather"
	}
	return db
}

func weatherUri() string {
	uri := os.Getenv("WEATHER_URI")
	if len(uri) == 0 {
		uri = "https://api.openweathermap.org"
	}
	return uri
}

func weatherAppid() string {
	appid := os.Getenv("WEATHER_APPID")
	if len(appid) == 0 {
		appid = "5b3f51e527ba4ee2ba87940ce9705cb5"
	}
	return appid
}
