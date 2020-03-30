package manifest

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/gojek/stevedore/log"
	"github.com/gojek/stevedore/pkg/helm"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/olekukonko/tablewriter"
)

var (
	red     = fmt.Sprintf
	green   = fmt.Sprintf
	yellow  = fmt.Sprintf
	nocolor = fmt.Sprintf
)

func init() {
	red = color.New(color.FgRed).SprintfFunc()
	green = color.New(color.FgGreen).SprintfFunc()
	yellow = color.New(color.FgYellow).SprintfFunc()
	nocolor = fmt.Sprintf
}

// Summarizer is the interface that wraps the
// grouped responses summary Display method
type Summarizer interface {
	Display(stevedore.GroupedResponses)
}

// TableSummarizer represents the default summarizer that
// displays the summary in a tabular format
type TableSummarizer struct {
	writer io.Writer
}

// Display outputs the tabular summary to the TableSummarizer.writer
func (s TableSummarizer) Display(group stevedore.GroupedResponses) {
	changedFilesTable := createTable(s.writer, []string{"FILENAME", "ADDITIONS", "MODIFICATIONS", "DESTRUCTIONS"}, false)
	helmReleaseTable := createTable(s.writer, []string{"RELEASE", "MANIFEST CHANGES"}, true)

	for _, file := range group.SortedFileNames() {
		responses := group[file]
		fileName := getFormattedFileName(file)
		responses.SortByReleaseName()

		fileHasDiff := false
		additions := 0
		modifications := 0
		deletions := 0

		for _, response := range responses {
			fileHasDiff = fileHasDiff || response.HasDiff

			summary := response.Summary()
			additions += len(summary.Added)
			modifications += len(summary.Modified)
			deletions += len(summary.Destroyed)

			if response.HasDiff {
				releaseDetail := nocolor("%s\n(%s)", response.ReleaseName, fileName)
				helmReleaseTable.Append(append([]string{releaseDetail}, formatRelease(summary)))
			}
		}

		if fileHasDiff {
			changedFilesTable.Append([]string{
				fileName,
				nocolor("%d", additions),
				nocolor("%d", modifications),
				nocolor("%d", deletions),
			})
		}
	}

	renderTable(s.writer, "Release changes:", helmReleaseTable)
	renderTable(s.writer, "File changes:", changedFilesTable)

}

func getFormattedFileName(file string) string {
	split := strings.Split(file, "/")
	return split[len(split)-1]
}

func createTable(writer io.Writer, headers []string, displayRowLine bool) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(headers)
	table.SetRowLine(displayRowLine)
	return table
}

func renderTable(writer io.Writer, tableTitle string, table *tablewriter.Table) {
	if table.NumLines() <= 0 {
		return
	}
	if _, err := fmt.Fprintln(writer, tableTitle); err != nil {
		log.Error("error in TableSummarizer when printing to ", writer)
	}
	table.Render()
}

func formatRelease(summary helm.Summary) string {
	var allRows []string
	allRows = append(allRows, groupAndSortResourcesByKind(summary.Added, green, "+")...)
	allRows = append(allRows, groupAndSortResourcesByKind(summary.Modified, yellow, "~")...)
	allRows = append(allRows, groupAndSortResourcesByKind(summary.Destroyed, red, "-")...)
	return strings.Join(allRows, "\n")
}

func groupAndSortResourcesByKind(resources helm.Resources, colorize func(format string, a ...interface{}) string, prefix string) []string {
	rows := make([]string, 0)
	for kind, resources := range resources.GroupByKind() {
		for _, resource := range resources {
			rows = append(rows, colorize("%s%s/%s", prefix, kind, resource.Name))
		}
	}
	sort.Strings(rows)
	return rows
}
