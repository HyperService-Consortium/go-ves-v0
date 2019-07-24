package wsrpc

// MessageType is defined as uint16
type MessageType = uint16

const (
	// CodeMessageRequest is from client to server, request for unicast its
	// message to other client
	CodeMessageRequest MessageType = iota

	// CodeMessageReply is from server to client
	CodeMessageReply

	// CodeRawProto if from client to server
	CodeRawProto

	// CodeClientHelloRequest is from client to server
	CodeClientHelloRequest

	// CodeClientHelloReply is from server to client
	CodeClientHelloReply // 4

	// CodeRequestComingRequest is from server to client
	CodeRequestComingRequest

	// CodeRequestComingReply is from client to server
	CodeRequestComingReply

	// CodeRequestGrpcServiceRequest is from client to server
	CodeRequestGrpcServiceRequest

	// CodeRequestGrpcServiceReply is from server to client
	CodeRequestGrpcServiceReply

	// CodeRequestNsbServiceRequest is from client to server
	CodeRequestNsbServiceRequest // 9

	// CodeRequestNsbServiceReply is from server to client
	CodeRequestNsbServiceReply

	// CodeSessionListRequest is from client to server
	CodeSessionListRequest

	// CodeSessionListReply is from server to client
	CodeSessionListReply

	// CodeTransactionListRequest is from client to server
	CodeTransactionListRequest

	// CodeTransactionListReply is from server to client
	CodeTransactionListReply // 14

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
	CodeSessionRequireTransactRequest // 19

	// CodeSessionRequireTransactReply is either from server to client or client to server
	CodeSessionRequireTransactReply

	// CodeAttestationReceiveRequest is either from server to client or client to server
	CodeAttestationReceiveRequest

	// CodeAttestationReceiveReply is either from server to client or client to server
	CodeAttestationReceiveReply

	// CodeAttestationSendingRequest is from server to client
	CodeAttestationSendingRequest

	// CodeAttestationSendingReply is client to server
	CodeAttestationSendingReply // 24

	// CodeCloseSessionRequest is from server to client
	CodeCloseSessionRequest

	// CodeCloseSessionReply is client to server
	CodeCloseSessionReply
)
