package scan

type Request struct {
	Payload       any      `json:"payload"`
	Preprocessors []string `json:"preprocessors"`
}
