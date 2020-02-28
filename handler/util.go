package handler

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type ErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"description"`
}

type RstResponse struct {
	User  string `json:"user"`
	Token string `json:"access_token"`
}

type TokenCheck struct {
	Result string `json:"result"`
	//ReturnMsg string   `json:"returnMsg"`
	UsrInfo UserInfo `json:"userInfo"`
}

type UserInfo struct {
	Id string `json:"userId"`
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var httpClientB = &http.Client{
	Transport: tr,
	Timeout:   0,
}

var httpClientG = &http.Client{
	Transport: httpClientB.Transport,
	Timeout:   time.Duration(10) * time.Second,
}

func getENV(env string) string {
	env_value := os.Getenv(env)

	if env_value == "" {
		fmt.Println("FATAL: NEED ENV", env)
		fmt.Println("Exit...........")
		os.Exit(2)
	}

	fmt.Println("ENV:", env, env_value)

	return env_value
}

func trRequest(method string, url string, users Users) (*http.Response, error) {
	var req *http.Request
	var err error

	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(users.Username, users.Password)
	req.Header.Set("X-CSRF-Token", "1")

	return tr.RoundTrip(req)

}

func getDFtoken(info string) string {
	if len(info) == 0 {
		return ""
	}

	tmpLs := strings.Split(info, "access_token=")

	if len(tmpLs) == 0 {
		return ""
	}

	rstLs := strings.Split(tmpLs[1], "&")

	if len(rstLs) == 0 {
		return ""
	}

	return rstLs[0]

}

func Get(client *redis.Client, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd("get", key)
	client.Process(cmd)
	return cmd
}
func Setkey(client *redis.Client, key string, value string) *redis.StringCmd {
	cmd := redis.NewStringCmd("set", key, value)
	client.Process(cmd)
	return cmd
}
