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

package integrationplatform

import (
	"context"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/install"
	"github.com/sirupsen/logrus"
)

var resources = []string{
	"platform-integration-context-jvm.yaml",
	"platform-integration-context-groovy.yaml",
	"platform-integration-context-kotlin.yaml",
	"platform-integration-context-spring-boot.yaml",
}

var knativeResources = []string{
	"platform-integration-context-knative.yaml",
}

// NewCreateAction returns a action that creates resources needed by the platform
func NewCreateAction() Action {
	return &createAction{}
}

type createAction struct {
	baseAction
}

func (action *createAction) Name() string {
	return "create"
}

func (action *createAction) CanHandle(platform *v1alpha1.IntegrationPlatform) bool {
	return platform.Status.Phase == v1alpha1.IntegrationPlatformPhaseCreating
}

func (action *createAction) Handle(ctx context.Context, platform *v1alpha1.IntegrationPlatform) error {
	logrus.Info("Installing platform resources")
	err := install.Resources(ctx, action.client, platform.Namespace, resources...)
	if err != nil {
		return err
	}

	if platform.Spec.Profile == v1alpha1.TraitProfileKnative {
		logrus.Info("Installing knative resources")
		err := install.Resources(ctx, action.client, platform.Namespace, knativeResources...)
		if err != nil {
			return err
		}
	}

	target := platform.DeepCopy()
	target.Status.Phase = v1alpha1.IntegrationPlatformPhaseStarting
	logrus.Info("Platform ", target.Name, " transitioning to state ", target.Status.Phase)

	return action.client.Update(ctx, target)
}
