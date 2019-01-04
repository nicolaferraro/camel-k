/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package install

import (
	"context"
	"github.com/apache/camel-k/deploy"
	"github.com/apache/camel-k/pkg/client"
	"github.com/apache/camel-k/pkg/util/kubernetes"
	"github.com/apache/camel-k/pkg/util/kubernetes/customclient"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// SetupClusterwideResources --
func SetupClusterwideResources(ctx context.Context, c client.Client, ) error {
	return SetupClusterwideResourcesOrCollect(ctx, c, nil)
}

// SetupClusterwideResourcesOrCollect --
func SetupClusterwideResourcesOrCollect(ctx context.Context, c client.Client, collection *kubernetes.Collection) error {

	// Install CRD for Integration Platform (if needed)
	if err := installCRD(ctx, c, "IntegrationPlatform", "crd-integration-platform.yaml", collection); err != nil {
		return err
	}

	// Install CRD for Integration Context (if needed)
	if err := installCRD(ctx, c, "IntegrationContext", "crd-integration-context.yaml", collection); err != nil {
		return err
	}

	// Install CRD for Integration (if needed)
	if err := installCRD(ctx, c, "Integration", "crd-integration.yaml", collection); err != nil {
		return err
	}

	// Installing ClusterRole
	clusterRoleInstalled, err := IsClusterRoleInstalled(ctx, c)
	if err != nil {
		return err
	}
	if !clusterRoleInstalled || collection != nil {
		err := installClusterRole(ctx, c, collection)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsCRDInstalled check if the given CRT kind is installed
func IsCRDInstalled(ctx context.Context, c client.Client, kind string) (bool, error) {
	lst, err := c.Discovery().ServerResourcesForGroupVersion("camel.apache.org/v1alpha1")
	if err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	for _, res := range lst.APIResources {
		if res.Kind == kind {
			return true, nil
		}
	}
	return false, nil
}

func installCRD(ctx context.Context, c client.Client, kind string, resourceName string, collection *kubernetes.Collection) error {
	crd := []byte(deploy.Resources[resourceName])
	if collection != nil {
		unstr, err := kubernetes.LoadRawResourceFromYaml(string(crd))
		if err != nil {
			return err
		}
		collection.Add(unstr)
		return nil
	}

	// Installing Integration CRD
	installed, err := IsCRDInstalled(ctx, c, kind)
	if err != nil {
		return err
	}
	if installed {
		return nil
	}

	crdJSON, err := yaml.ToJSON(crd)
	if err != nil {
		return err
	}
	restClient, err := customclient.GetClientFor(c, "apiextensions.k8s.io", "v1beta1")
	if err != nil {
		return err
	}
	// Post using dynamic client
	result := restClient.
		Post().
		Body(crdJSON).
		Resource("customresourcedefinitions").
		Do()
	// Check result
	if result.Error() != nil && !errors.IsAlreadyExists(result.Error()) {
		return result.Error()
	}

	return nil
}

// IsClusterRoleInstalled check if cluster role camel-k:edit is installed
func IsClusterRoleInstalled(ctx context.Context, c client.Client, ) (bool, error) {
	clusterRole := v1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "camel-k:edit",
		},
	}
	key, err := k8sclient.ObjectKeyFromObject(&clusterRole)
	if err != nil {
		return false, err
	}
	err = c.Get(ctx, key, &clusterRole)
	if err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func installClusterRole(ctx context.Context, c client.Client, collection *kubernetes.Collection) error {
	obj, err := kubernetes.LoadResourceFromYaml(c.GetScheme(), deploy.Resources["user-cluster-role.yaml"])
	if err != nil {
		return err
	}

	if collection != nil {
		collection.Add(obj)
		return nil
	}
	return c.Create(ctx, obj)
}
