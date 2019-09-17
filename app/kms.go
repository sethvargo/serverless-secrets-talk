// Copyright 2019 Seth Vargo
// Copyright 2019 Google, Inc.
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

package main

import (
	"context"
	"encoding/base64"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func kmsDecrypt(input string) (string, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", err
	}

	b, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	resp, err := client.Decrypt(ctx, &kmspb.DecryptRequest{
		Name:       "projects/sethvargo-devsecconseattle-19/locations/global/keyRings/serverless/cryptoKeys/secrets",
		Ciphertext: b,
	})
	if err != nil {
		return "", err
	}

	return string(resp.GetPlaintext()), nil
}
