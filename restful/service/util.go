package service

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/emicklei/go-restful"

	"deploy/restful/errors"
)

func ParseRequestBody(req *restful.Request, res *restful.Response, entityPointer interface{}) error {
	err := req.ReadEntity(entityPointer)
	if err != nil {
		log.Println(err)
		res.Header().Add("Content-Type", "application/json")
		res.WriteError(http.StatusOK, errors.BadRequest(`invalid request: %s`, err.Error()))
	}
	return err
}

// fail
func InternalError(res *restful.Response, err error) error {
	log.Println(err)
	res.Header().Add("Content-Type", "application/json")
	return res.WriteError(http.StatusOK, errors.New(err.Error(), http.StatusInternalServerError))
}

func BadRequest(res *restful.Response, err error) error {
	log.Println(err)
	res.Header().Add("Content-Type", "application/json")
	return res.WriteError(http.StatusOK, errors.BadRequest(`invalid request: %s`, err.Error()))
}

// success
func WriteEntity(res *restful.Response, value interface{}) error {
	res.Header().Add("Content-Type", "application/json")

	if value == nil {
		return res.WriteEntity(struct {
			ret string `json:"ret,omitempty"`
		}{})
	}
	return res.WriteEntity(value)
}

func WriteNullEntity(res *restful.Response) error {
	res.Header().Add("Content-Type", "application/json")
	return res.WriteEntity(struct {
		ret string `json:"ret,omitempty"`
	}{})
}

// setup a signal hander to gracefully exit
func SignalHandler() <-chan struct{} {
	stop := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			syscall.SIGINT,  // Ctrl+C
			syscall.SIGTERM, // Termination Request
			syscall.SIGSEGV, // FullDerp
			syscall.SIGABRT, // Abnormal termination
			syscall.SIGILL,  // illegal instruction
			syscall.SIGFPE)  // floating point - this is why we can't have nice things
		sig := <-c
		log.Printf("Signal (%v) detected, shutting down", sig)
		close(stop)
	}()
	return stop
}
