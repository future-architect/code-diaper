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
package formatter

import (
	"github.com/future-architect/code-diaper/crawler"
	"strings"
	"testing"
)

var input1 = []SearchResult{
	{
		Query: "test1",
		Repos: crawler.Repositories{
			{
				URL:   "https://github.com/ghost/dummy-repo1",
				Owner: "ghost",
				Name:  "dummy-repo1",
				HitFiles: crawler.Files{
					{
						URL:       "https://github.com/ghost/dummy-repo1/dummy1.md",
						Fragments: []string{"detect dummy1-1", "detect dummy1-2"},
					},
					{
						URL:       "https://github.com/ghost/dummy-repo1/dummy2.md",
						Fragments: []string{"detect dummy2-1", "detect dummy2-2"},
					},
				},
			},
		},
		HitCount: 1,
	},
	{
		Query: "test2",
		Repos: crawler.Repositories{
			{
				URL:   "https://github.com/ghost/dummy-repo2",
				Owner: "ghost",
				Name:  "dummy-repo2",
				HitFiles: crawler.Files{
					{
						URL:       "https://github.com/ghost/dummy-repo2/dummy3.md",
						Fragments: []string{"detect dummy3-1", "detect dummy3-2"},
					},
				},
			},
		},
		HitCount: 1,
	},
}

func TestFmtTop(t *testing.T) {
	expected := strings.Join([]string{"test1の検索結果: 1件", "test2の検索結果: 1件"}, "\n")

	actual, err := FmtTop(input1)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

func TestFmtDetail(t *testing.T) {
	expected := strings.Join([]string{"test1の詳細結果:ghost/dummy-repo1",
		"-->https://github.com/ghost/dummy-repo1/dummy1.md",
		"-->https://github.com/ghost/dummy-repo1/dummy2.md"}, "\n")

	actual, err := FmtDetail(input1[0])
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

var input2 = []SearchResult{
	{
		Query: "test1",
		Repos: crawler.Repositories{
			{
				URL:   "https://github.com/ghost/dummy-repo1",
				Owner: "ghost",
				Name:  "dummy-repo1",
				HitFiles: crawler.Files{
					{
						URL:       "https://github.com/ghost/dummy-repo1/dummy1.md",
						Fragments: []string{"detect dummy1-1", "detect dummy1-2"},
					},
					{
						URL:       "https://github.com/ghost/dummy-repo1/dummy2.md",
						Fragments: []string{"detect dummy2-1", "detect dummy2-2"},
					},
					{
						URL:       "https://github.com/ghost/dummy-repo1/dummy3.md",
						Fragments: []string{"detect dummy3-1", "detect dummy3-2"},
					},
					{
						URL:       "https://github.com/ghost/dummy-repo1/dummy4.md",
						Fragments: []string{"detect dummy4-1", "detect dummy4-2"},
					},
				},
			},
		},
		HitCount: 1,
	},
}

func TestFmtDetailOmmit(t *testing.T) {
	expected := strings.Join([]string{"test1の詳細結果:ghost/dummy-repo1",
		"-->https://github.com/ghost/dummy-repo1/dummy1.md",
		"-->https://github.com/ghost/dummy-repo1/dummy2.md",
		"-->https://github.com/ghost/dummy-repo1/dummy3.md",
		"-->..."}, "\n")

	actual, err := FmtDetail(input2[0])
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}
