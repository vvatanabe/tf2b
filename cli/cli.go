package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/vvatanabe/tf2b/config"
	"github.com/vvatanabe/tf2b/constant"
	"github.com/vvatanabe/tf2b/template"
)

func New() *cli.App {
	app := cli.NewApp()
	app.Name = "tf2b"
	app.Usage = "Notify terraform results to Backlog"
	app.Version = constant.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "ci", Usage: "name of CI to run tf2b"},
		&cli.StringFlag{Name: "config", Usage: "config path"},
		&cli.StringFlag{Name: "base-url", Usage: "base url for the backlog space. eg: foo.backlog.com"},
		&cli.StringFlag{Name: "project", Usage: "key of the project to which the repository belongs"},
		&cli.StringFlag{Name: "repo", Usage: "repository name"},
		&cli.IntFlag{Name: "pr", Usage: "pull request number"},
		&cli.StringFlag{Name: "format", Usage: "comment format. eg: markdown or backlog", Value: string(config.Markdown)},
	}
	app.Commands = []*cli.Command{
		{
			Name:  "plan",
			Usage: "Run terraform plan and post a comment to Backlog pull request",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "title, t",
					Usage: "Specify the title to use for notification",
					Value: template.DefaultPlanTitle,
				},
				&cli.StringFlag{
					Name:  "message, m",
					Usage: "Specify the message to use for notification",
				},
			},
			Action: cmdPlan,
		},
		{
			Name:  "apply",
			Usage: "Run terraform apply and post a comment to Backlog pull request",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "title, t",
					Usage: "Specify the title to use for notification",
					Value: template.DefaultApplyTitle,
				},
				&cli.StringFlag{
					Name:  "message, m",
					Usage: "Specify the message to use for notification",
				},
			},
			Action: cmdApply,
		},
		{
			Name:  "version",
			Usage: "Show version",
			Action: func(ctx *cli.Context) error {
				cli.ShowVersion(ctx)
				return nil
			},
		},
	}
	return app
}

func parseOpts(ctx *cli.Context, cfg *config.Config) error {

	if ci := ctx.String("ci"); ci != "" {
		cfg.CI = ci
	}

	if baseURL := ctx.String("base_url"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	if project := ctx.String("project"); project != "" {
		cfg.Project = project
	}

	if repo := ctx.String("repo"); repo != "" {
		cfg.Repository = repo
	}

	if pr := ctx.Int("pr"); pr != 0 {
		cfg.PRNumber = pr
	}

	if format := ctx.String("format"); format != "" {
		cfg.Format = config.TextFormat(format)
		switch cfg.Format {
		case config.Backlog: // ok
		case config.Markdown: // ok
		default:
			return fmt.Errorf("unknown format: %s", format)
		}
	}

	if title := ctx.String("title"); title != "" {
		cfg.Terraform.Plan.Title = title
		cfg.Terraform.Apply.Title = title
	}
	if message := ctx.String("message"); message != "" {
		cfg.Terraform.Plan.Message = message
		cfg.Terraform.Apply.Message = message
	}

	return nil
}
