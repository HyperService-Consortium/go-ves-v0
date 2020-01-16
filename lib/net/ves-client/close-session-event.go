package vesclient

type SessionCloseHandler = func(sessionID []byte)

type sessionCloseSubscriber struct {
	Aborted bool
	handler SessionCloseHandler
}

func newSessionCloseSubscriber(handler SessionCloseHandler) *sessionCloseSubscriber {
	return &sessionCloseSubscriber{handler: handler}
}

func (ses *sessionCloseSubscriber) Abort() bool {
	a := ses.Aborted
	ses.Aborted = true
	return a
}

func (ses *sessionCloseSubscriber) IsAborted() bool {
	return ses.Aborted
}

func (ses *sessionCloseSubscriber) Emit(sessionID []byte) {
	ses.handler(sessionID)
}

type SessionCloseSubscriber interface {
	Abort() bool
	IsAborted() bool
	Emit(sessionID []byte)
}

func (vc *VesClient) SubscribeCloseSession(handler SessionCloseHandler) SessionCloseSubscriber {
	var subscriber = newSessionCloseSubscriber(handler)
	vc.closeSessionRWMutex.Lock()
	vc.closeSessionSubscriber = append(vc.closeSessionSubscriber, subscriber)
	vc.closeSessionRWMutex.Unlock()
	return subscriber
}

func (vc *VesClient) emitClose(SessionId []byte) {
	vc.closeSessionRWMutex.Lock()
	var j = len(vc.closeSessionSubscriber)
	for i := j - 1; i >= 0; i-- {
		subscriber := vc.closeSessionSubscriber[i]
		if subscriber.IsAborted() {
			continue
		}
		j--
		vc.closeSessionSubscriber[j] = subscriber
		subscriber.Emit(SessionId)
	}
	vc.closeSessionSubscriber = vc.closeSessionSubscriber[j:]
	vc.closeSessionRWMutex.Unlock()
}
