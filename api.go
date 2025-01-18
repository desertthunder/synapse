package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Credentials struct {
	username string
	password string
}

type Client struct {
	Credentials *Credentials
}

func SetEnvironmentVariables() Credentials {
	file, err := os.Open(".env")
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to open .env %v", err.Error()))
	}
	defer file.Close()

	c := Credentials{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "=")
		if parts[0] == "BLUESKY_USERNAME" {
			c.username = parts[1]
			os.Setenv(parts[0], c.username)
		} else if parts[0] == "BLUESKY_PASSWORD" {
			c.password = parts[1]
			os.Setenv(parts[0], c.password)
		}
	}

	if err = sc.Err(); err != nil {
		DefaultLogger().Fatal(fmt.Sprintf("sww %v", err.Error()))
	}

	return c
}
