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
	"log"
	"testing"
	"time"

	"github.com/buraksezer/memcached/config"
	"github.com/buraksezer/memcached/internal/testutils"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestMemcached_ListenAndServe(t *testing.T) {
	port, err := testutils.GetFreePort()
	require.NoError(t, err)

	c := &config.Config{
		TCP: &config.TCP{
			BindAddr: "127.0.0.1",
			BindPort: port,
		},
	}

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(10*time.Second),
	)
	c.StartedCallback = func() {
		defer cancel()
		log.Printf("Memcached is ready to accept connections")
	}

	m, err := New(c)
	require.NoError(t, err)

	var errGr errgroup.Group
	errGr.Go(func() error {
		return m.ListenAndServe()
	})

	if <-ctx.Done(); true {
		require.NotErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	}

	err = m.Shutdown()
	require.NoError(t, err)
}
