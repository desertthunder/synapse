package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	CreateSessionMethod AtProtoMethod = "com.atproto.server.createSession"
	CreatePostMethod    AtProtoMethod = "com.atproto.repo.createRecord"
	ServiceURL          string        = "https://bsky.social"
)

type AtProtoMethod = string

type AtClient struct {
	Service     string
	Credentials AtCredentials
}

type SessionRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type AtCredentials struct {
	Handle          string
	Password        string
	AccessToken     string
	RefreshToken    string
	DID             string
	ServiceEndpoint string
}

type Service struct {
	ID              string `json:"id"`
	ServiceEndpoint string `json:"serviceEndpoint"`
	Type            string `json:"type"`
}

type VerificationMethod struct {
	Controller         string `json:"controller"`
	ID                 string `json:"id"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
	Type               string `json:"type"`
}

type DidDoc struct {
	Context            []string             `json:"@context"`
	AlsoKnownAs        []string             `json:"alsoKnownAs"`
	ID                 string               `json:"id"`
	Service            []Service            `json:"service"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
}

type Session struct {
	AccessJwt       string `json:"accessJwt"`
	RefreshJwt      string `json:"refreshJwt"`
	Handle          string `json:"handle"`
	Did             string `json:"did"`
	DidDoc          DidDoc `json:"didDoc"`
	Email           string `json:"email"`
	EmailConfirmed  bool   `json:"emailConfirmed"`
	EmailAuthFactor bool   `json:"emailAuthFactor"`
	Active          bool   `json:"active"`
	Status          string `json:"status"`
}

func SetEnvironmentVariables(f string) AtCredentials {
	file, err := os.Open(f)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to open %v %v", f, err.Error()))
	}
	defer file.Close()

	c := AtCredentials{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "=")
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		if k == "BLUESKY_USERNAME" {
			c.Handle = v
			err = os.Setenv(k, c.Handle)
			if err != nil {
				logger.Error(fmt.Sprintf("sww %v", err.Error()))
			}
		} else if k == "BLUESKY_PASSWORD" {
			c.Password = v
			err = os.Setenv(k, c.Password)
			if err != nil {
				logger.Error(fmt.Sprintf("sww %v", err.Error()))
			}
		}
	}

	if err = sc.Err(); err != nil {
		DefaultLogger().Fatal(fmt.Sprintf("sww %v", err.Error()))
	}

	return c
}

func GetCredentialsFromEnv() AtCredentials {
	handle := os.Getenv("BLUESKY_HANDLE")
	password := os.Getenv("BLUESKY_PASSWORD")

	return AtCredentials{
		Handle:   handle,
		Password: password,
	}
}

// Instantiate a new [AtClient]
func NewClient(c AtCredentials) *AtClient {
	return &AtClient{
		Service:     ServiceURL,
		Credentials: c,
	}
}

func (c *AtCredentials) SetSession(s Session) {
	c.AccessToken = s.AccessJwt
	c.RefreshToken = s.RefreshJwt
	c.DID = s.Did
	c.ServiceEndpoint = s.ServiceEndpoint()
}

func (c AtClient) BuildURL(path AtProtoMethod) string {
	return fmt.Sprintf("%s/xrpc/%s", c.Service, path)
}

func (c AtClient) CreateSession() (*Session, error) {
	u := c.BuildURL(CreateSessionMethod)
	r := SessionRequest{c.Credentials.Handle, c.Credentials.Password}
	s := Session{}
	buf := bytes.NewBuffer(nil)

	req, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("unable to build body: %s", err.Error())
	}

	buf.Write(req)

	rsp, err := http.Post(u, "application/json", buf)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate: %s", err.Error())
	} else {
		logger.Debugf("request to %v completed with status %v", u, rsp.Status)
		buf.Reset()
	}

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status %v", rsp.Status)
	}

	defer rsp.Body.Close()

	if _, err := buf.ReadFrom(rsp.Body); err != nil {
		return nil, fmt.Errorf("unable to read response %s", err.Error())
	}

	if err = json.Unmarshal(buf.Bytes(), &s); err != nil {
		return nil, fmt.Errorf("unable to marshal JSON %v", err.Error())
	} else {
		c.Credentials.SetSession(s)
		logger.Info(fmt.Sprintf("session created at %s", time.Now().Format("03:04 PM on 01/02/2006")))
	}

	return &s, nil
}

func (s Session) ServiceEndpoint() string {
	return s.DidDoc.Service[0].ServiceEndpoint
}

func (s Session) DebugToken(l int) string {
	return s.AccessJwt[:l/2] + "..." + s.AccessJwt[len(s.AccessJwt)-l/2:]
}

// function SaveToken stores the tokens in the database
//
// Note: Estimated expiry is split in half before storage.
// https://atproto.blue/en/latest/atproto_client/auth.html#session-string
func SaveTokens(s Session) {
	// Access Token
	// Expiry = Now + 1 hour
	// Refresh Token
	// Expiry = Now + 4 weeks
}

// function Login creates a [AtClient] and authenticates into BlueSky
func Login() {
	logger.Warn("make sure you use an app password to authenticate")
	cred := GetCredentialsFromEnv()

	if cred.Handle == "" {
		logger.Debug("no credentials in term env, attempting to set manually")
		cred = SetEnvironmentVariables(".env")
	}

	logger.Debugf("credentials set with handle: %v", cred.Handle)

	c := NewClient(cred)
	s, err := c.CreateSession()
	if err != nil {
		logger.Errorf("unable to create session %v", err.Error())
	}

	logger.Infof("session created with token %v", s.DebugToken(12))
}
