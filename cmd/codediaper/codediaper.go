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
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/future-architect/code-diaper/condition"
	"github.com/future-architect/code-diaper/diaper"
	"github.com/future-architect/code-diaper/reporter"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type SearchSentenceArgs []condition.Sentence

func (s *SearchSentenceArgs) String() string {
	return "search word list representation"
}

func (s *SearchSentenceArgs) Set(value string) error {
	*s = append(*s, condition.Sentence(value))
	return nil
}

func main() {
	ctx := context.Background()

	var envOps condition.Options
	if err := envconfig.Process("", &envOps); err != nil {
		panic(err)
	}

	fs := flag.NewFlagSet(fmt.Sprintf("%s (v%s)", "codediaper", "0.01"), flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var searchSentenceList = SearchSentenceArgs{}
	fs.Var(&searchSentenceList, "searchWord", "SearchList word that represents leak key word")

	var (
		githubToken   = fs.String("githubToken", "", "Github access token")
		skipOwnerList = fs.String("skipOwners", "", "Skip repository owner name list. comma separated. if contained in file path then skipped")
		skipRepoList  = fs.String("skipRepos", "", "Skip repository name list. comma separated. if matched exactly then skipped")
		skipLibList   = fs.String("skipLibs", "", "Skip library name list. comma separated. if contained in file path then skipped")
		slackEnabled  = fs.Bool("slackEnabled", false, "Slack notification enabled. default false")
		slackToken    = fs.String("slackToken", "", "Slack access token")
		slackChannel  = fs.String("slackChannel", "", "Slack channel ID")
	)

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	cliOps := condition.Options{
		GitHubToken: *githubToken,
		SearchList: []condition.Search{
			{
				QueryList:  searchSentenceList,
				SkipOwners: *skipOwnerList,
				SkipRepos:  *skipRepoList,
				SkipLibs:   *skipLibList,
			},
		},
		SlackToken:   *slackToken,
		SlackChannel: *slackChannel,
	}

	ops := envOps.Override(cliOps)
	message, err := diaper.Run(ctx, ops)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message.Summary)
	for _, v := range message.Details {
		fmt.Println(v)
	}

	if *slackEnabled {
		slack := reporter.NewSlackReporter(ops.SlackToken, ops.SlackChannel)

		ts, err := slack.Post(ctx, message.Summary)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range message.Details {
			if err := slack.PostThread(ctx, ts, v); err != nil {
				log.Fatal(err)
			}
		}
	}
}
