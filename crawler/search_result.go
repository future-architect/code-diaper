/**
 * Copyright (c) 2019-present Future Corporation
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
 */
package crawler

type Repository struct {
	URL        string
	Owner      string
	Name       string
	HitFiles   Files
	ForkSource string // parent is the repository this repository was forked from, source is the ultimate source for the network. https://developer.github.com/v3/repos/#response-4
}

type Repositories []Repository

// Fragments is represents github api result
// https://developer.github.com/v3/search/#text-match-metadata
type File struct {
	URL       string
	Fragments []string
}

type Files []File

func (rs Repositories) Index(r Repository) (index int) {
	for i, elm := range rs {
		if elm.URL == r.URL {
			return i
		}
	}
	return -1
}

func (rs Repositories) Exclude(targetIndexes []int) Repositories {
	var result Repositories
	for i, r := range rs {
		if !Contains(targetIndexes, i) {
			result = append(result, r)
		}
	}
	return result
}

func (rs Repositories) Merge(r Repository) Repositories {
	idx := rs.Index(r)
	if idx == -1 {
		return append(rs, r)
	}

	var result Repositories
	result = append(result, rs...)

	for _, f := range r.HitFiles {
		result[idx].HitFiles = result[idx].HitFiles.Merge(f)
	}
	return result
}

func (fs Files) Index(file File) int {
	for i, element := range fs {
		if element.URL == file.URL {
			return i
		}
	}
	return -1
}

func (fs Files) Merge(f File) Files {
	idx := fs.Index(f)
	if idx == -1 {
		return append(fs, f)
	}

	var result Files
	result = append(result, fs...)

	result[idx].Fragments = append(result[idx].Fragments, f.Fragments...)
	return result
}

func Contains(arr []int, e int) bool {
	for _, v := range arr {
		if e == v {
			return true
		}
	}
	return false
}
