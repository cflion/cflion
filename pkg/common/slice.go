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

package common

func DistinctIntSlice(a []int) []int {
	m := make(map[int]struct{}, len(a))
	for _, key := range a {
		m[key] = struct{}{}
	}
	r := make([]int, 0, len(m))
	i := 0
	for k := range m {
		r[i] = k
		i++
	}
	return r
}

func DistinctInt64Slice(a []int64) []int64 {
	m := make(map[int64]struct{}, len(a))
	for _, key := range a {
		m[key] = struct{}{}
	}
	r := make([]int64, 0, len(m))
	i := 0
	for k := range m {
		r[i] = k
		i++
	}
	return r
}

func DiffTwoInt64Slice(a1 []int64, a2 []int64) ([]int64, []int64) {
	return DiffLeftInt64Slice(a1, a2), DiffLeftInt64Slice(a2, a1)
}

func DiffLeftInt64Slice(a1 []int64, a2 []int64) []int64 {
	m := make(map[int64]struct{}, len(a2))
	for _, k := range a2 {
		m[k] = struct{}{}
	}
	diff := make([]int64, 0, len(a1))
	for _, k := range a1 {
		if _, ok := m[k]; !ok {
			diff = append(diff, k)
		}
	}
	return diff
}
