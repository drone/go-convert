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

package travis

var defaultInstall = map[string]string{
	"csharp":  "nuget restore solution-name.sln",
	"clojure": "lein deps",
	"crystal": "shards install",
	"dart":    "pub get",
	"elixir":  "mix local.rebar --force; mix local.hex --force; mix deps.get",
	"elm":     "npm install",
	"erlang":  "rebar get-deps",
	"go":      "go install",
}

var defaultScript = map[string]string{
	"c":       "./configure && make && make test",
	"csharp":  "msbuild /p:Configuration=Release solution-name.sln",
	"cpp":     "./configure && make && make test",
	"clojure": "lein test",
	"crystal": "crystal spec",
	"d":       "dub test --compiler=$DC",
	"dart":    "pub run test",
	"elixir":  "mix test",
	"elm":     "elm-format --validate . && elm-test",
	"erlang":  "rebar compile && rebar skip_deps=true eunit",
	"go":      "go test -v",
}
