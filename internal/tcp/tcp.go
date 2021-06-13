// Copyright 2021 Burak Sezer
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcp

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/buraksezer/memcached/config"
)

type Handler func(conn net.Conn) error

type Server struct {
	listener net.Listener
	handler  Handler
	config   *config.TCP

	wg       sync.WaitGroup
	started  func()
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(c *config.TCP, started func(), handler Handler) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		handler: handler,
		started: started,
		config:  c,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer s.wg.Done()

	err := s.handler(conn)
	if err != nil {
		log.Printf("handler returned an error: %v", err)
	}
}

func (s *Server) ListenAndServe() error {
	addr := net.JoinHostPort(s.config.BindAddr, s.config.BindPort)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = l
	if s.started != nil {
		s.started()
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		s.wg.Add(1)
		go s.handleConn(conn)
	}
}

func (s *Server) Shutdown() error {
	select {
	case <-s.ctx.Done():
		// Already closed
		return nil
	default:
	}

	s.cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.listener.Close()
	}()

	s.wg.Wait()
	return <-errCh
}
