package cli

import (
	"fmt"
	"strings"

	"github.com/vvatanabe/tfnotify/notifier/backlog"

	"github.com/urfave/cli/v2"
	"github.com/vvatanabe/tf2b/config"
	"github.com/vvatanabe/tf2b/git"
	"github.com/vvatanabe/tf2b/template"
	"github.com/vvatanabe/tfnotify/ci"
	"github.com/vvatanabe/tfnotify/terraform"
)

func cmdPlan(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}

	if err := parseOpts(ctx, &cfg); err != nil {
		return err
	}

	c, err := resolveCI(&cfg)
	if err != nil {
		return err
	}

	if err := cfg.Validation(); err != nil {
		return err
	}

	t := &tf2b{
		config:   cfg,
		ci:       c,
		title:    cfg.Terraform.Plan.Title,
		message:  cfg.Terraform.Plan.Message,
		parser:   terraform.NewPlanParser(),
		template: terraform.NewPlanTemplate(tmplPlan(&cfg)),
	}
	return t.Run()
}

func tmplPlan(cfg *config.Config) string {
	if cfg.Terraform.Plan.Template != "" {
		return cfg.Terraform.Plan.Template
	}
	switch cfg.Format {
	case config.Backlog:
		return template.DefaultPlanTemplateWithBacklog
	case config.Markdown:
		return template.DefaultPlanTemplateWithMarkdown
	default:
		return template.DefaultPlanTemplateWithMarkdown
	}
}

func cmdApply(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}

	if err := parseOpts(ctx, &cfg); err != nil {
		return err
	}

	c, err := resolveCI(&cfg)
	if err != nil {
		return err
	}

	if err := cfg.Validation(); err != nil {
		return err
	}

	t := &tf2b{
		config:   cfg,
		ci:       c,
		title:    cfg.Terraform.Apply.Title,
		message:  cfg.Terraform.Apply.Message,
		parser:   terraform.NewApplyParser(),
		template: terraform.NewApplyTemplate(tmplApply(&cfg)),
	}
	return t.Run()
}

func tmplApply(cfg *config.Config) string {
	if cfg.Terraform.Apply.Template != "" {
		return cfg.Terraform.Apply.Template
	}
	switch cfg.Format {
	case config.Backlog:
		return template.DefaultApplyTemplateWithBacklog
	case config.Markdown:
		return template.DefaultApplyTemplateWithMarkdown
	default:
		return template.DefaultApplyTemplateWithMarkdown
	}
}

func newConfig(ctx *cli.Context) (config.Config, error) {
	cfg := config.Config{}
	confPath, err := cfg.Find(ctx.String("config"))
	if err != nil {
		return cfg, err
	}
	if confPath != "" {
		if err := cfg.LoadFile(confPath); err != nil {
			return cfg, err
		}
	}

	if cfg.APIKey == "" {
		cfg.APIKey = "$" + backlog.EnvApiKey
	}
	return cfg, nil
}

func resolveCI(cfg *config.Config) (ci.CI, error) {
	name := strings.ToLower(cfg.CI)
	switch name {
	case "":
		return local(cfg)
	case "circleci", "circle-ci":
		return ci.Circleci()
	case "travis", "travisci", "travis-ci":
		return ci.Travisci()
	case "codebuild":
		return ci.Codebuild()
	case "teamcity":
		return ci.Teamcity()
	case "drone":
		return ci.Drone()
	case "jenkins":
		return ci.Jenkins()
	case "gitlabci", "gitlab-ci":
		return ci.Gitlabci()
	case "github-actions":
		return ci.GithubActions()
	case "cloud-build", "cloudbuild":
		return ci.Cloudbuild()
	default:
		return ci.CI{}, fmt.Errorf("CI service %v: not supported yet", name)
	}
}

func local(cfg *config.Config) (ci ci.CI, err error) {
	repo, err := git.OpenRepository(".")
	if err != nil {
		return
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = repo.BaseURL()
	}
	if cfg.Project == "" {
		cfg.Project = repo.Project()
	}
	if cfg.Repository == "" {
		cfg.Repository = repo.Name()
	}
	if cfg.PRNumber <= 0 {
		cfg.PRNumber, err = repo.PullRequestNumberOfCurrentBranch()
	}
	return
}
