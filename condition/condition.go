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
package condition

import (
	"github.com/kujtimiihoxha/go-brace-expansion"
	"strings"
	"unsafe"
)

type Options struct {
	GitHubToken  string   `json:"github_token"  envconfig:"GITHUB_API_TOKEN"`
	SlackToken   string   `json:"slack_token"   envconfig:"SLACK_API_TOKEN"`
	SlackChannel string   `json:"slack_channel" envconfig:"SLACK_CHANNEL"`
	SearchList   []Search `json:"search_list"`
}

type Search struct {
	QueryList  []Sentence `json:"queries"`
	SkipRepos  string     `json:"skip_repos"`
	SkipLibs   string     `json:"skip_libs"`
	SkipOwners string     `json:"skip_owners"`
}

func (o Options) ExpandSearch() []Search {
	var result []Search
	for _, v := range o.SearchList {
		result = append(result, v.Expand()...)
	}
	return result
}

func (s Search) StringWordList() []string {
	return *(*[]string)(unsafe.Pointer(&s.QueryList))
}

func (s Search) Expand() []Search {
	expand := gobrex.Expand(strings.Join(s.StringWordList(), "\n"))
	if len(expand) == 1 {
		// Not expand result
		return []Search{s}
	}

	var result []Search
	for _, v := range expand {
		result = append(result, Search{
			QueryList:  Sentences(strings.Split(v, "\n")),
			SkipOwners: s.SkipOwners,
			SkipRepos:  s.SkipRepos,
			SkipLibs:   s.SkipLibs,
		})
	}

	return result
}

type Sentence string

func (s Sentence) Parse() []string {
	var result []string
	split := strings.Split(string(s), "+")
	for _, v := range split {
		if strings.Contains(v, ":") {
			// Skip GitHub SearchList syntax("type:code" OR "extension:rb" OR ...)
			continue
		}
		spaceSplit := strings.Split(v, " ")
		result = append(result, spaceSplit...)
	}
	return result
}

func Sentences(arr []string) []Sentence {
	var res []Sentence
	for _, v := range arr {
		res = append(res, Sentence(v))
	}
	return res
}

func (o *Options) Override(overOptions Options) Options {
	result := Options{
		GitHubToken:  o.GitHubToken,
		SlackToken:   o.SlackToken,
		SlackChannel: o.SlackChannel,
	}

	if overOptions.GitHubToken != "" {
		result.GitHubToken = overOptions.GitHubToken
	}
	if len(overOptions.SearchList) != 0 {
		result.SearchList = overOptions.SearchList
	}
	if overOptions.SlackToken != "" {
		result.SlackToken = overOptions.SlackToken
	}
	if overOptions.SlackChannel != "" {
		result.SlackChannel = overOptions.SlackChannel
	}
	return result
}
