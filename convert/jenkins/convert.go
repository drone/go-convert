// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jenkins converts Jenkins pipelines to Harness pipelines.
package jenkins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/drone/go-convert/convert/gitlab"
)

// Converter converts a Drone pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	token         string
}

// New creates a new Converter that converts a Drone
// pipeline to a Harness v1 pipeline.
func New(options ...Option) *Converter {
	d := new(Converter)

	// loop through and apply the options.
	for _, option := range options {
		option(d)
	}

	// set the default kubernetes namespace.
	if d.kubeNamespace == "" {
		d.kubeNamespace = "default"
	}

	// set the runtime to kubernetes if the kubernetes
	// connector is configured.
	if d.kubeConnector != "" {
		d.kubeEnabled = true
	}

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return d.convert(b)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertBytes(b []byte) ([]byte, error) {
	return d.convert(b)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertString(s string) ([]byte, error) {
	return d.convert([]byte(s))
}

// ConvertFile downgrades a v1 pipeline.
func (d *Converter) ConvertFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return d.Convert(f)
}

// converts converts a Drone pipeline to a Harness pipeline.
func (d *Converter) convert(src []byte) ([]byte, error) {

	// gpt input
	req := &request{
		Model: "gpt-3.5-turbo",
		Messages: []*message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Convert this Jenkinsfile to a GitLab Yaml. Omit git clone or git checkout steps.\n\n```\n%s\n```\n", []byte(src)),
			},
		},
	}

	// gpt output
	res := new(response)

	// marshal the input to json
	err := d.do("https://api.openai.com/v1/chat/completions", "POST", req, res)
	if err != nil {
		return nil, err
	}

	if len(res.Choices) == 0 {
		return nil, errors.New("chat gpt returned a response with zero choices. conversion not possible.")
	}

	// extract the message
	code := extractCodeFence(res.Choices[0].Message.Content)

	// convert the pipeline yaml from the gitlab
	// format to the harness yaml format.
	converter := gitlab.New(
		gitlab.WithDockerhub(d.dockerhubConn),
		gitlab.WithKubernetes(d.kubeConnector, d.kubeNamespace),
	)
	pipeline, err := converter.ConvertString(code)
	if err != nil {
		// temporarily dump chat gpt response on error
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("\n---")
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString(code)
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("\n---")
		os.Stderr.WriteString("\n")
	}

	return pipeline, err
}

func extractCodeFence(s string) string {
	_, a, _ := strings.Cut(s, "```")
	b, _, _ := strings.Cut(a, "```")
	b = strings.TrimPrefix(b, "yaml")
	return b
}

//
// Chat GPT Client
// TODO move to separate package
//

type request struct {
	Model    string     `json:"model"`
	Messages []*message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	}
}

// helper function to make an http request
func (d *Converter) do(rawurl, method string, in, out interface{}) error {
	body, err := d.open(rawurl, method, in, out)
	if err != nil {
		return err
	}
	defer body.Close()
	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}
	return nil
}

// helper function to open an http request
func (d *Converter) open(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	if in != nil {
		decoded, derr := json.Marshal(in)
		if derr != nil {
			return nil, derr
		}
		buf := bytes.NewBuffer(decoded)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(len(decoded))
		req.Header.Set("Content-Length", strconv.Itoa(len(decoded)))
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+d.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("client error %d: %s", resp.StatusCode, string(out))
	}
	return resp.Body, nil
}
