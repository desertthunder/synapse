package main

import (
	"os"
	"strings"
	"testing"
)

func TestAPI(t *testing.T) {
	t.Run("test env setter", func(t *testing.T) {
		contents := strings.Join([]string{
			`BLUESKY_USERNAME = "name_of_your_bot"`,
			`BLUESKY_PASSWORD = "your_bot_password"`}, "\n")
		tmp, err := os.CreateTemp(".", ".env.test")
		if err != nil {
			t.Errorf("test setup failed, unable to create tmp env file %v", err.Error())
		}

		defer tmp.Close()

		numBytes, err := tmp.Write([]byte(contents))
		if err != nil {
			t.Errorf("test setup failed, unable to write to tmp env file %v", err.Error())
		} else {
			t.Logf("wrote %v bytes to %v", numBytes, tmp.Name())
		}

		gotLocal := SetEnvironmentVariables(tmp.Name())
		gotEnv := AtCredentials{
			Handle:   os.Getenv("BLUESKY_USERNAME"),
			Password: os.Getenv("BLUESKY_PASSWORD"),
		}

		want := `"name_of_your_bot"`
		if gotLocal.Handle != want && gotEnv.Handle != want {
			t.Errorf("wanted %v but got %v in mem and %v in env", want, gotLocal.Handle, gotEnv.Handle)
		}

		want = `"your_bot_password"`
		if gotLocal.Password != want && gotEnv.Password != want {
			t.Errorf("wanted %v but got %v in mem and %v in env", want, gotLocal.Password, gotEnv.Password)
		}

		os.Remove(tmp.Name())
	})
}
