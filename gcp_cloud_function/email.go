package notify

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"
)

type Email struct {
	From    string
	To      string
	Subject string
	Text    string
	Html    string
}
type Config struct {
	EmailFrom string
	EmailTo   string
}

const TEXT = `<p>{{.MsgText}}</p>
<p><a href="{{.Build.LogUrl}}">Build logs</a></p>
<p>Images: {{.Images}}</p>
`

// CreateEmail is responsible for formatting the build notification email
func CreateEmail(config Config, build Build) (*Email, error) {
	var err error
	var duration *time.Duration

	if duration, err = humanizeDuration(build); err != nil {
		return nil, err
	}
	msgText := fmt.Sprintf("Cloud Build %s finished with %s, in %s.", build.Id, build.Status, *duration)

	var html *string
	if html, err = generateHtml(build, msgText); err != nil {
		return nil, err
	}

	email := Email{
		From: config.EmailFrom,
		To:   config.EmailTo,
		Subject: fmt.Sprintf("Build finished %s %s %s/%s (%s)",
			build.ProjectId,
			build.Status,
			build.Source.RepoSource.RepoName,
			build.Source.RepoSource.BranchName,
			build.Id),
		Text: msgText,
		Html: *html,
	}
	return &email, nil
}

func humanizeDuration(build Build) (*time.Duration, error) {
	var end, start time.Time
	var err error
	if end, err = time.Parse(time.RFC3339, build.FinishTime); err != nil {
		return nil, err
	}
	if start, err = time.Parse(time.RFC3339, build.StartTime); err != nil {
		return nil, err
	}
	duration := end.Sub(start)
	return &duration, nil
}

func generateHtml(build Build, msgText string) (*string, error) {
	var t = template.Must(template.New("email").Parse(TEXT))
	buffer := bytes.Buffer{}

	err := t.Execute(&buffer, struct {
		MsgText string
		Build   Build
		Images  string
	}{
		MsgText: msgText,
		Build:   build,
		Images:  strings.Join(build.Images, ","),
	})
	if err != nil {
		return nil, err
	}
	html := buffer.String()
	return &html, nil
}
