package handler

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"

type Handler struct {
	logger infrastructure.Logger
	trace  infrastructure.ErrorTrace
}

func New(
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
) *Handler {
	//overriding the nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}

	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}
	var h = &Handler{
		logger: logger,
		trace:  trace,
	}
	return h
}
