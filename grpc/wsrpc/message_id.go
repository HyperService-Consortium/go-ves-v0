package wsrpc

// MessageType is defined as uint16
type MessageType = uint16

const (
	// CodeMessageRequest is from client to server, request for unicast its
	// message to other client
	CodeMessageRequest MessageType = iota

	// CodeMessageReply is from server to client
	CodeMessageReply

	// CodeClientHelloRequest is from client to server
	CodeClientHelloRequest

	// CodeClientHelloReply is from server to client
	CodeClientHelloReply

	// CodeRequestComingRequest is from server to client
	CodeRequestComingRequest

	// CodeRequestComingReply is from client to server
	CodeRequestComingReply

	// CodeRequestGrpcServiceRequest is from client to server
	CodeRequestGrpcServiceRequest

	// CodeRequestGrpcServiceReply is from server to client
	CodeRequestGrpcServiceReply

	// CodeRequestNsbServiceRequest is from client to server
	CodeRequestNsbServiceRequest

	// CodeRequestNsbServiceReply is from server to client
	CodeRequestNsbServiceReply

	// CodeSessionListRequest is from client to server
	CodeSessionListRequest

	// CodeSessionListReply is from server to client
	CodeSessionListReply

	// CodeTransactionListRequest is from client to server
	CodeTransactionListRequest

	// CodeTransactionListReply is from server to client
	CodeTransactionListReply

	// CodeUserRegisterRequest is from client to server
	CodeUserRegisterRequest

	// CodeUserRegisterReply is from server to client
	CodeUserRegisterReply

	// CodeSessionFinishedRequest is either from server to client or client to server
	CodeSessionFinishedRequest

	// CodeSessionFinishedReply is either from server to client or client to server
	CodeSessionFinishedReply

	// CodeSessionRequestForInitRequest is from server to client
	// CodeSessionRequestForInitRequest
	//
	// CodeSessionRequestForInitReply

	// CodeSessionRequireTransactRequest is either from server to client or client to server
	CodeSessionRequireTransactRequest

	// CodeSessionRequireTransactReply is either from server to client or client to server
	CodeSessionRequireTransactReply

	// CodeAttestationReceiveRequest is either from server to client or client to server
	CodeAttestationReceiveRequest

	// CodeAttestationReceiveReply is either from server to client or client to server
	CodeAttestationReceiveReply
)
