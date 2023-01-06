//go:build sweep
// +build sweep

package schemas

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/schemas"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
)

func init() {
	resource.AddTestSweepers("aws_schemas_discoverer", &resource.Sweeper{
		Name: "aws_schemas_discoverer",
		F:    sweepDiscoverers,
	})

	resource.AddTestSweepers("aws_schemas_registry", &resource.Sweeper{
		Name: "aws_schemas_registry",
		F:    sweepRegistries,
		Dependencies: []string{
			"aws_schemas_schema",
		},
	})

	resource.AddTestSweepers("aws_schemas_schema", &resource.Sweeper{
		Name: "aws_schemas_registry",
		F:    sweepSchemas,
	})
}

func sweepDiscoverers(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).SchemasConn()

	sweepResources := make([]sweep.Sweepable, 0)
	var sweeperErrs *multierror.Error

	input := &schemas.ListDiscoverersInput{}
	err = conn.ListDiscoverersPages(input, func(page *schemas.ListDiscoverersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, discoverer := range page.Discoverers {
			r := ResourceDiscoverer()
			d := r.Data(nil)
			d.SetId(aws.StringValue(discoverer.DiscovererId))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EventBridge Schemas Discoverer sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("listing EventBridge Schemas Discoverers: %w", err))
	}

	if err := sweep.SweepOrchestrator(sweepResources); err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("sweeping EventBridge Schemas Discoverers: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepRegistries(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).SchemasConn()

	sweepResources := make([]sweep.Sweepable, 0)
	var sweeperErrs *multierror.Error

	input := &schemas.ListRegistriesInput{}
	err = conn.ListRegistriesPages(input, func(page *schemas.ListRegistriesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, registry := range page.Registries {
			registryName := aws.StringValue(registry.RegistryName)
			if strings.HasPrefix(registryName, "aws.") {
				continue
			}

			r := ResourceRegistry()
			d := r.Data(nil)
			d.SetId(registryName)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EventBridge Schemas Registry sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("listing EventBridge Schemas Registries: %w", err))
	}

	if err := sweep.SweepOrchestrator(sweepResources); err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping EventBridge Schemas Registries: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepSchemas(region string) error {
	client, err := sweep.SharedRegionalSweepClient(region)
	if err != nil {
		return fmt.Errorf("getting client: %w", err)
	}
	conn := client.(*conns.AWSClient).SchemasConn()

	sweepResources := make([]sweep.Sweepable, 0)
	var sweeperErrs *multierror.Error

	input := &schemas.ListRegistriesInput{}
	err = conn.ListRegistriesPages(input, func(page *schemas.ListRegistriesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, registry := range page.Registries {
			registryName := aws.StringValue(registry.RegistryName)

			input := &schemas.ListSchemasInput{
				RegistryName: aws.String(registryName),
			}

			err = conn.ListSchemasPages(input, func(page *schemas.ListSchemasOutput, lastPage bool) bool {
				if page == nil {
					return !lastPage
				}

				for _, schema := range page.Schemas {
					schemaName := aws.StringValue(schema.SchemaName)
					if strings.HasPrefix(schemaName, "aws.") {
						log.Printf("[DEBUG] skipping EventBridge Schema (%s): managed by AWS", schemaName)
						continue
					}

					r := ResourceSchema()
					d := r.Data(nil)
					d.SetId(SchemaCreateResourceID(schemaName, registryName))

					sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
				}

				return !lastPage
			})
			if err != nil {
				sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("listing EventBridge Schemas Schemas from Registry %q: %w", registryName, err))
			}
		}

		return !lastPage
	})

	if sweep.SkipSweepError(err) {
		log.Printf("[WARN] Skipping EventBridge Schemas Schemas sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("listing EventBridge Schemas Schemas (%s): %w", region, err))
	}

	if err := sweep.SweepOrchestrator(sweepResources); err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping EventBridge Schemas Schemas (%s): %w", region, err))
	}

	return sweeperErrs.ErrorOrNil()
}
