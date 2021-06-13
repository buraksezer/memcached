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

package memcached

import (
	"context"
	"net"

	"github.com/buraksezer/memcached/config"
	"github.com/buraksezer/memcached/internal/tcp"
)

type Memcached struct {
	server *tcp.Server

	ctx    context.Context
	cancel context.CancelFunc
}

func New(c *config.Config) (*Memcached, error) {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Memcached{
		ctx:    ctx,
		cancel: cancel,
	}

	s, err := tcp.New(c.TCP, c.StartedCallback, m.dispatcher)
	if err != nil {
		return nil, err
	}

	m.server = s
	return m, nil
}

func (m *Memcached) dispatcher(conn net.Conn) error {
	return nil
}

func (m *Memcached) ListenAndServe() error {
	return m.server.ListenAndServe()
}

func (m *Memcached) Shutdown() error {
	m.cancel()

	return m.server.Shutdown()
}