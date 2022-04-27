/*
Copyright 2021 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"encoding/json"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"

	"github.com/deckhouse/deckhouse/go_lib/pwgen"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue:        "/modules/openvpn/generate_password",
	OnBeforeHelm: &go_hook.OrderedConfig{Order: 10},
}, generatePassword)

func generatePassword(input *go_hook.HookInput) error {
	if input.Values.Exists("openvpn.auth.externalAuthentication") {
		input.ConfigValues.Remove("openvpn.auth.password")
		if input.ConfigValues.Exists("openvpn.auth") && len(input.ConfigValues.Get("openvpn.auth").Map()) == 0 {
			input.ConfigValues.Remove("openvpn.auth")
		}

		return nil
	}

	if input.Values.Exists("openvpn.auth.password") {
		return nil
	}

	if !input.ConfigValues.Exists("openvpn.auth") {
		input.ConfigValues.Set("openvpn.auth", json.RawMessage("{}"))
	}

	generatedPass := pwgen.AlphaNum(20)

	input.ConfigValues.Set("openvpn.auth.password", generatedPass)

	return nil
}
