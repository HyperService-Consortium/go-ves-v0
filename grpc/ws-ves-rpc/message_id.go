package wsrpc

// MessageType is defined as uint16
type MessageType = uint16

const (
	// CodeMessageRequest is from client to server, request for unicast its
	// message to other client
	CodeMessageRequest MessageType = iota

	// CodeMessageReply is from server to client
	CodeMessageReply

	// CodeClientHelloRequest is from server to client
	CodeClientHelloRequest

	// CodeClientHelloReply is from client to server
	CodeClientHelloReply

	// CodeRequestComing is from server to client
	CodeRequestComing

	// CodeRequestGrpcServiceRequest is from client to server
	CodeRequestGrpcServiceRequest

	// CodeRequestGrpcServiceReply is from server to client
	CodeRequestGrpcServiceReply

	// CodeSessionListRequest is from client to server
	CodeSessionListRequest

	// CodeSessionListReply is from server to client
	CodeSessionListReply

	// CodeTransactionListRequest is from client to server
	CodeTransactionListRequest

	// CodeTransactionListReply is from server to client
	CodeTransactionListReply
)
