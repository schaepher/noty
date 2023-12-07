package channels

type Message struct {
	From      string
	Receivers []string
	Title     string
	Body      string
}
