# github-contrib-stats

this repository is originally forked from https://github.com/yshnb/github-stat but with different design purpose.
except statistics analysis of the numbers of merged/LGTM'ed/NonLGTM'ed(open) PRs in specified repos for specified user(s) , lots of other features are planed.

## Usage

1. `go get -u github.com/google/go-github/github`
2. `go get -u github.com/mgutz/ansi`
3. `go get -u golang.org/x/oauth2`
4. `go get -u github.com/olekukonko/tablewriter`
5. `git clone https://github.com/bruceauyeung/github-contrib-stats`
6. `cd ./github-contrib-stats`
7. `cp config.toml.dist config.toml`
8. fill in valid `accessToken` generated from https://github.com/settings/tokens (make sure **repo** scope is checked)
9. fill in other fields.

After, you can run the below command
```
$ go run main.go
metrics: pull request stat analysis
kubernetes/charts : listing open pull requests
kubernetes/charts : listing closed pull requests
 ......(some outputs ignored)
+--------------+--------+---------+------------+
|     NAME     | MERGED | LGTM'ED | NONLGTM'ED |
+--------------+--------+---------+------------+
| bruceauyeung |     10 |       0 |          6 |
| tanshanshan  |      8 |       0 |          6 |
+--------------+--------+---------+------------+
```
like this.