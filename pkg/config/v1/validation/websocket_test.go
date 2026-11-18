// Copyright 2026 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"testing"

	"github.com/stretchr/testify/require"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func TestValidateWebsocketPath(t *testing.T) {
	for _, tc := range []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "omitted"},
		{name: "root", path: "/"},
		{name: "custom", path: "/relay/egress"},
		{name: "trailing slash", path: "/relay/egress/"},
		{name: "escaped path", path: "/relay/hello%20world"},
		{name: "relative", path: "relay/egress", wantErr: true},
		{name: "full URL", path: "wss://example.com/relay", wantErr: true},
		{name: "authority", path: "//example.com/relay", wantErr: true},
		{name: "query", path: "/relay?token=value", wantErr: true},
		{name: "empty query", path: "/relay?", wantErr: true},
		{name: "fragment", path: "/relay#fragment", wantErr: true},
		{name: "space", path: "/relay path", wantErr: true},
		{name: "newline", path: "/relay\r\nHeader: value", wantErr: true},
		{name: "invalid escape", path: "/relay%zz", wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := validClientConfigWithAuth(v1.AuthClientConfig{Method: v1.AuthMethodToken})
			client.Transport.Protocol = "websocket"
			client.Transport.WSURL = tc.path
			require.NoError(t, client.Complete())
			server := validServerConfigWithAuth(v1.AuthServerConfig{Method: v1.AuthMethodToken})
			server.Transport.WSURL = tc.path
			require.NoError(t, server.Complete())

			validator := NewConfigValidator(nil)
			_, clientErr := validator.ValidateClientCommonConfig(client)
			_, serverErr := validator.ValidateServerConfig(server)
			if tc.wantErr {
				require.ErrorContains(t, clientErr, "invalid transport.ws_url")
				require.ErrorContains(t, serverErr, "invalid transport.ws_url")
				return
			}
			require.NoError(t, clientErr)
			require.NoError(t, serverErr)
		})
	}
}
