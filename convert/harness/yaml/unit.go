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

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/go-units"
)

// BytesSize stores a human-readable size in bytes,
// kibibytes, mebibytes, gibibytes, or tebibytes
// (eg. "44kiB", "17MiB").
type BytesSize int64

// UnmarshalYAML implements yaml unmarshalling.
func (b *BytesSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intType int64
	if err := unmarshal(&intType); err == nil {
		*b = BytesSize(intType)
		return nil
	}

	var stringType string
	if err := unmarshal(&stringType); err != nil {
		return err
	}

	intType, err := units.RAMInBytes(stringType)
	if err == nil {
		*b = BytesSize(intType)
	}
	return err
}

// MarshalJSON makes UnitBytes implement json.Marshaler
func (b *BytesSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON implements json unmarshalling.
func (b *BytesSize) UnmarshalJSON(data []byte) error {
	var intType int64
	if err := json.Unmarshal(data, &intType); err == nil {
		*b = BytesSize(intType)
		return nil
	}

	var stringType string
	if err := json.Unmarshal(data, &stringType); err != nil {
		return err
	}

	intType, err := units.RAMInBytes(stringType)
	if err == nil {
		*b = BytesSize(intType)
	}
	return err
}

// String returns a human-readable size in bytes,
// kibibytes, mebibytes, gibibytes, or tebibytes
// (eg. "44kiB", "17MiB").
func (b BytesSize) String() string {
	return units.BytesSize(float64(b))
}

// MilliSize will convert cpus to millicpus as int64.
// for instance "1" will be converted to 1000 and "100m" to 100
type MilliSize int64

// UnmarshalYAML implements yaml unmarshalling.
func (m *MilliSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var intType int64
	if err := unmarshal(&intType); err == nil {
		*m = MilliSize(intType * 1000)
		return nil
	}

	var stringType string
	if err := unmarshal(&stringType); err != nil {
		return err
	}
	if strings.HasSuffix(stringType, "m") {
		i, err := strconv.ParseInt(strings.TrimSuffix(stringType, "m"), 10, 64)
		if err != nil {
			return err
		}
		*m = MilliSize(i)
		return nil
	}
	return fmt.Errorf("cannot unmarshal cpu millis")
}

// UnmarshalJSON implements json unmarshalling.
func (m *MilliSize) UnmarshalJSON(data []byte) error {
	var intType int64
	if err := json.Unmarshal(data, &intType); err == nil {
		*m = MilliSize(intType * 1000)
		return nil
	}

	var stringType string
	if err := json.Unmarshal(data, &stringType); err != nil {
		return err
	}
	if strings.HasSuffix(stringType, "m") {
		i, err := strconv.ParseInt(strings.TrimSuffix(stringType, "m"), 10, 64)
		if err != nil {
			return err
		}
		*m = MilliSize(i)
		return nil
	}
	return fmt.Errorf("cannot unmarshal %s into cpu millis", string(data))
}

// MarshalJSON makes implements json.Marshaler
func (m *MilliSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// String returns a human-readable cpu millis,
// (eg. "1000m", "10m").
func (m MilliSize) String() string {
	if m == 0 {
		return "0m"
	} else {
		return strconv.FormatInt(int64(m), 10) + "m"
	}
}

//
//
//

// Duration is a wrapper around time.Duration which supports correct
// marshaling to YAML and JSON. In particular, it marshals into strings, which
// can be used as map keys in json.
type Duration struct {
	time.Duration
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	// First try the standard library parser
	if pd, err := time.ParseDuration(str); err == nil {
		d.Duration = pd
		return nil
	}

	// Fallback: support shorthand units not handled by time.ParseDuration, e.g., days (d) and weeks (w)
	lower := strings.ToLower(strings.TrimSpace(str))
	if strings.HasSuffix(lower, "d") {
		val := strings.TrimSuffix(lower, "d")
		if n, convErr := strconv.ParseFloat(val, 64); convErr == nil {
			hours := n * 24
			if pd, perr := time.ParseDuration(fmt.Sprintf("%fh", hours)); perr == nil {
				d.Duration = pd
				return nil
			}
		}
	}
	if strings.HasSuffix(lower, "w") {
		val := strings.TrimSuffix(lower, "w")
		if n, convErr := strconv.ParseFloat(val, 64); convErr == nil {
			hours := n * 7 * 24
			if pd, perr := time.ParseDuration(fmt.Sprintf("%fh", hours)); perr == nil {
				d.Duration = pd
				return nil
			}
		}
	}

	return fmt.Errorf("time: unknown unit in duration %q", str)
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	// First try the standard library parser
	if pd, err := time.ParseDuration(str); err == nil {
		d.Duration = pd
		return nil
	}

	// Fallback: support shorthand units not handled by time.ParseDuration, e.g., days (d) and weeks (w)
	// Examples: 1d -> 24h, 2d -> 48h, 1w -> 168h
	// We also support uppercase letters by normalizing to lowercase.
	lower := strings.ToLower(strings.TrimSpace(str))
	if strings.HasSuffix(lower, "d") {
		val := strings.TrimSuffix(lower, "d")
		if n, convErr := strconv.ParseFloat(val, 64); convErr == nil {
			hours := n * 24
			pd, perr := time.ParseDuration(fmt.Sprintf("%fh", hours))
			if perr == nil {
				d.Duration = pd
				return nil
			}
		}
	}
	if strings.HasSuffix(lower, "w") {
		val := strings.TrimSuffix(lower, "w")
		if n, convErr := strconv.ParseFloat(val, 64); convErr == nil {
			hours := n * 7 * 24
			pd, perr := time.ParseDuration(fmt.Sprintf("%fh", hours))
			if perr == nil {
				d.Duration = pd
				return nil
			}
		}
	}

	return fmt.Errorf("time: unknown unit in duration %q", str)
}

// MarshalJSON implements the json.Marshaler interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	if d.Duration == 0 {
		return json.Marshal("")
	}
	return json.Marshal(d.Duration.String())
}
