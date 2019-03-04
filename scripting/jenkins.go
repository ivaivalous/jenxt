package scripting /* import "ivo.qa/jenxt/scripting" */

import (
	"errors"
	"io/ioutil"
	"ivo.qa/jenxt/config"
	"net/http"
	"net/url"
	"strings"
)

// ExecuteOnJenkins runs a script against a Jenkins server.
// This is done with two HTTP requests: one to get a crumb ID,
// and a second w=one with the actual payload.
// The function returns Jenkins' output of the script's execution.
func ExecuteOnJenkins(server *config.RemoteServer, script string) (response string, err error) {
	crumb, err := getCrumb(server)
	if err != nil {
		return "", err
	}

	form := url.Values{}
	form.Add("script", script)
	req, err := http.NewRequest("POST", getURL(server.URL), strings.NewReader(form.Encode()))
	req.Header.Add(ContentTypeHeader, ContentTypeForm)

	req.SetBasicAuth(server.Username, server.PasswordRaw)
	req.Header.Set("Jenkins-Crumb", crumb)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBytes, _ := ioutil.ReadAll(resp.Body)
	response = cleanResult(string(responseBytes))

	return
}

// getURL builds the URL to the Jenkins script API
func getURL(baseURL string) string {
	return baseURL + "/scriptText"
}

// getCrumbURL builds the URL to the Jenkins crumbs API
func getCrumbURL(baseURL string) string {
	return baseURL + "/crumbIssuer/api/xml?xpath=//crumb"
}

// getCrumb contacts Jenkins to get a crumb ID.
// This ID is used to avert some types of security risks.
// It needs to be sent with any API requests that follow.
// An error will be returned if the server is unreachable or
// in case authentication fails.
func getCrumb(server *config.RemoteServer) (crumb string, err error) {
	req, _ := http.NewRequest("GET", getCrumbURL(server.URL), nil)
	req.SetBasicAuth(server.Username, server.PasswordRaw)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("authentication failed")
		return
	}

	responseBytes, _ := ioutil.ReadAll(resp.Body)
	crumb = strings.TrimSpace(string(responseBytes))
	crumb = strings.Replace(crumb, "<crumb>", "", 1)
	crumb = strings.Replace(crumb, "</crumb>", "", 1)

	return crumb, nil
}

// cleanJesult removes unneeded data from Jenkins script execution responses.
// Responses from a script execution are typically prefixed with "Result: ".
// This function removes this and trims the output.
func cleanResult(response string) string {
	return strings.TrimSpace(strings.Replace(response, "Result: ", "", 1))
}
