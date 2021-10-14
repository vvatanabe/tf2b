package cli

import (
	"os"

	"github.com/vvatanabe/tfnotify/io"

	"github.com/vvatanabe/tf2b/config"
	"github.com/vvatanabe/tfnotify/ci"
	"github.com/vvatanabe/tfnotify/errors"
	"github.com/vvatanabe/tfnotify/notifier/backlog"
	"github.com/vvatanabe/tfnotify/terraform"
)

type tf2b struct {
	config   config.Config
	ci       ci.CI
	title    string
	message  string
	parser   terraform.Parser
	template terraform.Template
}

// Run sends the notification with notifier
func (t *tf2b) Run() error {

	prNum := t.ci.PR.Number
	if t.config.PRNumber > 0 {
		prNum = t.config.PRNumber
	}

	cfg := backlog.Config{
		APIKey:  t.config.APIKey,
		BaseURL: t.config.BaseURL,
		Project: t.config.Project,
		Repo:    t.config.Repository,
		PR: backlog.PullRequest{
			Number:  prNum,
			Title:   t.title,
			Message: t.message,
		},
		CI:       t.ci.URL,
		Parser:   t.parser,
		Template: t.template,
	}
	client, err := backlog.NewClient(cfg)
	if err != nil {
		return err
	}
	notifier := client.Notify

	return errors.NewExitError(notifier.Notify(io.Tee(os.Stdin, os.Stdout)))
}
