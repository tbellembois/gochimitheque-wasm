package types

// ValidateConfig is a jQuery validate
// parameters struct as defined
// https://jqueryvalidation.org/validate/
type ValidateConfig struct {
	Ignore     string                     `json:"ignore"`
	ErrorClass string                     `json:"errorClass"`
	Rules      map[string]ValidateRule    `json:"rules"`
	Messages   map[string]ValidateMessage `json:"messages"`
}
type ValidateRule struct {
	Email    bool           `json:"email,omitempty"`
	EqualTo  string         `json:"equalTo,omitempty"`
	Required interface{}    `json:"required"`
	Remote   ValidateRemote `json:"remote"`
}
type ValidateRemote struct {
	URL        string                 `json:"url"`
	Type       string                 `json:"type"`
	BeforeSend interface{}            `json:"beforeSend"`
	Data       map[string]interface{} `json:"data"`
}

type ValidateMessage struct {
	Email     string `json:"email,omitempty"`
	EqualTo   string `json:"equalTo,omitempty"`
	Required  string `json:"required"`
	MinLength string `json:"minlength"`
}

func (jq Jquery) ValidateAddRequired() {

	jq.Object.Call("rules", "add", "required")

}

func (jq Jquery) ValidateRemoveRequired() {

	jq.Object.Call("rules", "remove", "required")

}

func (jq Jquery) Validate(config ValidateConfig) {

	configMap := StructToMap(config)
	jq.Object.Call("validate", configMap)

}

func (jq Jquery) Valid() bool {

	return jq.Object.Call("valid").Bool()

}
