CodeDiaper
====
<img src="https://img.shields.io/badge/go-v1.12-green.svg" />

You can search for a specific string from all the source code on GitHub and check if it has been posted illegally.

## Usage

This package uses below services.

* GitHub API
* Slack API(Optional)
* Google Cloud Functions(Optional)

## Motivation

I want to detect when a developer accidentally submits a confidential code to GitHub or misconfigures the Public setting.
COPYRIGHT is described as a comment of the code in many confidential codes.
This tool aims to detect illegal posts by specifying such strings.
It seems that this can be achieved using the standard GitHub API, 
but it only tells you what is contained somewhere in the file. With this tool, you can more accurately detect suspicious code.


## QuickStart(Command Line)

### Requirements

* [Go](https://golang.org/dl/) more than 1.11

### Steps

1. Get [GitHub API Token](https://github.blog/2013-05-16-personal-api-tokens/)
2. Install
`go get -u github.com/pj-cancan/code-diaper/cmd/codediaper`
3. Run
```sh
codediaper -githubToken <Your GitHub Token> \
  -searchWord="Copyright+{2019,2018,2017}+Future+Corporation" \
  -skipOwners=future-architect \
  -skipRepos=vuls,ap4r,uroborosql \
  -skipLibs=lib/ap4r \
  -slackEnabled=false
```
4. Result
You can see search result. "Copyright 2019 Future Corporation", "Copyright 2018 Future Corporation", etc.

## QuickStart(Google Cloud Functions)

### Requirements

* [Go](https://golang.org/dl/) more than 1.11
* [Cloud SDK](https://cloud.google.com/sdk/install/)

### Steps

1. Get [GitHub API Token](https://github.blog/2013-05-16-personal-api-tokens/)
2. [Get Slack API Token](https://get.slack.help/hc/en-us/articles/215770388-Create-and-regenerate-API-tokens)
3. Set Cloud Scheduler
```sh
# Mac/Linux
gcloud beta scheduler jobs create pubsub code-diaper --project <YOUR GCP PROJECT> \
  --schedule "55 23 * * *" \
  --topic topic-code-diaper \
  --message-body='{"search":[{"word_list":"<YOUR SEARCH WORD>", "skip_owners":<YOUR SKIP OWNER LIST>", skip_repos":"<YOUR SKIP LIST>"}]}' \
  --time-zone "Asia/Tokyo" \
  --description "This job invokes CloudFunction of code-diaper"

# Windows
gcloud beta scheduler jobs create pubsub code-diaper --project <YOUR GCP PROJECT> ^
  --schedule "55 23 * * *" ^
  --topic topic-code-diaper ^
  --message-body="{\"search_list\":[{\"queries\":[\"<YOUR SEARCH WORD>\"], "skip_owners":<YOUR SKIP OWNER LIST>", \"skip_repos\":\"<YOUR SKIP LIST>\"}]}" ^
  --time-zone "Asia/Tokyo" ^
  --description "This job invokes CloudFunction of code-diaper"

```
4. Deploy to Cloud Functions
```sh
gcloud functions deploy codeDiaper --project <YOUR GCP PROJECT> \
  --entry-point Subscribe \
  --trigger-resource topic-code-diaper \
  --trigger-event google.pubsub.topic.publish \
  --timeout=540s \
  --runtime go111 \
  --set-env-vars GITHUB_API_TOKEN=<github-api-token> \
  --set-env-vars SLACK_API_TOKEN=<slack-api-token> \
  --set-env-vars SLACK_CHANNEL=<slack-channel-name>
```
5. Go to the [Cloud Scheduler page](https://cloud.google.com/scheduler/docs/tut-pub-sub) and click the *run now* button of *code-diaper*


## Example

// TODO

## Options

| CLI Arg       | Env              | Notes                                         | Type                | Example          |
|---------------|------------------|-----------------------------------------------|---------------------|------------------|
| githubToken   | GITHUB_API_TOKEN | GitHub Access Token                           | Required            |                  |
| searchWord    | SEARCH_WORDS     | GitHub Search word. Comma separated.          | Required            | apple+orange     |
| skipOwnerList | SKIP_OWNER_LIST  | Skip Owner name list. Comma separated.        | Optional            | future-architect |
| skipRepoList  | SKIP_REPO_LIST   | Skip repository name list. Comma separated.   | Optional            | repo1,repo2      |
| skipLibList   | SKIP_LIB_LIST    | Skip library name list. Comma separated.      | Optional            | lib/emoji        |
| slackEnabled  | ---              | Skip library name list                        | Optional            | true / false     |
| slackToken    | SLACK_API_TOKEN  | Slack Access Token                            | Optional            |                  |
| slackChannel  | SLACK_CHANNEL    | Slack Channel ID                              | Optional            |                  |

Tips:

The GitHub API has a limit on the maximum number of searches for a term. Therefore,
it is necessary to set keywords that will reduce the number of searches as much as possible.

This is a trade-off. If too many keywords are set, there is a risk of missing leaked codes.

If there are many false positives, you can exclude them by adding a skip list.


## Developer Guide

Install git pre-commit hook script before developing.

```bash
# Windows
git clone https://github.com/pj-cancan/code-diaper
copy /Y .\githooks\*.* .\.git\hooks

# Mac/Linux
git clone https://github.com/pj-cancan/code-diaper
cp githooks/* .git/hooks
chmod +x .git/hooks/pre-commit
```

## License

This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
