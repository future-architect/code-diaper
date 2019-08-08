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
	"bytes"
	"github.com/future-architect/code-diaper/crawler"
	"strings"
	"text/template"
)

const TopMessage = `
{{ range $i, $sr := . -}}
	{{- $sr.Query}}の検索結果: {{$sr.HitCount}}件
{{ end -}}
`

const DetailMessage = `
{{ range $i, $repo := .Repos -}}
{{ $.Query -}}の詳細結果:{{- $repo.Owner }}/{{- $repo.Name }}{{printf "\n" }}
	{{- range $j, $file := $repo.HitFiles -}}
		{{- if lt $j 3 -}}
-->{{ $file.URL }}{{printf "\n" }}
		{{- else if eq $j 3 -}}
-->...{{printf "\n" }}
		{{- end -}}
	{{- end -}}
{{ end -}}
`

type SearchResult struct {
	Query    string
	Repos    crawler.Repositories
	HitCount int
}

func NewSearchResult(searchWord string, reps crawler.Repositories) SearchResult {
	return SearchResult{
		Query:    searchWord,
		Repos:    reps,
		HitCount: len(reps),
	}
}

func FmtTop(list []SearchResult) (string, error) {
	var buff bytes.Buffer
	topTemplate := template.Must(template.New("top").Parse(TopMessage))
	err := topTemplate.Execute(&buff, list)
	return strings.TrimSpace(buff.String()), err
}

func FmtDetail(sr SearchResult) (string, error) {
	var buff bytes.Buffer
	topTemplate := template.Must(template.New("detail").Parse(DetailMessage))
	err := topTemplate.Execute(&buff, sr)
	return strings.TrimSpace(buff.String()), err
}
