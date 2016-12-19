package githubstat

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ProxyClient struct {
	client *github.Client
}

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

// create a github client only once.
// call Client() and create client only once.
func (c *ProxyClient) getClient() *github.Client {
	if nil == c.client {
		ts := &tokenSource{
			&oauth2.Token{AccessToken: Config.AccessToken},
		}

		tc := oauth2.NewClient(oauth2.NoContext, ts)
		c.client = github.NewClient(tc)
	}

	return c.client
}
