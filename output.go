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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/itchyny/gojq"
)

func PrintVersion(appVersion string) {
	fmt.Printf(`%s %s(sacloud-usage-lib: %s)
Compiler: %s %s
`,
		os.Args[0],
		appVersion,
		Version,
		runtime.Compiler,
		runtime.Version())
}

func OutputMetrics(w io.Writer, metrics map[string]interface{}, query string) error {
	if query == "" {
		v, _ := json.Marshal(metrics)
		fmt.Fprintln(w, string(v))
		return nil
	}

	parsed, err := gojq.Parse(query)
	if err != nil {
		return err
	}
	iter := parsed.Run(metrics)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return err
		}
		if v == nil {
			return fmt.Errorf("%s not found in result", query)
		}
		j2, _ := json.Marshal(v)
		fmt.Fprintln(w, string(j2))
	}

	return nil
}
