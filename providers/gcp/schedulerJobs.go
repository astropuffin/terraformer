// Copyright 2018 The Terraformer Authors.
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

package gcp

import (
	"context"
	"log"
	"strings"

	"google.golang.org/api/cloudscheduler/v1beta1"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"golang.org/x/oauth2/google"
)

var schedulerJobsAllowEmptyValues = []string{""}

var schedulerJobsAdditionalFields = map[string]string{}

type SchedulerJobsGenerator struct {
	GCPService
}

// Run on SchedulerJobsList and create for each TerraformResource
func (g SchedulerJobsGenerator) createResources(jobsList *cloudscheduler.ProjectsLocationsJobsListCall, ctx context.Context) []terraform_utils.Resource {
	resources := []terraform_utils.Resource{}
	if err := jobsList.Pages(ctx, func(page *cloudscheduler.ListJobsResponse) error {
		for _, obj := range page.Jobs {
			t := strings.Split(obj.Name, "/")
			name := t[len(t)-1]
			resources = append(resources, terraform_utils.NewResource(
				obj.Name,
				name,
				"google_cloud_scheduler_job",
				"google",
				map[string]string{
					"name":    name,
					"project": g.GetArgs()["project"],
					"region":  g.GetArgs()["region"],
				},
				schedulerJobsAllowEmptyValues,
				schedulerJobsAdditionalFields,
			))
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return resources
}

// Generate TerraformResources from GCP API,
func (g *SchedulerJobsGenerator) InitResources() error {
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, cloudscheduler.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	cloudSchedulerService, err := cloudscheduler.New(c)
	if err != nil {
		log.Fatal(err)
	}

	jobsList := cloudSchedulerService.Projects.Locations.Jobs.List("projects/" + g.GetArgs()["project"] + "/locations/" + g.GetArgs()["region"])

	g.Resources = g.createResources(jobsList, ctx)
	g.PopulateIgnoreKeys()
	return nil
}
