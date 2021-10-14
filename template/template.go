package template

const (
	// DefaultPlanTitle is a default title for terraform plan
	DefaultPlanTitle = "Plan result"
	// DefaultApplyTitle is a default title for terraform apply
	DefaultApplyTitle = "Apply result"

	// DefaultPlanTemplateWithBacklog is a default template with Backlog for terraform plan
	DefaultPlanTemplateWithBacklog = `
** {{ .Title }}
{{ .Message }}
{{if .Result}}
*** Summary
{code}
{{ .Result }}
{/code}
{{end}}
*** Details
{code}
{{ .Body }}
{/code}
`

	// DefaultApplyTemplateWithBacklog is a default template with Backlog for terraform apply
	DefaultApplyTemplateWithBacklog = `
** {{ .Title }}
{{ .Message }}
{{if .Result}}
*** Summary
{code}
{{ .Result }}
{/code}
{{end}}
*** Details
{code}
{{ .Body }}
{/code}
`

	// DefaultPlanTemplateWithMarkdown is a default template with Markdown for terraform plan
	DefaultPlanTemplateWithMarkdown = `
## {{ .Title }}
{{ .Message }}
{{if .Result}}
### Summary
` + "```" + `
{{ .Result }}
` + "```" + `
{{end}}
### Details
` + "```" + `
{{ .Body }}
` + "```" + `
`

	// DefaultApplyTemplateWithMarkdown is a default template with Markdown for terraform apply
	DefaultApplyTemplateWithMarkdown = `
## {{ .Title }}
{{ .Message }}
{{if .Result}}
### Summary
` + "```" + `
{{ .Result }}
` + "```" + `
{{end}}
### Details
` + "```" + `
{{ .Body }}
` + "```" + `
`
)
