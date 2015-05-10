# github-stat

tutorial for myself to write in golang and to use [go-github](https://github.com/google/go-github)

## Usage

1. `git clone https://github.com/yshnb/github-stat.git`
2. `cd ./github-stat`
3. `cp config.json.dist config.json`
4. fill in valid `accessToken` generated in Github

After, you can run the below command
```
$ go run main.go --metrics=star google/go-github
target repository: google/go-github
metrics: star
star: 1079
```
like this.

## Available metrics

For now, it can only get the numbers of github star in each repository.

My plan ....

- issue
- pull_requests (open/closed)
- code changes(additions/deletions)
- etc.

