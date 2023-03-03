// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package yaml

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestBytesSize(t *testing.T) {
	tests := []struct {
		yaml string
		size int64
	}{
		{
			yaml: "1KiB",
			size: 1024,
		},
		{
			yaml: "100MiB",
			size: 104857600,
		},
		{
			yaml: "100mb",
			size: 104857600,
		},
		{
			yaml: "1024",
			size: 1024,
		},
	}
	for _, test := range tests {
		in := []byte(test.yaml)
		out := BytesSize(0)
		err := yaml.Unmarshal(in, &out)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := int64(out), test.size; got != want {
			t.Errorf("Want byte size %d, got %d", want, got)
		}
	}
}
