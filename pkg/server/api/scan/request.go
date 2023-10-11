package scan

type Request struct {
	Payload       interface{} `json:"payload"`
	Preprocessors []string    `json:"preprocessors"`
}
