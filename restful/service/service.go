package service

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	uuid "github.com/satori/go.uuid"
)

var (
	// For serving
	DefaultId      = uuid.NewV4().String()
	DefaultName    = "daemon"
	DefaultVersion = "latest"
	DefaultAddress = ":8080"
)

type Service interface {
	Init(c *Meta) error
	Run() error
	String() string

	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

func NewService(name, version string) Service {
	return newService(name, version)
}

type service struct {
	opts Options

	mux *http.ServeMux

	sync.Mutex
	running bool
	exit    chan chan error
}

func newService(name, ver string) Service {
	options := newOptions(name, ver)
	s := &service{
		opts: options,
		mux:  http.NewServeMux(),
	}
	return s
}

func (s *service) Init(c *Meta) error {
	// 加载翻译文件
	log.Printf("Initializing %s ...", s.String())
	rand.Seed(time.Now().UTC().UnixNano())

	// 初始化服务注册
	if len(c.ServiceAddress) > 0 {
		s.opts.Address = c.ServiceAddress
	}
	if len(c.ServiceNamespace) > 0 {
		s.opts.Namespace = c.ServiceNamespace
	}
	return nil
}

func (s *service) Run() error {
	log.Printf("Starting %s ...", s.String())

	if err := s.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	select {
	// wait on kill signal
	case <-ch:
		// wait on context cancel
	case <-s.opts.Context.Done():
	}
	return s.stop()
}

func (s *service) String() string {
	return fmt.Sprintf("Service %s(%s): %s", s.opts.Name, s.opts.Version, s.opts.Id)
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Product Daemon API",
			Description: "Resource for Product Daemon",
			Version:     "2.0.0",
		},
	}
	swo.Schemes = []string{"http"}
}

func (s *service) Handle(pattern string, handler http.Handler) {
	wc, ok := handler.(*restful.Container)
	if ok && wc != nil {
		config := restfulspec.Config{
			WebServices:                   wc.RegisteredWebServices(),
			APIPath:                       "/apidocs.json",
			PostBuildSwaggerObjectHandler: enrichSwaggerObject,
		}
		wc.Add(restfulspec.NewOpenAPIService(config))

		cors := restful.CrossOriginResourceSharing{
			ExposeHeaders: []string{"X-My-Header"},
			AllowedHeaders: []string{"Content-Type", "Accept", "Content-Length",
				"Accept-Encoding", "X-CSRF-Token", "Authorization", "Access-Control-Allow-Headers", "auth-session"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			CookiesAllowed: true,
			Container:      wc,
		}
		wc.Filter(cors.Filter)

		// http://localhost:8080/apidocs.json
		http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("."))))
	}
	s.mux.Handle(pattern, handler)
}

func (s *service) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *service) start() error {
	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	l, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.opts.Address = l.Addr().String()

	var h http.Handler
	h = s.mux
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	var httpSrv *http.Server
	httpSrv = &http.Server{}
	httpSrv.Handler = h

	go httpSrv.Serve(l)

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	s.exit = make(chan chan error, 1)
	s.running = true

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()

	log.Printf("Listening on %v\n", l.Addr().String())
	return nil
}

func (s *service) stop() error {
	s.Lock()
	defer s.Unlock()

	if !s.running {
		return nil
	}

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			return err
		}
	}

	ch := make(chan error, 1)
	s.exit <- ch
	s.running = false

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			if chErr := <-ch; chErr != nil {
				return chErr
			}
			return err
		}
	}
	return <-ch
}
