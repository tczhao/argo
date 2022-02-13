package cron

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/argoproj/argo-workflows/v3/cmd/argo/commands/client"
	cronworkflowpkg "github.com/argoproj/argo-workflows/v3/pkg/apiclient/cronworkflow"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/workflow/util"
)

type cliUpdateOpts struct {
	strict   bool   // --strict
}

func NewUpdateCommand() *cobra.Command {
	var (
		cliUpdateOpts cliUpdateOpts
		submitOpts    wfv1.SubmitOpts
	)
	command := &cobra.Command{
		Use:   "udpate FILE1 FILE2...",
		Short: "update a cron workflow",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			UpdateCronWorkflows(cmd.Context(), args, &cliUpdateOpts, &submitOpts)
		},
	}

	util.PopulateSubmitOpts(command, &submitOpts, false)
	command.Flags().BoolVar(&cliUpdateOpts.strict, "strict", true, "perform strict workflow validation")
	return command
}

func UpdateCronWorkflows(ctx context.Context, filePaths []string, cliOpts *cliUpdateOpts, submitOpts *wfv1.SubmitOpts) {
	ctx, apiClient := client.NewAPIClient(ctx)
	serviceClient, err := apiClient.NewCronWorkflowServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	fileContents, err := util.ReadManifest(filePaths...)
	if err != nil {
		log.Fatal(err)
	}

	var cronWorkflows []wfv1.CronWorkflow
	for _, body := range fileContents {
		cronWfs := unmarshalCronWorkflows(body, cliOpts.strict)
		cronWorkflows = append(cronWorkflows, cronWfs...)
	}

	if len(cronWorkflows) == 0 {
		log.Println("No CronWorkflows found in given files")
		os.Exit(1)
	}

	for _, cronWf := range cronWorkflows {

		newWf := wfv1.Workflow{Spec: cronWf.Spec.WorkflowSpec}
		err := util.ApplySubmitOpts(&newWf, submitOpts)
		if err != nil {
			log.Fatal(err)
		}
		cronWf.Spec.WorkflowSpec = newWf.Spec
		// We have only copied the workflow spec to the cron workflow but not the metadata
		// that includes name and generateName. Here we copy the metadata to the cron
		// workflow's metadata and remove the unnecessary and mutually exclusive part.
		if generateName := newWf.ObjectMeta.GenerateName; generateName != "" {
			cronWf.ObjectMeta.GenerateName = generateName
			cronWf.ObjectMeta.Name = ""
		}
		if name := newWf.ObjectMeta.Name; name != "" {
			cronWf.ObjectMeta.Name = name
			cronWf.ObjectMeta.GenerateName = ""
		}
		if cronWf.Namespace == "" {
			cronWf.Namespace = client.Namespace()
		}
		updated, err := serviceClient.UpdateCronWorkflow(ctx, &cronworkflowpkg.UpdateCronWorkflowRequest{
			Namespace:    cronWf.Namespace,
			CronWorkflow: &cronWf,
		})
		if err != nil {
			log.Fatalf("Failed to update cron workflow: %v", err)
		}
		fmt.Print(getCronWorkflowGet(updated))
	}
}
