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

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"reflect"
	"strings"
	"time"
)

const MaxPageSize = 100

type Crawler interface {
	Search(ctx context.Context, words []string) (Repositories, error)
	FulfillForkSource(ctx context.Context, repos Repositories) (Repositories, error)
}

type gitHubCrawler struct {
	client *github.Client
	option *github.SearchOptions
}

func NewGitHubCrawler(token string) Crawler {

	tokenClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	}))

	return &gitHubCrawler{
		client: github.NewClient(tokenClient),
		option: &github.SearchOptions{
			Sort:  "updated",
			Order: "desc",
			ListOptions: github.ListOptions{
				PerPage: MaxPageSize,
			},
			TextMatch: true,
		},
	}
}

// According to the API reference, up to 100 can be obtained with one API call
// https://developer.github.com/v3/search/
// https://developer.github.com/v3/search/#constructing-a-search-query
func (c *gitHubCrawler) Search(ctx context.Context, words []string) (Repositories, error) {

	result := Repositories{}

	apiCallCnt := 0
	for {
		q := strings.Join(words, "+")
		codeSearchResult, resp, err := c.client.Search.Code(ctx, q+"+in:file", c.option)
		if abuseRateLimitErr, ok := err.(*github.AbuseRateLimitError); ok {

			// The time at which the current rate limit window resets in UTC epoch seconds.
			// https://developer.github.com/v3/#rate-limiting
			fmt.Printf("retry after %v\n", *abuseRateLimitErr.RetryAfter)
			time.Sleep(*abuseRateLimitErr.RetryAfter)
			continue

		} else if _, ok := err.(*github.RateLimitError); ok {
			fmt.Printf("RateLimit Exceed\n")
			break
		} else if err != nil {
			fmt.Printf("Something happend: %+v \n type: %+v\n", err.Error(), reflect.TypeOf(err))
			return nil, err
		}

		if apiCallCnt == 0 {
			fmt.Printf("%+v Hits. Continue searching\n", *codeSearchResult.Total)
		}
		apiCallCnt++

		for _, cr := range codeSearchResult.CodeResults {
			files := make(Files, 0, len(cr.TextMatches))
			for _, match := range cr.TextMatches {
				f := File{
					Fragments: []string{match.GetFragment()},
					URL:       cr.GetHTMLURL(),
				}
				files = files.Merge(f)
			}

			r := Repository{
				URL:      *cr.Repository.HTMLURL,
				Owner:    strings.Split(*cr.Repository.FullName, "/")[0],
				Name:     strings.Split(*cr.Repository.FullName, "/")[1],
				HitFiles: files,
			}
			result = result.Merge(r)
		}

		if resp.NextPage == 0 {
			// finish
			break
		}

		// update search condition
		c.option.Page = resp.NextPage

		time.Sleep(1000 * time.Millisecond)
	}

	return result, nil
}

func (c *gitHubCrawler) FulfillForkSource(ctx context.Context, repos Repositories) (Repositories, error) {

	var result Repositories
	for _, v := range repos {
		source, err := c.fetchForkSource(ctx, v.Owner, v.Name)
		if err != nil {
			return nil, err
		}
		v.ForkSource = source

		result = append(result, v)
	}
	return result, nil
}

func (c *gitHubCrawler) fetchForkSource(ctx context.Context, owner, repoName string) (string, error) {
	for {
		repo, _, err := c.client.Repositories.Get(ctx, owner, repoName)

		if abuseRateLimitErr, ok := err.(*github.AbuseRateLimitError); ok {
			// The time at which the current rate limit window resets in UTC epoch seconds.
			// https://developer.github.com/v3/#rate-limiting
			time.Sleep(*abuseRateLimitErr.RetryAfter)
			continue
		} else if _, ok := err.(*github.RateLimitError); ok {
			fmt.Printf("RateLimit Exceed\n")
			break
		} else if err != nil {
			fmt.Printf("something happend: %+v \n type: %+v\n", err.Error(), reflect.TypeOf(err))
			return "", err
		}
		if repo.Source == nil {
			return "", nil
		}

		return *repo.Source.FullName, nil
	}

	return "", nil
}
