/*
 *
 * Copyright 2019 gRPC authors.
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

// 给各部件提供了k-v存储

// Package attributes defines a generic key/value store used in various gRPC
// components.
//
// All APIs in this package are EXPERIMENTAL.
package attributes

import "fmt"

// Key必须可以hashable，是为了可以存储到map中
// Tip Go的官方库是以简洁著称的，建议在复杂使用场景中，浅浅地封装一层
type Attributes struct {
	m map[interface{}]interface{}
}

// 按顺序写入map，格式为key1,value1,key2,value2
// 所以出现相同的key时，老的会覆盖新的
func New(kvs ...interface{}) *Attributes {
	if len(kvs)%2 != 0 {
		panic(fmt.Sprintf("attributes.New called with unexpected input: len(kvs) = %v", len(kvs)))
	}
	a := &Attributes{m: make(map[interface{}]interface{}, len(kvs)/2)}
	for i := 0; i < len(kvs)/2; i++ {
		a.m[kvs[i*2]] = kvs[i*2+1]
	}
	return a
}

// WithValues会对Attributes进行扩容，注意，这里是新建出对应的Attributes
func (a *Attributes) WithValues(kvs ...interface{}) *Attributes {
	if len(kvs)%2 != 0 {
		panic(fmt.Sprintf("attributes.New called with unexpected input: len(kvs) = %v", len(kvs)))
	}
	n := &Attributes{m: make(map[interface{}]interface{}, len(a.m)+len(kvs)/2)}
	for k, v := range a.m {
		n.m[k] = v
	}
	for i := 0; i < len(kvs)/2; i++ {
		n.m[kvs[i*2]] = kvs[i*2+1]
	}
	return n
}

func (a *Attributes) Value(key interface{}) interface{} {
	return a.m[key]
}
