/*
Copyright 2020 Cortex Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package batchapi

import (
	"fmt"

	"github.com/cortexlabs/cortex/pkg/lib/errors"
)

const (
	ErrJobNotFound     = "batchapi.job_spec_not_found"
	ErrMissingJobState = "batchapi.missing_job_state"
	ErrAPINotDeployed  = "batchapi.api_not_deployed"
)

func ErrorJobSpecNotFound() error {
	return errors.WithStack(&errors.Error{
		Kind:    ErrJobNotFound,
		Message: fmt.Sprintf("unable to find job specification file"),
	})
}

// TODO remove
func ErrorMissingJobState() error {
	return errors.WithStack(&errors.Error{
		Kind:    ErrMissingJobState,
		Message: fmt.Sprintf("missing job state"),
	})
}
