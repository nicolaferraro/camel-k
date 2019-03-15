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

package quarkus

import (
	"github.com/apache/camel-k/pkg/builder"
	"github.com/apache/camel-k/pkg/util/maven"
	"github.com/pkg/errors"
	"os"
	"path"
)

// DefaultSteps --
var DefaultSteps = []builder.Step{
	//builder.NewStep("project/generate", builder.ProjectGenerationPhase, builder.GenerateProject),
	//builder.NewStep("project/inject-dependencies", builder.ProjectGenerationPhase+1, builder.InjectDependencies),
	//builder.NewStep("project/sanitize-dependencies", builder.ProjectGenerationPhase+2, builder.SanitizeDependencies),
	builder.NewStep("project/quarkus-dependencies", builder.ProjectGenerationPhase+3, InjectDependencies),
	//builder.NewStep("build/compute-dependencies", builder.ProjectBuildPhase, builder.ComputeDependencies),
	builder.NewStep("build/quarkus-package", builder.ProjectBuildPhase, BuildProject),
	builder.NewStep("build/compute-dependencies", builder.ProjectBuildPhase+1, ComputeDependencies),
	//builder.NewStep("packager/incremental", builder.ApplicationPackagePhase, builder.IncrementalPackager),
	//builder.NewStep("publisher/s2i", builder.ApplicationPublishPhase, Publisher),
}

// BuildProject --
func BuildProject(ctx *builder.Context) error {
	p := path.Join(ctx.Path, "maven")

	err := maven.CreateStructure(p, ctx.Project)
	if err != nil {
		return err
	}

	// Work-around native-image.properties generation, as there is no classes to compile
	classesPath := path.Join(p, "target", "classes")
	err = os.MkdirAll(classesPath, 0777)
	if err != nil {
		return err
	}

	opts := make([]string, 0, 2)
	opts = append(opts, maven.ExtraOptions(ctx.Request.Platform.Build.LocalRepository)...)
	opts = append(opts, "package")

	err = maven.Run(p, opts...)
	if err != nil {
		return errors.Wrap(err, "failure while building Quarkus project")
	}

	return nil
}