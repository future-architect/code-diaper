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
	"testing"
)

var input1 = crawler.Repositories(
	crawler.Repositories{
		crawler.Repository{
			URL:   "https://github.com/ghost/dummy1",
			Owner: "ghost",
			Name:  "dummy1",
			HitFiles: crawler.Files{
				{URL: "https://github.com/ghost/dummy1/fizz1.md", Fragments: []string{"abcdef"}},
				{URL: "https://github.com/ghost/dummy1/buzz1.md", Fragments: []string{"abcdef"}},
			},
		}, crawler.Repository{
			URL:   "https://github.com/ghost/dummy2",
			Owner: "ghost",
			Name:  "dummy2",
			HitFiles: crawler.Files{
				{URL: "https://github.com/ghost/dummy2/fizz2.md", Fragments: []string{"abcdef"}},
				{URL: "https://github.com/ghost/dummy2/buzz2.md", Fragments: []string{"abcdef"}},
			},
		}, crawler.Repository{
			URL:   "https://github.com/future-architect/vuls",
			Owner: "future-architect",
			Name:  "vuls",
			HitFiles: crawler.Files{
				{URL: "https://github.com/future-architect/vuls/main.go", Fragments: []string{"abcdef"}},
			},
		}, crawler.Repository{
			URL:   "https://github.com/future-architect/uroborosql",
			Owner: "future-architect",
			Name:  "uroborosql",
			HitFiles: crawler.Files{
				{URL: "https://github.com/future-architect/uroborosql/src/main/java/jp/co/future/uroborosql/utils/util.java", Fragments: []string{"* FutureTask.java\n* \n* Copyright (c) 2000-2019 Example Corporation."}},
			},
		},
	},
)

func TestDoNoFilter(t *testing.T) {
	expected := 4
	actual := NewSkipFilter([]condition.Sentence{}, []string{}, []string{}, []string{}).Do(input1)

	if len(actual) != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

func TestRepositoryFilter(t *testing.T) {
	expected := 1
	actual := NewSkipFilter([]condition.Sentence{}, []string{"dummy1"}, []string{"uroborosql/utils"}, []string{"ghost"}).Do(input1)

	if len(actual) != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

func TestSearchWord(t *testing.T) {
	expected1 := 0
	actual1 := NewSkipFilter([]condition.Sentence{"Copyright 2019 Future Corporation"}, []string{}, []string{}, []string{}).Do(input1)
	if len(actual1) != expected1 {
		t.Errorf("got: %v\nwant: %v", actual1, expected1)
	}

	expected2 := 1
	actual2 := NewSkipFilter([]condition.Sentence{"Copyright 2019 Example Corporation"}, []string{}, []string{}, []string{}).Do(input1)
	if len(actual2) != expected2 {
		t.Errorf("got: %v\nwant: %v", actual2, expected2)
	}

}
