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

package resources

import (
	"github.com/cortexlabs/cortex/pkg/lib/errors"
	"github.com/cortexlabs/cortex/pkg/operator/config"
	ok8s "github.com/cortexlabs/cortex/pkg/operator/k8s"
	"github.com/cortexlabs/cortex/pkg/operator/resources/batchapi"
	"github.com/cortexlabs/cortex/pkg/operator/resources/syncapi"
	"github.com/cortexlabs/cortex/pkg/types/spec"
	"github.com/cortexlabs/cortex/pkg/types/userconfig"
)

func FindDeployedResourceByName(resourceName string) (*userconfig.Resource, error) {
	virtualService, err := config.K8s.GetVirtualService(ok8s.Name(resourceName))
	if err != nil {
		return nil, err
	}

	if virtualService == nil {
		return nil, nil
	}

	return &userconfig.Resource{
		Name: virtualService.GetLabels()["apiName"],
		Kind: userconfig.KindFromString(virtualService.GetLabels()["apiKind"]),
	}, nil
}

func IsResourceUpdating(resource userconfig.Resource) (bool, error) {
	if resource.Kind == userconfig.SyncAPIKind {
		return syncapi.IsAPIUpdating(resource.Name)
	}

	return false, ErrorKindNotSupported(resource.Kind)
}

func UpdateAPI(apiConfig *userconfig.API, projectID string, force bool) (*spec.API, string, error) {
	deployedResource, err := FindDeployedResourceByName(apiConfig.Name)
	if err != nil {
		return nil, "", err
	}

	if deployedResource != nil && deployedResource.Kind != apiConfig.Kind {
		return nil, "", ErrorCannotChangeKindOfDeployedAPI(apiConfig.Name, apiConfig.Kind, deployedResource.Kind)
	}

	if apiConfig.Kind == userconfig.SyncAPIKind {
		return syncapi.UpdateAPI(apiConfig, projectID, force)
	}

	if apiConfig.Kind == userconfig.BatchAPIKind {
		return batchapi.UpdateAPI(apiConfig, projectID)
	}

	return nil, "", errors.Wrap(ErrorKindNotSupported(apiConfig.Kind), apiConfig.Identify()) // unexpected
}

func RefreshAPI(apiName string, force bool) (string, error) {
	deployedResource, err := FindDeployedResourceByName(apiName)
	if err != nil {
		return "", err
	} else if deployedResource == nil {
		return "", ErrorAPINotDeployed(apiName)
	}

	if deployedResource.Kind == userconfig.SyncAPIKind {
		return syncapi.RefreshAPI(apiName, force)
	}

	return "", errors.Wrap(ErrorKindNotSupported(deployedResource.Kind), deployedResource.Identify()) // unexpected
}

func DeleteAPI(apiName string, keepCache bool) error {
	deployedResource, err := FindDeployedResourceByName(apiName)
	if err != nil {
		return err
	} else if deployedResource == nil {
		return ErrorAPINotDeployed(apiName)
	}

	if deployedResource.Kind == userconfig.SyncAPIKind {
		return syncapi.DeleteAPI(apiName, keepCache)
	} else {
		return batchapi.DeleteAPI(apiName, keepCache) // TODO
	}
}