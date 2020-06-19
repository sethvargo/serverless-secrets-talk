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
	"fmt"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
)

func berglasAccess(obj string) (string, error) {
	ctx := context.Background()

	projectID, err := valueFromMetadata(ctx, "project/project-id")
	if err != nil {
		return "", fmt.Errorf("failed to get project: %w", err)
	}

	resp, err := berglas.Access(ctx, &berglas.AccessRequest{
		Bucket: fmt.Sprintf("%s-secrets", projectID),
		Object: obj,
	})
	if err != nil {
		return "", err
	}

	return string(resp), nil
}
