/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// 首先，这个package的含义为连通性
// 这里提供了一个文档，建议大家先仔细阅读一下
// Tip 这里的五种状态，是gRPC连接状态的核心，我们要理解这个状态机转换

// Package connectivity defines connectivity semantics.
// For details, see https://github.com/grpc/grpc/blob/master/doc/connectivity-semantics-and-api.md.
// All APIs in this package are experimental.
package connectivity

import (
	"context"

	"google.golang.org/grpc/grpclog"
)

type State int

func (s State) String() string {
	switch s {
	case Idle:
		return "IDLE"
	case Connecting:
		return "CONNECTING"
	case Ready:
		return "READY"
	case TransientFailure:
		return "TRANSIENT_FAILURE"
	case Shutdown:
		return "SHUTDOWN"
	default:
		grpclog.Errorf("unknown connectivity state: %d", s)
		return "Invalid-State"
	}
}

const (
	Idle State = iota
	Connecting
	Ready
	TransientFailure
	Shutdown
)

type Reporter interface {
	CurrentState() State
	WaitForStateChange(context.Context, State) bool
}
