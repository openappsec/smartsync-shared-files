// Copyright (C) 2022 Check Point Software Technologies Ltd. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"context"
	"net/http"

	"openappsec.io/ctxutils"
	"openappsec.io/errors/errorloader"
	"openappsec.io/log"
)

// CreateErrorBody loads the error body from the location of err-responses.json file
func CreateErrorBody(ctx context.Context, errorName string) string {
	errorResponse, err := errorloader.GetError(ctx, errorName)
	if err != nil {
		log.WithContext(ctx).Errorf("Failed to create error body, using default. Error: %s", err.Error())
		traceID := ctxutils.ExtractString(ctx, ctxutils.ContextKeyEventTraceID)
		errorBody := errorloader.NewErrorResponse(traceID, http.StatusText(http.StatusInternalServerError))
		return (&errorBody).Error()
	}

	return errorResponse.Error()
}
