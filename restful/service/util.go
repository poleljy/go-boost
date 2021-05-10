package service

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"github.com/poleljy/go-boost/restful/errors"
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

func InvalidRequest(res *restful.Response, err error) error {
	log.Println(err)
	res.Header().Add("Content-Type", "application/json")
	return res.WriteError(http.StatusOK, errors.BadRequest(`invalid request: %s`, err.Error()))
}

func NotFound(res *restful.Response, err error) error {
	log.Println(err)
	res.Header().Add("Content-Type", "application/json")
	return res.WriteError(http.StatusOK, errors.NotFound(`%s`, err.Error()))
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

func WriteOKEntity(res *restful.Response) error {
	res.Header().Add("Content-Type", "application/json")
	return res.WriteError(http.StatusOK, errors.New("", http.StatusOK))
}
