package cmd

import (
	"log"
	"os"
	"time"

	"github.com/clintjedwards/cursor/api"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var cmdPipelinesList = &cobra.Command{
	Use:   "list",
	Short: "List all pipelines",
	Run:   runPipelinesListCmd,
}

func runPipelinesListCmd(cmd *cobra.Command, args []string) {
	conn := initClientConn()
	client := api.NewCursorMasterClient(conn)
	ctx := generateClientContext()

	pipelines, err := client.ListPipelines(ctx, &api.ListPipelinesRequest{})
	if err != nil {
		log.Fatalf("could not list pipelines: %v", err)
	}

	tableData := [][]string{}

	for _, pipeline := range pipelines.Pipelines {
		tableData = append(tableData, []string{
			pipeline.Id,
			pipeline.Name,
			pipeline.Description,
			humanize.Time(time.Unix(pipeline.Created, 0)),
			humanize.Time(time.Unix(pipeline.LastCompiled, 0)),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "Created", "Last Compiled"})
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(tableData)
	table.Render()
}

func init() {
	cmdPipelines.AddCommand(cmdPipelinesList)
}
