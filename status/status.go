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

// Package status implements errors returned by gRPC.  These errors are
// serialized and transmitted on the wire between server and client, and allow
// for additional data to be transmitted via the Details field in the status
// proto.  gRPC service handlers should return an error created by this
// package, and gRPC clients should expect a corresponding error to be
// returned from the RPC call.
//
// This package upholds the invariants that a non-nil error may not
// contain an OK code, and an OK code must result in a nil error.
package status

import (
	"context"
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/status"
)

// Status references google.golang.org/grpc/internal/status. It represents an
// RPC status code, message, and details.  It is immutable and should be
// created with New, Newf, or FromProto.
// https://godoc.org/google.golang.org/grpc/internal/status

// status 错误代码的相关实现
// Tip 这里的代码目录结构可以学习一下：
// 外部可被引用的部分，相对简洁
// 内部的实现放在internal中，避免外部包的引用
type Status = status.Status

// Tip 可以对比error包中的Wrap，思考一下各自适合什么场景
func New(c codes.Code, msg string) *Status {
	return status.New(c, msg)
}

func Newf(c codes.Code, format string, a ...interface{}) *Status {
	return New(c, fmt.Sprintf(format, a...))
}

// 这里的Err()里值得看看
func Error(c codes.Code, msg string) error {
	return New(c, msg).Err()
}

func Errorf(c codes.Code, format string, a ...interface{}) error {
	return Error(c, fmt.Sprintf(format, a...))
}

func ErrorProto(s *spb.Status) error {
	return FromProto(s).Err()
}

// 将原始的proto类型的错误，包装到status的错误中
func FromProto(s *spb.Status) *Status {
	return status.FromProto(s)
}

// 从error中提取status信息
// 这里并不是根据类型来判断，而是通过是否定义了接口
// Tips 用接口定义的方式扩展性很高
func FromError(err error) (s *Status, ok bool) {
	if err == nil {
		return nil, true
	}
	if se, ok := err.(interface {
		GRPCStatus() *Status
	}); ok {
		return se.GRPCStatus(), true
	}
	return New(codes.Unknown, err.Error()), false
}

// FromError需要处理2个返回值，这个函数是为了方便，
// Tips 两个返回值和一个返回值各有利弊，建议好好思考利弊
// 例如单返回值带来的便利性，但返回的信息量，例如具体错误原因，会被屏蔽
func Convert(err error) *Status {
	s, _ := FromError(err)
	return s
}

// 查询错误里的Code
func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if se, ok := err.(interface {
		GRPCStatus() *Status
	}); ok {
		return se.GRPCStatus().Code()
	}
	return codes.Unknown
}

// 集成了常用的Context两种错误情况
func FromContextError(err error) *Status {
	switch err {
	case nil:
		return nil
	case context.DeadlineExceeded:
		return New(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return New(codes.Canceled, err.Error())
	default:
		return New(codes.Unknown, err.Error())
	}
}

