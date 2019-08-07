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
package reporter

import (
	"context"
	"github.com/nlopes/slack"
)

type SlackReporter struct {
	api     *slack.Client
	channel string
}

func NewSlackReporter(token, channel string) *SlackReporter {
	return &SlackReporter{
		api:     slack.New(token),
		channel: channel,
	}
}

// Post is send func for slack.
func (s SlackReporter) Post(ctx context.Context, msg string) (string, error) {
	_, ts, err := s.api.PostMessageContext(ctx, s.channel, slack.MsgOptionText(msg, false))
	return ts, err
}

// PostThread is send func for slack thread. timestamp is parent message timestamp.
func (s SlackReporter) PostThread(ctx context.Context, timeStamp, msg string) error {
	_, _, err := s.api.PostMessageContext(ctx, s.channel, slack.MsgOptionText(msg, false), slack.MsgOptionTS(timeStamp))
	return err
}
