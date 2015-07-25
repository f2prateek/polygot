package polygot

import (
	"log"
	"strings"

	"github.com/f2prateek/go-counter"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Polygot struct {
	client *github.Client
}

func New(token string) *Polygot {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	return &Polygot{client}
}

func (p *Polygot) Counts() (map[string]int64, error) {
	done := make(chan struct{})
	defer close(done)

	// Get the authenticated user.
	user, _, err := p.client.Users.Get("")
	log.Println("fetching user")

	if err != nil {
		return nil, err
	}

	eventC := p.events(*user.Login)
	return p.count(eventC)
}

// Publish events performed by the user to the `eventC` channel.
func (p *Polygot) events(user string) <-chan github.Event {
	// Buffered so we can fetch events while looking up language for a repository.
	eventC := make(chan github.Event, 10)

	go func() {
		defer close(eventC)

		opt := &github.ListOptions{}
		for page := 1; ; page++ {
			opt.Page = page
			log.Println("fetching activity for", user, "on page", opt.Page)
			events, resp, err := p.client.Activity.ListEventsPerformedByUser(user, false, opt)
			if err != nil {
				break
			}

			for _, event := range events {
				eventC <- event
			}

			if resp.NextPage == 0 {
				break
			}
		}
	}()

	return eventC
}

func (p *Polygot) count(eventC <-chan github.Event) (map[string]int64, error) {
	c := counter.New()
	values := make(map[string]string)

	for event := range eventC {
		l, ok := values[*event.Repo.Name]
		if !ok {
			log.Println("fetching repository information for", *event.Repo.Name)
			parts := strings.Split(*event.Repo.Name, "/")
			repo, _, err := p.client.Repositories.Get(parts[0], parts[1])
			if err != nil {
				return c.Values(), err
			}
			if repo.Language == nil {
				l = "unknown"
			} else {
				l = *repo.Language
			}
			values[*event.Repo.Name] = l
		}

		c.Increment(l)
	}

	return c.Values(), nil
}
