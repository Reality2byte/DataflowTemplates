/*
 * Copyright (C) 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License. You may obtain a copy of
 * the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations under
 * the License.
 */

// The run-it-smoke-tests workflow runs all Direct Runner tests.
package main

import (
	"flag"
	"log"

	"github.com/GoogleCloudPlatform/DataflowTemplates/cicd/internal/flags"
	"github.com/GoogleCloudPlatform/DataflowTemplates/cicd/internal/workflows"
)

func main() {
	flags.RegisterCommonFlags()
	flags.RegisterItFlags()
	flag.Parse()

	// Run mvn install before running integration tests
	mvnFlags := workflows.NewMavenFlags()
	err := workflows.MvnCleanInstall().Run(
		mvnFlags.IncludeDependencies(),
		mvnFlags.SkipDependencyAnalysis(),
		mvnFlags.SkipCheckstyle(),
		mvnFlags.SkipJib(),
		mvnFlags.SkipTests(),
		mvnFlags.SkipJacoco(),
		mvnFlags.SkipShade(),
		mvnFlags.ThreadCount(8),
		mvnFlags.InternalMaven())
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// Run integration tests
	mvnFlags = workflows.NewMavenFlags()
	err = workflows.MvnVerify().Run(
		mvnFlags.IncludeDependencies(),
		mvnFlags.SkipDependencyAnalysis(),
		mvnFlags.SkipCheckstyle(),
		mvnFlags.SkipJib(),
		mvnFlags.RunIntegrationSmokeTests(),
		mvnFlags.ThreadCount(8),
		mvnFlags.IntegrationTestParallelism(4),
		mvnFlags.StaticBigtableInstance("teleport"),
		mvnFlags.StaticSpannerInstance("teleport"),
		mvnFlags.InternalMaven(),
		flags.Region(),
		flags.Project(),
		flags.ArtifactBucket(),
		flags.StageBucket(),
		flags.PrivateConnectivity(),
		flags.SpannerHost(),
		flags.FailureMode(),
		flags.RetryFailures(),
		flags.StaticOracleHost(),
		flags.StaticOracleSysPassword(),
		flags.CloudProxyHost(),
		flags.CloudProxyMySqlPort(),
		flags.CloudProxyPostgresPort(),
		flags.CloudProxyPassword())
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	log.Println("Build Successful!")
}
