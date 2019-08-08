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
package filter

import (
	"github.com/future-architect/code-diaper/condition"
	"github.com/future-architect/code-diaper/crawler"
	"strings"
)

type Filter interface {
	Do(rs crawler.Repositories) crawler.Repositories
}

type skipFilter struct {
	sList          []condition.Sentence
	skipRepoNames  []string
	skipLibNames   []string
	skipOwnerNames []string
}

func NewSkipFilter(sList []condition.Sentence, skipRepoNames, skipLibNames, skipOwnerNames []string) Filter {
	return skipFilter{
		sList:          sList,
		skipRepoNames:  skipRepoNames,
		skipLibNames:   skipLibNames,
		skipOwnerNames: skipOwnerNames,
	}
}

func (f skipFilter) Do(rs crawler.Repositories) crawler.Repositories {
	return f.doByPath(f.doByRepo(f.doByOwner(f.doBySearchWord(rs))))
}

func (f skipFilter) doBySearchWord(rs crawler.Repositories) crawler.Repositories {
	if len(f.sList) == 0 {
		return rs
	}
	result := rs
	for _, s := range f.sList {
		result = f.doBySentence(s, result) // update filtering by sentence
	}
	return result
}

func (f skipFilter) doBySentence(s condition.Sentence, rs crawler.Repositories) crawler.Repositories {
	if s == "" {
		return rs
	}
	searchWords := s.Parse()

	result := make(crawler.Repositories, 0, len(rs))

	// check one line
	for _, v := range rs {
		var files []crawler.File

		for _, f := range v.HitFiles {
			var containsAllKeyWord = false
			for _, fragment := range f.Fragments {

				// When there is line break, split the search target
				lines := strings.Split(fragment, "\n")

				for _, searchLine := range lines {
					if allContains(searchLine, searchWords) {
						containsAllKeyWord = true
						break
					}
				}
			}
			if containsAllKeyWord {
				files = append(files, f)
			}
		}

		if len(files) > 0 {
			result = append(result, crawler.Repository{
				URL:      v.URL,
				Owner:    v.Owner,
				Name:     v.Name,
				HitFiles: files,
			})
		}
	}
	return result
}

func (f skipFilter) doByRepo(rs crawler.Repositories) crawler.Repositories {
	if len(f.skipRepoNames) == 0 || f.skipRepoNames[0] == "" {
		return rs
	}

	// filter by repository name
	var deleteIndexes []int
	for i, r := range rs {
		if f.matchRepoName(r) {
			deleteIndexes = append(deleteIndexes, i)
		}
	}
	return rs.Exclude(deleteIndexes)
}

func (f skipFilter) doByOwner(rs crawler.Repositories) crawler.Repositories {
	if len(f.skipOwnerNames) == 0 || f.skipOwnerNames[0] == "" {
		return rs
	}

	// filter by repository owner name
	var result crawler.Repositories
	for _, r := range rs {
		if !f.matchOwnerName(r) {
			result = append(result, r)
		}
	}
	return result
}

func (f skipFilter) doByPath(rs crawler.Repositories) crawler.Repositories {
	if len(f.skipLibNames) == 0 || f.skipLibNames[0] == "" {
		return rs
	}

	// filter by filepath
	var result crawler.Repositories
	for _, r := range rs {

		var containsFiles crawler.Files
		for _, file := range r.HitFiles {
			if f.containsLib(file) {
				continue
			}
			containsFiles = append(containsFiles, file)
		}

		if len(containsFiles) >= 1 {
			r.HitFiles = containsFiles
			result = append(result, r)
		}
	}
	return result
}

func (f skipFilter) matchRepoName(repo crawler.Repository) bool {
	for _, v := range f.skipRepoNames {
		if repo.Name == v || repo.ForkSource == v {
			return true
		}
	}
	return false
}

func (f skipFilter) containsLib(file crawler.File) bool {
	for _, v := range f.skipLibNames {
		if strings.Contains(file.URL, v) {
			return true
		}
	}
	return false
}

func (f skipFilter) matchOwnerName(repo crawler.Repository) bool {
	for _, v := range f.skipOwnerNames {
		if repo.Owner == v {
			return true
		}
	}
	return false
}

func allContains(fragment string, searchWords []string) bool {
	for _, search := range searchWords {
		if !strings.Contains(fragment, search) {
			return false
		}
	}
	return true
}
