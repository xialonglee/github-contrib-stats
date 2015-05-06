package githubstat

import (
	"fmt"
	"strings"
)

type Format struct{}

func (f *Format) separatorOutput() {
	fmt.Println("------------------------------------------")
}

func (f *Format) fieldOutput(field, value string) {
	base := " " + field + " "
	padding := strings.Repeat(" ", 15-len(base))
	fmt.Printf("%s%s: %s\n", base, padding, value)
}

func (f *Format) FormatOutput(stats []ReposStat) {
	fmt.Print("\n")
	for i := 0; i < len(stats); i++ {
		f.separatorOutput()
		f.fieldOutput("repo", stats[i].Name)
		f.fieldOutput("open", fmt.Sprintf("%d", stats[i].OpenPullRequest))
		f.fieldOutput("closed", fmt.Sprintf("%d", stats[i].ClosedPullRequest))
	}
	f.separatorOutput()
}
