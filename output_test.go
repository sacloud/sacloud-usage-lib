// Copyright 2023 The sacloud/sacloud-usage-lib Authors
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

package usage

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOutputMetrics(t *testing.T) {
	log.SetOutput(io.Discard)

	type args struct {
		metrics map[string]interface{}
		query   string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "without query",
			args: args{
				metrics: map[string]interface{}{
					"90pt":    1.,
					"avg":     2.,
					"max":     3.,
					"min":     4.,
					"routers": []interface{}{},
				},
				query: "",
			},
			wantW:   `{"90pt":1,"avg":2,"max":3,"min":4,"routers":[]}`,
			wantErr: false,
		},
		{
			name: "with query",
			args: args{
				metrics: map[string]interface{}{
					"90pt":    1.,
					"avg":     2.,
					"max":     3.,
					"min":     4.,
					"routers": []interface{}{},
				},
				query: ".avg",
			},
			wantW:   `2`,
			wantErr: false,
		},
		{
			name: "invalid query",
			args: args{
				metrics: map[string]interface{}{
					"90pt":    1.,
					"avg":     2.,
					"max":     3.,
					"min":     4.,
					"routers": []interface{}{},
				},
				query: "invalid-query",
			},
			wantW:   ``,
			wantErr: true,
		},
		{
			name: "query returns no value",
			args: args{
				metrics: map[string]interface{}{
					"90pt":    1.,
					"avg":     2.,
					"max":     3.,
					"min":     4.,
					"routers": []interface{}{},
				},
				query: ".not_exists",
			},
			wantW:   ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := OutputMetrics(w, tt.args.metrics, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("outputMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				require.Equal(t, tt.wantW+"\n", w.String())
			}
		})
	}
}
