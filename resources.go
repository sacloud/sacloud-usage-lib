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
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/sacloud/iaas-api-go/types"
)

type MonitorValue struct {
	Time  time.Time
	Value float64
}

type Resource struct {
	ID   types.ID
	Name string
	Zone string

	Monitors []MonitorValue
	Label    string

	AdditionalInfo map[string]interface{}
}

func (r *Resource) toMetrics() map[string]interface{} {
	sum := float64(0)
	monitors := make([]interface{}, 0)
	for _, p := range r.Monitors {
		m := map[string]interface{}{
			r.Label: p.Value,
			"time":  p.Time.String(),
		}
		monitors = append(monitors, m)
		sum += p.Value

		info := strings.Builder{}
		for k, v := range r.AdditionalInfo {
			info.WriteString(fmt.Sprintf("%s:%v ", k, v))
		}
		strInfo := strings.TrimSpace(info.String())

		log.Printf("%s zone:%s %s %s:%f time:%s", r.Name, r.Zone, strInfo, r.Label, p.Value, p.Time.String())
	}

	avg := sum / float64(len(r.Monitors))
	log.Printf("%s average_%s:%f", r.Name, r.Label, avg)

	metrics := map[string]interface{}{
		"name":     r.Name,
		"zone":     r.Zone,
		"avg":      avg,
		"monitors": monitors,
	}

	for k, v := range r.AdditionalInfo {
		metrics[k] = v
	}

	return metrics
}

type Resources struct {
	Resources []*Resource
	Label     string
	Option    *Option
}

func (rs *Resources) Metrics() map[string]interface{} {
	var fs sort.Float64Slice
	resources := make([]interface{}, 0)
	total := float64(0)
	label := rs.Label
	if label == "" {
		label = "resources"
	}

	for _, t := range rs.Resources {
		metrics := t.toMetrics()
		avg := metrics["avg"].(float64)

		fs = append(fs, avg)
		total += avg

		resources = append(resources, metrics)
	}

	if len(fs) == 0 {
		result := map[string]interface{}{}
		result["max"] = float64(0)
		result["avg"] = float64(0)
		result["min"] = float64(0)
		if rs.Option != nil {
			for _, p := range rs.Option.percentiles {
				result[fmt.Sprintf("%spt", p.str)] = float64(0)
			}
		}
		result[label] = resources
		return result
	}

	sort.Sort(fs)
	fl := float64(len(fs))
	result := map[string]interface{}{}
	result["max"] = fs[len(fs)-1]
	result["avg"] = total / fl
	result["min"] = fs[0]
	if rs.Option != nil {
		for _, p := range rs.Option.percentiles {
			result[fmt.Sprintf("%spt", p.str)] = fs[round(fl*(p.float))]
		}
	}
	result[label] = resources
	return result
}

func round(f float64) int64 {
	return int64(math.Round(f)) - 1
}
