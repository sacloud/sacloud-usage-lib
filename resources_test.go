// Copyright (c) 2023 The sacloud/sacloud-usage-lib Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package usage

import (
	"io"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResources_Metrics(t *testing.T) {
	log.SetOutput(io.Discard)

	type args struct {
		resources *Resources
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "empty",
			args: args{
				resources: &Resources{Label: "routers"},
			},
			want: map[string]interface{}{
				"avg":     0.,
				"max":     0.,
				"min":     0.,
				"routers": []interface{}{},
			},
		},
		{
			name: "single resource - single value",
			args: args{
				resources: &Resources{
					Label: "routers",
					Resources: []*Resource{
						{
							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []MonitorValue{
								{Time: time.Unix(1, 0), Value: 3},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
				},
			},
			want: map[string]interface{}{
				"avg": 3.,
				"max": 3.,
				"min": 3.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  3.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 3.,
								"time":    time.Unix(1, 0).String(),
							},
						},
					},
				},
			},
		},
		{
			name: "with percentiles",
			args: args{
				resources: &Resources{
					Label: "routers",
					Resources: []*Resource{
						{

							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []MonitorValue{
								{Time: time.Unix(1, 0), Value: 3},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
					Option: &Option{percentiles: []percentile{{str: "90", float: 0.9}}},
				},
			},
			want: map[string]interface{}{
				"avg":  3.,
				"max":  3.,
				"min":  3.,
				"90pt": 3.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  3.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 3.,
								"time":    time.Unix(1, 0).String(),
							},
						},
					},
				},
			},
		},
		{
			name: "single resource - multi values",
			args: args{
				resources: &Resources{
					Label: "routers",
					Resources: []*Resource{
						{
							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []MonitorValue{
								{Time: time.Unix(1, 0), Value: 1},
								{Time: time.Unix(2, 0), Value: 2},
								{Time: time.Unix(3, 0), Value: 3},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
				},
			},
			want: map[string]interface{}{
				"avg": 2.,
				"max": 2.,
				"min": 2.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  2.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 1.,
								"time":    time.Unix(1, 0).String(),
							},
							map[string]interface{}{
								"traffic": 2.,
								"time":    time.Unix(2, 0).String(),
							},
							map[string]interface{}{
								"traffic": 3.,
								"time":    time.Unix(3, 0).String(),
							},
						},
					},
				},
			},
		},
		{
			name: "multi resources - single value",
			args: args{
				resources: &Resources{
					Label: "routers",
					Resources: []*Resource{
						{
							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []MonitorValue{
								{Time: time.Unix(3, 0), Value: 2},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
						{
							ID:   2,
							Name: "test2",
							Zone: "is1b",
							Monitors: []MonitorValue{
								{Time: time.Unix(3, 0), Value: 4},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
				},
			},
			want: map[string]interface{}{
				"avg": 3.,
				"max": 4.,
				"min": 2.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  2.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 2.,
								"time":    time.Unix(3, 0).String(),
							},
						},
					},
					map[string]interface{}{
						"name": "test2",
						"zone": "is1b",
						"avg":  4.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 4.,
								"time":    time.Unix(3, 0).String(),
							},
						},
					},
				},
			},
		},
		{
			name: "multi resources - multi values",
			args: args{
				resources: &Resources{
					Label: "routers",
					Resources: []*Resource{
						{
							ID:   1,
							Name: "test1",
							Zone: "is1a",
							Monitors: []MonitorValue{
								{Time: time.Unix(1, 0), Value: 1},
								{Time: time.Unix(2, 0), Value: 2},
								{Time: time.Unix(3, 0), Value: 3},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
						{
							ID:   2,
							Name: "test2",
							Zone: "is1b",
							Monitors: []MonitorValue{
								{Time: time.Unix(4, 0), Value: 4},
								{Time: time.Unix(5, 0), Value: 5},
								{Time: time.Unix(6, 0), Value: 6},
							},
							Label:          "traffic",
							AdditionalInfo: nil,
						},
					},
					Option: &Option{percentiles: []percentile{{str: "90", float: 0.9}}},
				},
			},
			want: map[string]interface{}{
				"avg":  3.5,
				"max":  5.,
				"min":  2.,
				"90pt": 5.,
				"routers": []interface{}{
					map[string]interface{}{
						"name": "test1",
						"zone": "is1a",
						"avg":  2.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 1.,
								"time":    time.Unix(1, 0).String(),
							},
							map[string]interface{}{
								"traffic": 2.,
								"time":    time.Unix(2, 0).String(),
							},
							map[string]interface{}{
								"traffic": 3.,
								"time":    time.Unix(3, 0).String(),
							},
						},
					},
					map[string]interface{}{
						"name": "test2",
						"zone": "is1b",
						"avg":  5.,
						"monitors": []interface{}{
							map[string]interface{}{
								"traffic": 4.,
								"time":    time.Unix(4, 0).String(),
							},
							map[string]interface{}{
								"traffic": 5.,
								"time":    time.Unix(5, 0).String(),
							},
							map[string]interface{}{
								"traffic": 6.,
								"time":    time.Unix(6, 0).String(),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.resources.Metrics()
			require.EqualValues(t, tt.want, got)
		})
	}
}
