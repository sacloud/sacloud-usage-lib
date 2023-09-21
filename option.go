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
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

type Option struct {
	Time          uint     `long:"time" description:"Get average traffic for a specified amount of time" default:"3"`
	Prefix        []string `long:"prefix" description:"prefix for router names. prefix accepts more than one." required:"true"`
	Zones         []string `long:"zone" description:"zone name" required:"true"`
	PercentileSet string   `long:"percentile-set" default:"99,95,90,75" description:"percentiles to dispaly"`
	Version       bool     `short:"v" long:"version" description:"Show version"`
	Query         string   `long:"query" description:"jq style query to result and display"`
	EnvFrom       string   `long:"env-from" description:"load environment values from this file"`
	percentiles   []percentile
}

type optionProvider interface {
	option() *Option
}

func (o *Option) option() *Option {
	return o
}

type percentile struct {
	str   string
	float float64
}

func ParseOption(o interface{}) error {
	psr := flags.NewParser(o, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()

	v, ok := o.(optionProvider)
	if !ok {
		return nil
	}
	opts := v.option()

	if opts.Version {
		return nil
	}

	if err != nil {
		return err
	}

	if opts.Time < 1 {
		opts.Time = 1
	}

	if opts.EnvFrom != "" {
		if err := godotenv.Load(opts.EnvFrom); err != nil {
			return err
		}
	}

	m := make(map[string]struct{})
	for _, z := range opts.Zones {
		if _, ok := m[z]; ok {
			return fmt.Errorf("zone %q is duplicated", z)
		}
		m[z] = struct{}{}
	}

	var percentiles []percentile
	percentileStrings := strings.Split(opts.PercentileSet, ",")
	for _, s := range percentileStrings {
		if s == "" {
			continue
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("could not parse --percentile-set: %v", err)
		}
		f /= 100
		percentiles = append(percentiles, percentile{s, f})
	}
	opts.percentiles = percentiles

	return nil
}
