package init

import (
	"bytes"
	"fmt"
)

// Response to represent the response of stevedore Init
// swagger:model
type Response struct {
	// Namespace in which stevedore has initialized
	Namespace string `json:"namespace"`
	Message   string `json:"message"`
}

// Responses is collection of Response
type Responses []Response

// String to format the Responses
func (responses Responses) String() string {
	buff := bytes.NewBufferString("Stevedore initialised in below namespace(s):\n")
	for _, response := range responses {
		buff.WriteString(fmt.Sprintf("%s: %v\n", response.Namespace, response.Message))
	}
	return buff.String()
}
