package githubstat

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var _client *github.Client

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

// create a github client only once.
// call Client() and create client only once.
func client() *github.Client {
	if nil == _client {
		ts := &tokenSource{
			&oauth2.Token{AccessToken: config.AccessToken},
		}

		tc := oauth2.NewClient(oauth2.NoContext, ts)
		_client = github.NewClient(tc)
	}

	return _client
}
