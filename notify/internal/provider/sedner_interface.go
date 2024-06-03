package provider

type Sender interface {
	Send(in SendIn) error
}

type SendIn struct {
	To    string
	Phone string
	From  string
	Title string
	Type  int
	Data  []byte
}
