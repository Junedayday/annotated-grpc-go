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

package grpclog

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc/internal/grpclog"
)

// 四个级别的Log
type LoggerV2 interface {
	Info(args ...interface{})
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})
	Warning(args ...interface{})
	Warningln(args ...interface{})
	Warningf(format string, args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalln(args ...interface{})
	Fatalf(format string, args ...interface{})
	// V reports whether verbosity level l is at least the requested verbose level.
	V(l int) bool
}

// 提供可自己实现的接口
func SetLoggerV2(l LoggerV2) {
	grpclog.Logger = l
	grpclog.DepthLogger, _ = l.(grpclog.DepthLoggerV2)
}

// 四种级别的日志
const (
	infoLog int = iota
	warningLog
	errorLog
	fatalLog
)

var severityName = []string{
	infoLog:    "INFO",
	warningLog: "WARNING",
	errorLog:   "ERROR",
	fatalLog:   "FATAL",
}

// gRPC日志的默认实现
type loggerT struct {
	// log为官方的log包
	m []*log.Logger
	v int
}

// 定义了3个Writer，没有Fatal，是因为Fatal其实是error+os.Exit
// 0表示默认是Info级别
func NewLoggerV2(infoW, warningW, errorW io.Writer) LoggerV2 {
	return NewLoggerV2WithVerbosity(infoW, warningW, errorW, 0)
}

// 定义了四个日志，其中高等级的包含低等级的打印
// Tip io包中的MultiWriter和MultiReader，是一个很有用的工具，你可以看看其中的实现
func NewLoggerV2WithVerbosity(infoW, warningW, errorW io.Writer, v int) LoggerV2 {
	var m []*log.Logger
	m = append(m, log.New(infoW, severityName[infoLog]+": ", log.LstdFlags))
	m = append(m, log.New(io.MultiWriter(infoW, warningW), severityName[warningLog]+": ", log.LstdFlags))
	ew := io.MultiWriter(infoW, warningW, errorW)
	m = append(m, log.New(ew, severityName[errorLog]+": ", log.LstdFlags))
	m = append(m, log.New(ew, severityName[fatalLog]+": ", log.LstdFlags))
	return &loggerT{m: m, v: v}
}

// 默认的日志实现
func newLoggerV2() LoggerV2 {
	// 先初始化为不写
	errorW := ioutil.Discard
	warningW := ioutil.Discard
	infoW := ioutil.Discard

	// 根据环境变量的设置，默认为error
	logLevel := os.Getenv("GRPC_GO_LOG_SEVERITY_LEVEL")
	switch logLevel {
	case "", "ERROR", "error":
		errorW = os.Stderr
	case "WARNING", "warning":
		warningW = os.Stderr
	case "INFO", "info":
		infoW = os.Stderr
	}

	var v int
	vLevel := os.Getenv("GRPC_GO_LOG_VERBOSITY_LEVEL")
	if vl, err := strconv.Atoi(vLevel); err == nil {
		v = vl
	}
	return NewLoggerV2WithVerbosity(infoW, warningW, errorW, v)
}

func (g *loggerT) Info(args ...interface{}) {
	g.m[infoLog].Print(args...)
}

func (g *loggerT) Infoln(args ...interface{}) {
	g.m[infoLog].Println(args...)
}

func (g *loggerT) Infof(format string, args ...interface{}) {
	g.m[infoLog].Printf(format, args...)
}

func (g *loggerT) Warning(args ...interface{}) {
	g.m[warningLog].Print(args...)
}

func (g *loggerT) Warningln(args ...interface{}) {
	g.m[warningLog].Println(args...)
}

func (g *loggerT) Warningf(format string, args ...interface{}) {
	g.m[warningLog].Printf(format, args...)
}

func (g *loggerT) Error(args ...interface{}) {
	g.m[errorLog].Print(args...)
}

func (g *loggerT) Errorln(args ...interface{}) {
	g.m[errorLog].Println(args...)
}

func (g *loggerT) Errorf(format string, args ...interface{}) {
	g.m[errorLog].Printf(format, args...)
}

func (g *loggerT) Fatal(args ...interface{}) {
	g.m[fatalLog].Fatal(args...)
}

func (g *loggerT) Fatalln(args ...interface{}) {
	g.m[fatalLog].Fatalln(args...)
}

func (g *loggerT) Fatalf(format string, args ...interface{}) {
	g.m[fatalLog].Fatalf(format, args...)
}

func (g *loggerT) V(l int) bool {
	return l <= g.v
}


// 这个接口实现了打印了stack depth，也就是栈深度
type DepthLoggerV2 interface {
	InfoDepth(depth int, args ...interface{})
	WarningDepth(depth int, args ...interface{})
	ErrorDepth(depth int, args ...interface{})
	FatalDepth(depth int, args ...interface{})
}
