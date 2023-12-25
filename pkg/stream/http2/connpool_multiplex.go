/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http2

import (
	"context"
	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	str "mosn.io/mosn/pkg/stream"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/utils"
	"mosn.io/pkg/variable"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	Init = iota
	Connecting
	Connected
	GoAway // received GoAway frame
)

var (
	defaultMaxConn         = 1
	connNumberLimit uint64 = 65535 // port limit
)

func init() {
	conn := os.Getenv("HTTP2_HOST_MAX_CONNECTIONS")
	if len(conn) > 0 {
		num, _ := strconv.Atoi(conn)
		if num > 0 && num <= int(connNumberLimit) {
			defaultMaxConn = num
			log.DefaultLogger.Infof("[stream] [http2] max host connections %d", defaultMaxConn)
		}
	}
}

// poolMultiplex is used for multiplex protocols like sofa, dubbo, etc.
// a single pool is connections which can be reused in a single host
type poolMultiplex struct {
	*connPool

	clientMux              sync.RWMutex
	activeClients          []*activeClientWrapper
	currentCheckAndInitIdx int64

	shutdown bool // pool is already shutdown
}

type activeClientWrapper struct {
	activeClient

	pool *poolMultiplex

	state       uint32 // for async connection
	indexInPool int
}

func (ac *activeClientWrapper) OnEvent(event api.ConnectionEvent) {
	ac.pool.onConnectionEvent(ac, event)
}

func isValidMaxNum(maxConns uint64) bool {
	// xDS cluster if not limit max connection will recv:
	// max_connections:{value:4294967295}  max_pending_requests:{value:4294967295}  max_requests:{value:4294967295}  max_retries:{value:4294967295}
	// if not judge max, will oom
	return maxConns > 0 && maxConns < connNumberLimit
}

// SetDefaultMaxConnNumPerHostPortForMuxPool set the max connections for each host:port
// users could use this function or cluster threshold config to configure connection no.
func SetDefaultMaxConnNumPerHostPortForMuxPool(maxConns int) {
	if isValidMaxNum(uint64(maxConns)) {
		defaultMaxConn = maxConns
	}
}

func NewConnMultiplexPool(ctx context.Context, host types.Host) types.ConnectionPool {
	pool := &connPool{
		tlsHash: host.TLSHashValue(),
	}
	pool.host.Store(host)
	return NewPoolMultiplex(pool)
}

// NewPoolMultiplex generates a multiplex conn pool
func NewPoolMultiplex(p *connPool) types.ConnectionPool {
	maxConns := p.Host().ClusterInfo().ResourceManager().Connections().Max()

	// the cluster threshold config has higher privilege than global config
	// if a valid number is provided, should use it
	if !isValidMaxNum(maxConns) {
		// override maxConns by default max conns
		maxConns = uint64(defaultMaxConn)
	}

	return &poolMultiplex{
		connPool:      p,
		activeClients: make([]*activeClientWrapper, maxConns),
	}
}

func (p *poolMultiplex) init(index int) {
	utils.GoWithRecover(func() {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("[stream] [http2] [connpool] init host %s", p.Host().AddressString())
		}

		p.clientMux.Lock()
		defer p.clientMux.Unlock()

		// if the pool is already shut down, do nothing directly return
		if p.shutdown {
			return
		}
		ctx := context.Background()
		client := p.newActiveClient(ctx, p.connPool)
		if client != nil {
			client.state = Connected
			client.indexInPool = index
			p.activeClients[index] = client
		} else {
			p.activeClients[index] = nil
		}
	}, nil)
}

func (p *poolMultiplex) newActiveClient(ctx context.Context, pool *connPool) *activeClientWrapper {
	ac := &activeClientWrapper{
		pool: p,
		activeClient: activeClient{
			pool: pool,
		},
	}

	host := pool.Host()
	data := host.CreateConnection(ctx)
	data.Connection.AddConnectionEventListener(ac)
	_ = variable.Set(ctx, types.VariableConnectionID, data.Connection.ID())
	codecClient := pool.createStreamClient(ctx, data)
	codecClient.SetStreamConnectionEventListener(ac)

	ac.client = codecClient
	ac.host = data

	if err := data.Connection.Connect(); err != nil {
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			log.DefaultLogger.Debugf("http2 underlying connection error: %v", err)
		}
		return nil
	}

	host.HostStats().UpstreamConnectionTotal.Inc(1)
	host.HostStats().UpstreamConnectionActive.Inc(1)
	host.ClusterInfo().Stats().UpstreamConnectionTotal.Inc(1)
	host.ClusterInfo().Stats().UpstreamConnectionActive.Inc(1)

	// bytes total adds all connections data together, but buffered data not
	codecClient.SetConnectionCollector(host.ClusterInfo().Stats().UpstreamBytesReadTotal, host.ClusterInfo().Stats().UpstreamBytesWriteTotal)

	return ac
}

// CheckAndInit init the connection pool
func (p *poolMultiplex) CheckAndInit(ctx context.Context) bool {
	var clientIdx int64 = 0 // most use cases, there will only be 1 connection
	if len(p.activeClients) > 1 {
		if clientIdx = getClientIDFromDownStreamCtx(ctx); clientIdx == invalidClientID {
			clientIdx = atomic.AddInt64(&p.currentCheckAndInitIdx, 1) % int64(len(p.activeClients))
			// set current client index to downstream context
			_ = variable.Set(ctx, types.VariableConnectionPoolIndex, clientIdx)
		}
	}

	var client *activeClientWrapper

	p.clientMux.RLock()
	client = p.activeClients[clientIdx]
	if client == nil {
		// unlock read lock
		p.clientMux.RUnlock()
		client = p.loadOrPutClient(client, clientIdx)
	} else {
		// unlock read lock
		p.clientMux.RUnlock()
	}

	if atomic.LoadUint32(&client.state) == Connected {
		return true
	}

	// init connection when client is Init or GoAway.
	if atomic.CompareAndSwapUint32(&client.state, Init, Connecting) ||
		atomic.CompareAndSwapUint32(&client.state, GoAway, Connecting) {
		p.init(int(clientIdx))
	}

	return false
}

func (p *poolMultiplex) loadOrPutClient(client *activeClientWrapper, clientIdx int64) *activeClientWrapper {
	p.clientMux.Lock()
	client = p.activeClients[clientIdx]
	if client == nil {
		client = &activeClientWrapper{pool: p,
			activeClient: activeClient{pool: p.connPool},
		}
		client.state = Init
		p.activeClients[clientIdx] = client
	}
	// unlock write lock
	p.clientMux.Unlock()
	return client
}

// NewStream Create a client stream and call's by proxy
func (p *poolMultiplex) NewStream(ctx context.Context, responseDecoder types.StreamReceiveListener) (types.Host, types.StreamSender, types.PoolFailureReason) {
	var (
		ok        bool
		clientIdx int64 = 0
	)

	if len(p.activeClients) > 1 {
		clientIdxInter, _ := variable.Get(ctx, types.VariableConnectionPoolIndex)
		if clientIdx, ok = clientIdxInter.(int64); !ok {
			// this client is not initialized
			return p.Host(), nil, types.ConnectionFailure
		}
	}

	ac := p.activeClients[clientIdx]

	host := p.Host()
	if ac == nil {
		return host, nil, types.ConnectionFailure
	}

	if atomic.LoadUint32(&ac.state) != Connected {
		return host, nil, types.ConnectionFailure
	}

	if !host.ClusterInfo().ResourceManager().Requests().CanCreate() {
		host.HostStats().UpstreamRequestPendingOverflow.Inc(1)
		host.ClusterInfo().Stats().UpstreamRequestPendingOverflow.Inc(1)
		return host, nil, types.Overflow
	}

	atomic.AddUint64(&ac.totalStream, 1)
	host.HostStats().UpstreamRequestTotal.Inc(1)
	host.HostStats().UpstreamRequestActive.Inc(1)
	host.ClusterInfo().Stats().UpstreamRequestTotal.Inc(1)
	host.ClusterInfo().Stats().UpstreamRequestActive.Inc(1)
	host.ClusterInfo().ResourceManager().Requests().Increase()
	streamEncoder := ac.client.NewStream(ctx, responseDecoder)
	streamEncoder.GetStream().AddEventListener(ac)

	_ = variable.Set(ctx, types.VariableUpstreamConnectionID, ac.client.ConnID())

	return host, streamEncoder, ""
}

// Shutdown stop the keepalive, so the connection will be idle after requests finished
func (p *poolMultiplex) Shutdown() {
	utils.GoWithRecover(func() {
		{
			p.clientMux.Lock()
			if p.shutdown {
				p.clientMux.Unlock()
				return
			}
			p.shutdown = true
			p.clientMux.Unlock()
		}
	}, nil)
}

func (p *poolMultiplex) createStreamClient(context context.Context, connData types.CreateConnectionData) str.Client {
	return p.connPool.createStreamClient(context, connData)
}

func (p *poolMultiplex) Close() {
	for i := 0; i < len(p.activeClients); i++ {
		cc := p.activeClients[i]
		if cc != nil {
			cc.client.Close()
		}
	}
}

func (p *poolMultiplex) onConnectionEvent(ac *activeClientWrapper, event api.ConnectionEvent) {
	host := p.Host()
	// event.ConnectFailure() contains types.ConnectTimeout and types.ConnectTimeout
	if event.IsClose() {
		host.HostStats().UpstreamConnectionClose.Inc(1)
		host.HostStats().UpstreamConnectionActive.Dec(1)

		host.ClusterInfo().Stats().UpstreamConnectionClose.Inc(1)
		host.ClusterInfo().Stats().UpstreamConnectionActive.Dec(1)

		switch event {
		case api.LocalClose:
			host.HostStats().UpstreamConnectionLocalClose.Inc(1)
			host.ClusterInfo().Stats().UpstreamConnectionLocalClose.Inc(1)

			if ac.closeWithActiveReq {
				host.HostStats().UpstreamConnectionLocalCloseWithActiveRequest.Inc(1)
				host.ClusterInfo().Stats().UpstreamConnectionLocalCloseWithActiveRequest.Inc(1)
			}

		case api.RemoteClose:
			host.HostStats().UpstreamConnectionRemoteClose.Inc(1)
			host.ClusterInfo().Stats().UpstreamConnectionRemoteClose.Inc(1)

			if ac.closeWithActiveReq {
				host.HostStats().UpstreamConnectionRemoteCloseWithActiveRequest.Inc(1)
				host.ClusterInfo().Stats().UpstreamConnectionRemoteCloseWithActiveRequest.Inc(1)

			}
		default:
			// do nothing
		}
		// only delete the active client when the state is not GoAway
		// since the goaway state client has already been overwritten.
		if atomic.LoadUint32(&ac.state) != GoAway {
			p.clientMux.Lock()
			p.activeClients[ac.indexInPool] = nil
			p.clientMux.Unlock()
		}
	} else if event == api.ConnectTimeout {
		host.HostStats().UpstreamRequestTimeout.Inc(1)
		host.ClusterInfo().Stats().UpstreamRequestTimeout.Inc(1)
		ac.client.Close()
	} else if event == api.ConnectFailed {
		host.HostStats().UpstreamConnectionConFail.Inc(1)
		host.ClusterInfo().Stats().UpstreamConnectionConFail.Inc(1)
	}
}

const invalidClientID = -1

func getClientIDFromDownStreamCtx(ctx context.Context) int64 {
	clientIdxInter, _ := variable.Get(ctx, types.VariableConnectionPoolIndex)
	clientIdx, ok := clientIdxInter.(int64)
	if !ok {
		return invalidClientID
	}
	return clientIdx
}
