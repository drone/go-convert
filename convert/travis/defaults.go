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
	"groovy":  "gradle assemble",
	"haskell": "cabal install --only-dependencies --enable-tests",
	"java":    "mvn install -DskipTests=true -Dmaven.javadoc.skip=true -B -V",
	"node_js": "npm install",
	"perl":    "cpanm --quiet --installdeps --notest .",
	"python":  "pip install -r requirements.txt",
	"ruby":    "bundle install --jobs=3 --retry=3",
	"rust":    "cargo build --verbose",

	// default to ruby
	"": "bundle install --jobs=3 --retry=3",
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
	"groovy":  "gradle check",
	"haskell": "cabal configure --enable-tests && cabal build && cabal test",
	"java":    "mvn test -B",
	"node_js": "npm test",
	"matlab":  `matlab -batch "results = runtests('IncludeSubfolders',true); assertSuccess(results);"`,
	"perl":    "perl Build.PL && ./Build test",
	"perl6":   "PERL6LIB=lib prove --ext .t --ext .t6 -v -r --exec=perl6 t/",
	"php":     "phpunit",
	"ruby":    "rake",
	"rust":    "cargo build --verbose; cargo test --verbose",

	// default to ruby
	"": "rake",
}

var defaultImage = map[string]string{
	"c":  "gcc",
	"go": "golang",
}

var defaultServiceImage = map[string]string{
	"mysql":         "mysql",
	"postgresql":    "postgres",
	"mariadb":       "mariadb",
	"mongodb":       "mongo",
	"couchdb":       "couchdb",
	"rabbitmq":      "rabbitmq",
	"redis":         "redis",
	"memcached":     "memcached",
	"elasticsearch": "elasticsearch",
	"neo4j":         "neo4j",
	"riak":          "basho/riak-kv",
	"rethinkdb":     "rethinkdb",
	"cassandra":     "cassandra",
}

var defaultServicePorts = map[string][]string{
	"mysql":         {"3306"},
	"postgresql":    {"5432"},
	"mongodb":       {"27017"},
	"couchdb":       {"5984"},
	"redis":         {"6379"},
	"rabbitmq":      {"5672"},
	"memcached":     {"11211"},
	"elasticsearch": {"9200"},
	"neo4j":         {"7687"},
	"riak":          {"8087", "8098"},
	"rethinkdb":     {"28015"},
	"mariadb":       {"3306"},
	"cassandra":     {"7000"},
}

var defaultServiceEnvs = map[string]map[string]string{
	"postgresql": {
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
	},
	"mariadb": {
		"MARIADB_DATABASE":                  "test",
		"MARIADB_USER":                      "root",
		"MARIADB_ALLOW_EMPTY_ROOT_PASSWORD": "yes",
		"MYSQL_USER":                        "root",
		"MYSQL_ALLOW_EMPTY_PASSWORD":        "yes",
		"MYSQL_DATABASE":                    "test",
	},
	"mysql": {
		"MYSQL_USER":                 "root",
		"MYSQL_ALLOW_EMPTY_PASSWORD": "yes",
		"MYSQL_DATABASE":             "test",
	},
	"couchdb": {
		"COUCHDB_USER":     "admin",
		"COUCHDB_PASSWORD": "travis",
	},
	"rabbitmq": {
		"RABBITMQ_DEFAULT_VHOST": "/",
		"RABBITMQ_DEFAULT_USER":  "guest",
		"RABBITMQ_DEFAULT_PASS":  "guest",
	},
}
