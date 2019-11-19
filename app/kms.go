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
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"github.com/pkg/errors"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func kmsDecrypt(input string) (string, error) {
	ctx := context.Background()

	projectID, err := valueFromMetadata(ctx, "project/project-id")
	if err != nil {
		return "", errors.Wrap(err, "failed to get project ID")
	}

	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", err
	}

	b, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	resp, err := client.Decrypt(ctx, &kmspb.DecryptRequest{
		Name:       fmt.Sprintf("projects/%s/locations/global/keyRings/serverless/cryptoKeys/secrets", projectID),
		Ciphertext: b,
	})
	if err != nil {
		return "", err
	}

	return string(resp.GetPlaintext()), nil
}
