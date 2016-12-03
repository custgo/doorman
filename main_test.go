package main

//
// add token:
// curl -v -d "username=yourname&password=yourpasswd" http://yourhost/token
//

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heiing/logs"
)

// 测试示例，需要适当修改才能运行
func TestMain(t *testing.T) {
	config, _ := readConfigFile(filepath.Join(logs.GetExecPath(), "config.json"))
	baseUrl := "http://127.0.0.1" + config.Listen
	body := strings.NewReader("username=hzm&password=123")
	res1, err := http.Post(baseUrl+config.TokenEndpoint, "application/x-www-form-urlencoded", body)
	if nil != err {
		t.Fatal("Error", err)
	}
	if 201 != res1.StatusCode {
		t.Fatal("Expect", 201, "Actual", res1.StatusCode)
	}
	res1_bytes, _ := ioutil.ReadAll(res1.Body)
	token := strings.TrimSpace(string(res1_bytes))
	t.Log("Token Created:", token)

	time.Sleep(1 * time.Second)

	loginUrl := baseUrl + config.LoginEndpoint + "?token=" + token + "&url=%2Findex.html"
	t.Log("Login URL:", loginUrl)
	req, _ := http.NewRequest("GET", loginUrl, nil)
	client := &http.Client{
		CheckRedirect: StopRedirect,
	}
	res2, err := client.Do(req)

	if nil != err {
		t.Log(err)
	}

	if nil == res2 {
		t.Fatal("Error: nil Response")
	}

	if 302 != res2.StatusCode {
		t.Fatal("Expect", 302, "Actual", res2.StatusCode)
	}
	cookies := res2.Cookies()
	for _, cookie := range cookies {
		t.Log("Cookie:", cookie.Name, "=", cookie.Value)
	}
}
