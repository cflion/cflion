//  Copyright (c) 2018 The cflion Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package util

var empty = struct {
}{}

type Set struct {
	values map[interface{}]struct{}
}

func NewSet() *Set {
	return &Set{values: make(map[interface{}]struct{})}
}

func NewSetBySlice(a []interface{}) *Set {
	set := &Set{values: make(map[interface{}]struct{}, len(a))}
	for _, value := range a {
		set.values[value] = empty
	}
	return set
}

func (set *Set) ConvertToSlice() []interface{} {
	a := make([]interface{}, len(set.values))
	for k := range set.values {
		a = append(a, k)
	}
	return a
}

func DiffSet(set1 *Set, set2 *Set) (*Set, *Set) {
	diff1 := NewSet()
	diff2 := NewSet()
	for k := range set1.values {
		if _, ok := set2.values[k]; !ok {
			diff1.values[k] = empty
		}
	}
	for k := range set2.values {
		if _, ok := set1.values[k]; !ok {
			diff2.values[k] = empty
		}
	}
	return diff1, diff2
}
