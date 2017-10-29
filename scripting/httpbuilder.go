package scripting

import (
	"encoding/json"
	"fmt"
	"jenxt/config"
	"net/http"
)

const DEFAULT_LABEL = "default"

type ExecutionResult struct {
	ServerName   string      `json:"server"`
	ResponseBody interface{} `json:"response"`
}

type FullResult struct {
	Results []ExecutionResult `json:"results"`
}

func toJson(p interface{}) string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

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

func convertResponseToJSON(response string) map[string]interface{} {
	responseBytes := []byte(response)
	var asMap map[string]interface{}

	json.Unmarshal(responseBytes, &asMap)
	return asMap
}

func (m Meta) PrintInfo() {
	info := fmt.Sprintf("Registered endpoint %s", m.getEndpoint())
	fmt.Println(info)
}

func (m Meta) getEndpoint() string {
	return fmt.Sprintf("/%s", m.Expose)
}
