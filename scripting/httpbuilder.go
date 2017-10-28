package scripting

import (
	"fmt"
	"net/http"
)

func (m Meta) GetHandler() (endpoint string, handler func(w http.ResponseWriter, r *http.Request)) {
	return m.getEndpoint(), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, m.Expose)
	}
}

func (m Meta) getEndpoint() string {
	return fmt.Sprintf("/%s", m.Expose)
}
