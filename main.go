package main

import "./githubstat"

func main() {
	//	flags := os.Args()
	stats := githubstat.StatPullRequests()
	format := &githubstat.Format{}
	format.FormatOutput(stats)
}
