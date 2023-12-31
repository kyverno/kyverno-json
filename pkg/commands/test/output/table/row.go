package table

type Row struct {
	IsFailure bool
	ID        int    `header:"id"`
	Policy    string `header:"policy"`
	Rule      string `header:"rule"`
	Payload   string `header:"payload"`
	Result    string `header:"result"`
	Reason    string `header:"reason"`
}
