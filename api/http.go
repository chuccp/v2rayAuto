package api

import (
	"net/http"
	"strconv"
)

type tlsHttp struct {
	httpServer *http.Server
	hasCLose   bool
}

func (s *tlsHttp) Close() error {
	if s.httpServer != nil {
		s.hasCLose = true
		err := s.httpServer.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *tlsHttp) Start(port int, certFile string, keyFile string, handler http.Handler) error {
	s.hasCLose = false
	s.httpServer = &http.Server{Addr: ":" + strconv.Itoa(port), Handler: handler}
	err := s.httpServer.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		if s.hasCLose {
			return nil
		}
		return err
	}
	return nil
}
