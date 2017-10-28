package scripting

import (
	"errors"
	"io/ioutil"
	"jenxt/config"
	"net/http"
	"net/url"
	"strings"
)

func ExecuteOnJenkins(server config.RemoteServer, script string) (response string, err error) {
	crumb, err := getCrumb(server)
	if err != nil {
		return "", err
	}

	form := url.Values{}
	form.Add("script", script)
	req, err := http.NewRequest("POST", getURL(server.URL), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(server.Username, server.Password)
	req.Header.Set("Jenkins-Crumb", crumb)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBytes, _ := ioutil.ReadAll(resp.Body)
	response = string(responseBytes)

	return
}

func getURL(baseURL string) string {
	return baseURL + "/scriptText"
}

func getCrumbURL(baseURL string) string {
	return baseURL + "/crumbIssuer/api/xml?xpath=//crumb"
}

func getCrumb(server config.RemoteServer) (crumb string, err error) {
	req, _ := http.NewRequest("GET", getCrumbURL(server.URL), nil)
	req.SetBasicAuth(server.Username, server.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("Authentication failed")
		return
	}

	responseBytes, _ := ioutil.ReadAll(resp.Body)
	crumb = strings.TrimSpace(string(responseBytes))
	crumb = strings.Replace(crumb, "<crumb>", "", 1)
	crumb = strings.Replace(crumb, "</crumb>", "", 1)

	return crumb, nil
}
