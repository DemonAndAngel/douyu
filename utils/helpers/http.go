package helpers

import (
	"io/ioutil"
	"net/http"
	urlPkg "net/url"
	"strings"
)

func Request(method, url string, headers map[string]string, content urlPkg.Values) (code int, body []byte, err error) {
	if headers == nil {
		headers = map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}
	}
	if method == "GET" {
		return requestGet(url, headers, content)
	} else {
		return requestPost(url, headers, content)
	}
}

func requestGet(url string, headers map[string]string, params urlPkg.Values) (code int, body []byte, err error) {
	client := &http.Client{}
	path := ""
	if len(params) > 0 {
		for k, _ := range params {
			if path != "" {
				path += "&"
			}
			path += k + "=" + urlPkg.QueryEscape(params.Get(k))
		}
	}
	if path != "" {
		a := strings.Split(url, "?")
		if len(a) > 1 {
			url += "&"
		} else {
			url += "?"
		}
		url += path
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	code = resp.StatusCode
	if code != http.StatusOK {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func requestPost(url string, headers map[string]string, content urlPkg.Values) (code int, body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(content.Encode()))
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	code = resp.StatusCode
	if code != http.StatusOK {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err = ioutil.ReadAll(resp.Body)
	return
}
