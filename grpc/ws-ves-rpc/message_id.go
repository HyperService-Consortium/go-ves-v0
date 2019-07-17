package wsrpc

type MessageType = uint16

const (
	// request/reply
	CodeMessageRequest MessageType = iota
	CodeMessageReply
)
