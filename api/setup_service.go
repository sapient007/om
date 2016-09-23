package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
)

type SetupInput struct {
	IdentityProvider                 string
	AdminUserName                    string
	AdminPassword                    string
	AdminPasswordConfirmation        string
	DecryptionPassphrase             string
	DecryptionPassphraseConfirmation string
	EULAAccepted                     bool
}

type SetupOutput struct{}

type SetupService struct {
	client httpClient
}

func NewSetupService(client httpClient) SetupService {
	return SetupService{
		client: client,
	}
}

func (ss SetupService) Setup(input SetupInput) (SetupOutput, error) {
	var setup struct {
		Setup struct {
			IdentityProvider                 string `json:"identity_provider"`
			AdminUserName                    string `json:"admin_user_name"`
			AdminPassword                    string `json:"admin_password"`
			AdminPasswordConfirmation        string `json:"admin_password_confirmation"`
			DecryptionPassphrase             string `json:"decryption_passphrase"`
			DecryptionPassphraseConfirmation string `json:"decryption_passphrase_confirmation"`
			EULAAccepted                     string `json:"eula_accepted"`
		} `json:"setup"`
	}

	setup.Setup.IdentityProvider = input.IdentityProvider
	setup.Setup.AdminUserName = input.AdminUserName
	setup.Setup.AdminPassword = input.AdminPassword
	setup.Setup.AdminPasswordConfirmation = input.AdminPasswordConfirmation
	setup.Setup.DecryptionPassphrase = input.DecryptionPassphrase
	setup.Setup.DecryptionPassphraseConfirmation = input.DecryptionPassphraseConfirmation
	setup.Setup.EULAAccepted = strconv.FormatBool(input.EULAAccepted)

	payload, err := json.Marshal(setup)
	if err != nil {
		return SetupOutput{}, err
	}

	request, err := http.NewRequest("POST", "/api/v0/setup", bytes.NewReader(payload))
	if err != nil {
		return SetupOutput{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := ss.client.Do(request)
	if err != nil {
		return SetupOutput{}, fmt.Errorf("could not make api request to setup endpoint: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		out, err := httputil.DumpResponse(response, true)
		if err != nil {
			return SetupOutput{}, fmt.Errorf("request failed: unexpected response: %s", err)
		}

		return SetupOutput{}, fmt.Errorf("request failed: unexpected response:\n%s", out)
	}

	return SetupOutput{}, nil
}
