package test

import (
	"log/slog"
	"net/http"
	"sync"
)

var (
	handler     http.Handler
	handlerLock sync.RWMutex
)

func SetHandler(h http.Handler) {
	handlerLock.Lock()
	defer handlerLock.Unlock()

	if h == nil {
		slog.Error("Trying to set nil handler")
		return
	}

	slog.Info("Setting handler for tests")
	handler = h
}

func Handler() http.Handler {
	handlerLock.RLock()
	defer handlerLock.RUnlock()

	if handler == nil {
		slog.Error("Handler is nil, it may not have been properly initialized")
	}

	return handler
}
