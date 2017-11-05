package scripting

import (
	"encoding/json"
	"fmt"
	"jenxt/config"
	"net/http"
)

const (
	DEFAULT_LABEL       = "default"
	CONTENT_TYPE_HEADER = "Content-Type"
	CT_JSON             = "application/json"
	CT_FORM             = "application/x-www-form-urlencoded"
)

// ExecutionResult is the main element in Jenxt
// execution responses. It lists the name of the server
// the request was made to, as well as a string or JSON
// response body.
type ExecutionResult struct {
	ServerName   string      `json:"server"`
	ResponseBody interface{} `json:"response"`
	IsError      bool        `json:"error,omitempty"`
}

// FullResult describes the holder element of
// the Jenxt execution response body.
type FullResult struct {
	Results []ExecutionResult `json:"results"`
}

// append adds the information of a server's response to
// a FinalResult list of results
func (f *FullResult) append(serverName string, response interface{}) {
	f.Results = append(f.Results, ExecutionResult{
		ServerName:   serverName,
		ResponseBody: response,
	})
}

// appendError adds the information of a server's error response to
// a FinalResult list of results
func (f *FullResult) appendError(serverName string, response interface{}) {
	f.Results = append(f.Results, ExecutionResult{
		ServerName:   serverName,
		ResponseBody: response,
		IsError:      true,
	})
}

// toJson converts an arbitrary interface to a JSON string
func toJson(p interface{}) string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

// GetHandler runs a handler dynamically based on the requested path
func GetHandler(c config.Configuration, scripts *Scripts) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		GenerateHandler(c, r.URL.Path, scripts)(w, r)
	}
}

// GenerateHandler generates a handler dynamically based on the requested path
func GenerateHandler(c config.Configuration, path string, scripts *Scripts) func(w http.ResponseWriter, r *http.Request) {
	if script, ok := getScriptByPath(path, scripts); ok {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(CONTENT_TYPE_HEADER, CT_JSON)

			label := DEFAULT_LABEL
			if labelParameter := r.URL.Query().Get("label"); len(labelParameter) != 0 {
				label = string(labelParameter)
			}

			result := FullResult{}
			for _, s := range c.GetServersForLabel(label) {
				response, err := ExecuteOnJenkins(s, script.Script)
				if err != nil {
					result.appendError(s.Name, err.Error())
					continue
				}

				if script.JSONResponse {
					result.append(s.Name, convertResponseToJSON(response))
				} else {
					result.append(s.Name, response)
				}
			}

			fmt.Fprintf(w, toJson(result))
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "404")
	}
}

// getScriptByPath gets the script matching a requested path.
// If no matching script could be found, ok will be false.
func getScriptByPath(path string, scripts *Scripts) (meta Meta, ok bool) {
	// This should be improved with some memorization
	for _, script := range *scripts {
		if script.getEndpoint() == path {
			return script, true
		}
	}

	return Meta{}, false
}

// convertResponseToJSON builds a map out or an arbitrary JSON-formatted
// string. Normally, Jenkins responses are just strings. The user may wish to parse
// a response into JSON so they can work with it directly (without having to parse it
// separately) in the calling program.
// Make sure `response` is JSON-formatted, otherwise the method might output a nil.
func convertResponseToJSON(response string) map[string]interface{} {
	responseBytes := []byte(response)
	var asMap map[string]interface{}

	json.Unmarshal(responseBytes, &asMap)
	return asMap
}

// PrintInfo outputs information of an endpoint that has been added.
// It's used in the main application logic to let the user know
// what scripts have been loaded.
func (m Meta) PrintInfo() {
	info := fmt.Sprintf("Registered endpoint %s", m.getEndpoint())
	fmt.Println(info)
}

// getEndpoint gets the endpoint path described in a Meta
func (m Meta) getEndpoint() string {
	return fmt.Sprintf("/%s", m.Expose)
}
