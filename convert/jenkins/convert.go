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
	"time"

	"github.com/drone/go-convert/convert/drone"
	"github.com/drone/go-convert/convert/github"
	"github.com/drone/go-convert/convert/gitlab"
)

// Converter converts a Drone pipeline to a Harness
// v1 pipeline.
type Converter struct {
	format        Format
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	debug         bool
	token         string
	attempts      int
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

	// set the minimum number of attempts
	if d.attempts == 0 {
		d.attempts = 1
	}

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return d.ConvertBytes(b)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertBytes(b []byte) ([]byte, error) {
	return d.retry(b)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertString(s string) ([]byte, error) {
	return d.ConvertBytes([]byte(s))
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

// retry attempts the conversion with a backoff
func (d *Converter) retry(src []byte) ([]byte, error) {
	var out []byte
	var err error
	for i := 0; i < d.attempts; i++ {
		// puase before retry
		if i != 0 {
			// print status for debug purposes
			fmt.Fprintln(os.Stderr, "attempt failed")
			fmt.Fprintln(os.Stderr, err)
			// 10 seconds before retry
			time.Sleep(time.Second * 10)
		}
		// attempt the conversion
		if out, err = d.convert(src); err == nil {
			break
		}
	}
	return out, err
}

// convert converts a Drone pipeline to a Harness pipeline.
func (d *Converter) convert(src []byte) ([]byte, error) {

	// gpt input
	req := &request{
		Model: "gpt-3.5-turbo",
		Messages: []*message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Convert this Jenkinsfile to a %s Yaml.\n\n```\n%s\n```\n", d.format.String(), []byte(src)),
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

	if d.format == FromDrone {
		// convert the pipeline yaml from the drone
		// format to the harness yaml format.
		converter := drone.New(
			drone.WithDockerhub(d.dockerhubConn),
			drone.WithKubernetes(d.kubeConnector, d.kubeNamespace),
		)
		pipeline, err := converter.ConvertString(code)
		if err != nil && d.debug {
			// dump data for debug mode
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString("---")
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString(res.Choices[0].Message.Content)
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString("---")
			os.Stdout.WriteString("\n")
			os.Stdout.Write(pipeline)
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString("---")
			os.Stdout.WriteString("\n")
		}
		return pipeline, err
	}

	if d.format == FromGitlab {
		// convert the pipeline yaml from the gitlab
		// format to the harness yaml format.
		converter := gitlab.New(
			gitlab.WithDockerhub(d.dockerhubConn),
			gitlab.WithKubernetes(d.kubeConnector, d.kubeNamespace),
		)
		pipeline, err := converter.ConvertString(code)
		if err != nil {
			// dump data for debug mode
			if err != nil && d.debug {
				os.Stdout.WriteString("\n")
				os.Stdout.WriteString("---")
				os.Stdout.WriteString("\n")
				os.Stdout.WriteString(res.Choices[0].Message.Content)
				os.Stdout.WriteString("\n")
				os.Stdout.WriteString("---")
				os.Stdout.WriteString("\n")
				os.Stdout.Write(pipeline)
				os.Stdout.WriteString("\n")
				os.Stdout.WriteString("---")
				os.Stdout.WriteString("\n")
			}
		}
		return pipeline, err
	}

	// convert the pipeline yaml from the github
	// format to the harness yaml format.
	converter := github.New(
		github.WithDockerhub(d.dockerhubConn),
		github.WithKubernetes(d.kubeConnector, d.kubeNamespace),
	)
	pipeline, err := converter.ConvertString(code)
	if err != nil {
		// dump data for debug mode
		if err != nil && d.debug {
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString("---")
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString(res.Choices[0].Message.Content)
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString("---")
			os.Stdout.WriteString("\n")
			os.Stdout.Write(pipeline)
			os.Stdout.WriteString("\n")
			os.Stdout.WriteString("---")
			os.Stdout.WriteString("\n")
		}
	}

	return pipeline, err
}

func extractCodeFence(s string) string {
	// trim space
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "```")
	// find and trim the code fence prefix
	if _, c, ok := strings.Cut(s, "```"); ok {
		s = c
		// find and trim the code fence suffix
		if c, _, ok := strings.Cut(s, "```"); ok {
			s = c
		}
	}
	return strings.TrimPrefix(s, "yaml")
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
