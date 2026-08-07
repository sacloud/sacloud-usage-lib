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
	"os"
	"runtime"

	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/common/packages/envvar"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

func SacloudAPICaller(productName, version string) (iaas.APICaller, error) {
	accessToken := envvar.StringFromEnvMulti([]string{"SAKURA_ACCESS_TOKEN", "SAKURACLOUD_ACCESS_TOKEN"}, "")
	if accessToken == "" {
		return nil, fmt.Errorf("environment variable %q is required", "SAKURACLOUD_ACCESS_TOKEN")
	}
	accessTokenSecret := envvar.StringFromEnvMulti([]string{"SAKURA_ACCESS_TOKEN_SECRET", "SAKURACLOUD_ACCESS_TOKEN_SECRET"}, "")
	if accessTokenSecret == "" {
		return nil, fmt.Errorf("environment variable %q is required", "SAKURACLOUD_ACCESS_TOKEN_SECRET")
	}

	userAgent := fmt.Sprintf(
		"sacloud/%s/v%s (%s/%s; +https://github.com/sacloud/%s) %s",
		productName,
		version,
		runtime.GOOS,
		runtime.GOARCH,
		productName,
		iaas.DefaultUserAgent,
	)

	var client saclient.Client
	if err := client.SetEnviron(os.Environ()); err != nil {
		return nil, err
	}
	if err := client.SetWith(
		saclient.WithoutProfile(),
		saclient.WithUserAgent(userAgent),
	); err != nil {
		return nil, err
	}
	if err := client.Populate(); err != nil {
		return nil, err
	}

	cfg, err := client.EndpointConfig()
	if err != nil {
		return nil, err
	}
	if cfg.APIRootURL != "" {
		iaas.SakuraCloudAPIRoot = cfg.APIRootURL
	}
	if len(cfg.Zones) > 0 {
		iaas.SakuraCloudZones = cfg.Zones
	}
	if cfg.DefaultZone != "" {
		iaas.APIDefaultZone = cfg.DefaultZone
	}
	return iaas.NewClientFromSaclient(&client), nil
}
