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
package diaper

import (
	"context"
	"errors"
	"github.com/future-architect/code-diaper/condition"
	"github.com/future-architect/code-diaper/crawler"
	"github.com/future-architect/code-diaper/filter"
	"github.com/future-architect/code-diaper/formatter"
	"strings"
)

func Run(ctx context.Context, ops condition.Options) (*Message, error) {

	searchList := ops.ExpandSearch()

	var resultList []formatter.SearchResult
	for _, search := range searchList {
		detect, err := RunSearch(ctx, ops.GitHubToken, search)
		if err != nil {
			return nil, err
		}
		resultList = append(resultList, formatter.NewSearchResult(strings.Join(search.StringWordList(), "&"), detect))
	}

	if len(resultList) == 0 {
		msg := ""
		for _, v := range searchList {
			if msg != "" {
				msg += ","
			}
			msg += "「" + strings.Join(v.StringWordList(), "&") + "」"
		}
		return &Message{
			Summary: "GitHub Search Result is 0. Query:" + msg,
			Details: nil,
		}, nil
	}

	summary, err := formatter.FmtTop(resultList)
	if err != nil {
		return nil, err
	}

	var details []string
	for _, v := range resultList {
		if len(v.Repos) == 0 {
			continue
		}
		detail, err := formatter.FmtDetail(v)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}

	return &Message{
		Summary: summary,
		Details: details,
	}, nil
}

func RunSearch(ctx context.Context, githubToken string, s condition.Search) (crawler.Repositories, error) {

	if githubToken == "" {
		return nil, errors.New("required parameter: GitHubToken")
	}
	if len(s.QueryList) == 0 {
		return nil, errors.New("required parameter: SearchWord must be at least one")
	}

	repos := strings.Split(s.SkipRepos, ",")
	libs := strings.Split(s.SkipLibs, ",")
	owners := strings.Split(s.SkipOwners, ",")
	skipFilter := filter.NewSkipFilter(s.QueryList, repos, libs, owners)

	gc := crawler.NewGitHubCrawler(githubToken)

	originalResult, err := gc.Search(ctx, s.StringWordList())
	if err != nil {
		return nil, err
	}

	// if repository that has skip name is forked and renamed then it is too skipped.
	return gc.FulfillForkSource(ctx, skipFilter.Do(originalResult))
}
