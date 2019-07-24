package vesclient

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Myriad-Dreamin/go-ves/log"
)

type handler struct {
	funcs []func()
}

func (h *handler) register(atexit func()) {
	h.funcs = append(h.funcs, atexit)
}

func (h *handler) atExit() {
	osQuitSignalChan := make(chan os.Signal)
	signal.Notify(osQuitSignalChan, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT,
		syscall.SIGKILL, syscall.SIGILL, syscall.SIGTERM,
	)
	for {
		select {
		case osc := <-osQuitSignalChan:
			log.Infoln("handlering:", osc)
			for _, f := range h.funcs {
				f()
			}
			return
		}
	}
}

var phandler *handler

func init() {
	phandler = new(handler)
	go phandler.atExit()
}
