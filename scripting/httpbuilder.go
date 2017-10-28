package scripting

import (
	"fmt"
	"jenxt/config"
	"net/http"
)

const DEFAULT_LABEL = "default"

func (m Meta) GetHandler(c config.Configuration) (endpoint string, handler func(w http.ResponseWriter, r *http.Request)) {
	return m.getEndpoint(), func(w http.ResponseWriter, r *http.Request) {
		label := DEFAULT_LABEL
		if labelParameter := r.URL.Query().Get("label"); len(labelParameter) != 0 {
			label = string(labelParameter)
		}

		for _, s := range c.GetServersForLabel(label) {
			response, err := ExecuteOnJenkins(s, m.Script)
			if err != nil {
				fmt.Fprintf(w, s.Name+": "+err.Error())
			}

			fmt.Fprintf(w, s.Name+": "+response)
		}
	}
}

func (m Meta) PrintInfo() {
	info := fmt.Sprintf("Registered endpoint %s", m.getEndpoint())
	fmt.Println(info)
}

func (m Meta) getEndpoint() string {
	return fmt.Sprintf("/%s", m.Expose)
}
