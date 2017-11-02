package scripting

import (
	"encoding/json"
	"fmt"
	"jenxt/config"
	"net/http"
)

const DEFAULT_LABEL = "default"

// ExecutionResult is the main element in Jenxt
// execution responses. It lists the name of the server
// the request was made to, as well as a string or JSON
// response body.
type ExecutionResult struct {
	ServerName   string      `json:"server"`
	ResponseBody interface{} `json:"response"`
}

// FullResult describes the holder element of
// the Jenxt execution response body.
type FullResult struct {
	Results []ExecutionResult `json:"results"`
}

// toJson converts an arbitrary interface to a JSON string
func toJson(p interface{}) string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

// GetHandler builds a function that can be passed to http.HandleFunc.
// This creates the endpoints users can access.
// When the returned function is called, an HTTP request is made to
// the required Jenkins servers and a response is built listing
// all server responses.
func (m Meta) GetHandler(c config.Configuration) (endpoint string, handler func(w http.ResponseWriter, r *http.Request)) {
	return m.getEndpoint(), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		label := DEFAULT_LABEL
		if labelParameter := r.URL.Query().Get("label"); len(labelParameter) != 0 {
			label = string(labelParameter)
		}

		result := FullResult{}
		for _, s := range c.GetServersForLabel(label) {
			response, err := ExecuteOnJenkins(s, m.Script)
			if err != nil {
				result.Results = append(result.Results, ExecutionResult{ServerName: s.Name, ResponseBody: err.Error()})
				continue
			}

			if m.JSONResponse {
				result.Results = append(result.Results, ExecutionResult{ServerName: s.Name, ResponseBody: convertResponseToJSON(response)})
			} else {
				result.Results = append(result.Results, ExecutionResult{ServerName: s.Name, ResponseBody: response})
			}
		}

		fmt.Fprintf(w, toJson(result))
	}
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
