package playground

type Request struct {
	Payload       string   `json:"payload"`
	Preprocessors []string `json:"preprocessors"`
	Policy        string   `json:"policy"`
}
