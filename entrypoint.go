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
package function

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"github.com/future-architect/code-diaper/condition"
	"github.com/future-architect/code-diaper/diaper"
	"github.com/future-architect/code-diaper/reporter"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/context"
	"log"
)

// CloudFunction entry point
func Subscribe(ctx context.Context, msg *pubsub.Message) error {
	log.Println("start")

	var envOps condition.Options
	if err := envconfig.Process("", &envOps); err != nil {
		return err
	}

	var msgOps condition.Options
	if err := json.Unmarshal(msg.Data, &msgOps); err != nil {
		return err
	}

	ops := envOps.Override(msgOps)

	message, err := diaper.Run(ctx, ops)
	if err != nil {
		return err
	}

	slack := reporter.NewSlackReporter(ops.SlackToken, ops.SlackChannel)

	ts, err := slack.Post(ctx, message.Summary)
	if err != nil {
		return err
	}

	for _, v := range message.Details {
		if err := slack.PostThread(ctx, ts, v); err != nil {
			return err
		}
	}

	log.Println("finish")
	return nil
}
