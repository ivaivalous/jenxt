package scripting

import (
	"fmt"
	"jenxt/config"
	"net/http"
)

func (m Meta) GetHandler(s config.RemoteServer) (endpoint string, handler func(w http.ResponseWriter, r *http.Request)) {
	return m.getEndpoint(), func(w http.ResponseWriter, r *http.Request) {
		response, err := ExecuteOnJenkins(s, m.Script)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		fmt.Fprintf(w, response)
	}
}

func (m Meta) PrintInfo() {
	info := fmt.Sprintf("Registered endpoint %s", m.getEndpoint())
	fmt.Println(info)
}

func (m Meta) getEndpoint() string {
	return fmt.Sprintf("/%s", m.Expose)
}
