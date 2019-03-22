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
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/builder"
	"github.com/apache/camel-k/pkg/util/maven"
)

const Version = "0.11.0"

// InjectDependencies --
func InjectDependencies(ctx *builder.Context) error {
	// Add Quarkus BOM
	ctx.Project.DependencyManagement.Dependencies = append(ctx.Project.DependencyManagement.Dependencies, maven.Dependency{
		GroupID:     "io.quarkus",
		ArtifactID:  "quarkus-bom",
			Version: Version,
			Type:    "pom",
			Scope:   "import",
	})

	// Add Quarkus dependencies
	ctx.Project.AddDependencies(
		maven.Dependency{
			GroupID:    "io.quarkus",
			ArtifactID: "quarkus-arc",
		},
		maven.Dependency{
			GroupID:    "org.apache.camel.k",
			ArtifactID: "camel-k-quarkus-deployment",
			Version:    ctx.Request.RuntimeVersion,
			Scope:      "provided",
		},
		maven.Dependency{
			GroupID:    "org.apache.camel.k",
			ArtifactID: "camel-k-quarkus-runtime",
			Version:    ctx.Request.RuntimeVersion,
			Scope:      "runtime",
		},
		maven.Dependency{
			GroupID:    "org.apache.camel.k",
			ArtifactID: "camel-k-runtime-quarkus",
			Version:    ctx.Request.RuntimeVersion,
			Scope:      "runtime",
		},
	)

	// Add Quarkus plugin
	ctx.Project.Build.Plugins = append(ctx.Project.Build.Plugins, maven.Plugin{
		GroupID:    "io.quarkus",
		ArtifactID: "quarkus-maven-plugin",
		Version:    Version,
		Executions:[]maven.Execution{
			{
				//ID: "build",
				Goals:[]string{
					"build",
				},
			},
		},
	})

	return nil
}

// ComputeDependencies --
func ComputeDependencies(ctx *builder.Context) error {
	target := path.Join(ctx.Path, "maven", "target")
	lib := path.Join(ctx.Path, "maven", "target", "lib")

	var files []string
	err := filepath.Walk(lib, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		return err
	}

	rex := regexp.MustCompile("([\\w-]+)-([\\w.]+(?:-[\\w]+)?).jar")

	for _, file := range files {
		s := rex.FindStringSubmatchIndex(file)
		if len(s) != 6 {
			continue
		}

		groupId := file[:s[2]-1]
		artifactId := file[s[2]:s[3]]
		version := file[s[4]:s[5]]

		ctx.Artifacts = append(ctx.Artifacts, v1alpha1.Artifact{
			ID:       groupId + ":" + artifactId + ":" + version,
			Location: path.Join(lib, file),
			Target:   path.Join("dependencies", file),
		})
	}

	runner := ctx.Project.ArtifactID + "-" + ctx.Project.Version + "-runner.jar"
	ctx.Artifacts = append(ctx.Artifacts, v1alpha1.Artifact{
		ID:       ctx.Project.GroupID + ":" + ctx.Project.ArtifactID + ":" + ctx.Project.Version,
		Location: path.Join(target, runner),
		Target:   runner,
	})

	return nil
}
