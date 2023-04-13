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

package yaml

// when: << pipeline.parameters.run_integration_tests >>

// when:
//   or:
//   - equal: [ main, << pipeline.git.branch >> ]
//   - equal: [ staging, << pipeline.git.branch >> ]

// when:
//   and:
//   - not:
// 	     matches:
// 		   pattern: "^main$"
// 		   value: << pipeline.git.branch >>
//   - or:
// 	   - equal: [ canary, << pipeline.git.tag >> ]
// 	   - << pipeline.parameters.deploy-canary >>

// when:
//   equal: [ *macos-executor, << parameters.os >> ]
