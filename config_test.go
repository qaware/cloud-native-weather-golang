package main

import (
	"testing"
)

func TestPort(t *testing.T) {
	port := port()
	if port != ":8080" {
		t.Fatalf(`Unexpected port %v`, port)
	}
}

func TestPostgresHost(t *testing.T) {
	host := postgresHost()
	if host != "localhost" {
		t.Fatalf(`Unexpected Postgres host %v`, host)
	}
}

func TestPostgresPort(t *testing.T) {
	port := postgresPort()
	if port != "5432" {
		t.Fatalf(`Unexpected Postgres port %v`, port)
	}
}

func TestPostgresUser(t *testing.T) {
	user := postgresUser()
	if user != "postgres" {
		t.Fatalf(`Unexpected Postgres user %v`, user)
	}
}

func TestPostgresPassword(t *testing.T) {
	password := postgresPassword()
	if password != "postgres" {
		t.Fatalf(`Unexpected Postgres password %v`, password)
	}
}

func TestWeatherUri(t *testing.T) {
	uri := weatherUri()
	if uri != "https://api.openweathermap.org" {
		t.Fatalf(`Unexpected Weather URI %v`, uri)
	}
}

func TestWeatherAppid(t *testing.T) {
	appid := weatherAppid()
	if appid != "5b3f51e527ba4ee2ba87940ce9705cb5" {
		t.Fatalf(`Unexpected App ID %v`, appid)
	}
}
