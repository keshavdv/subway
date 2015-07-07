package msg

type Envelope struct {
	Type    string
	Payload []byte
}

type Ping struct {
}

type Pong struct {
}
