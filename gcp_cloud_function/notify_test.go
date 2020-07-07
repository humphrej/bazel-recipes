package notify

import (
	"context"
	"os"
	"testing"
)

var PubSubExample = []byte(`{"id": "TEST_ID"
,"projectId": "TEST_PROJECTID"
,"status": "SUCCESS"
,"startTime": "2020-07-01T16:46:44.568167054Z"
,"finishTime": "2020-07-01T16:50:42.144462Z"
,"logUrl": "https://console.cloud.google.com/cloud-build/builds/134da164-b72a-4395-9c34-9eeeb7fe7cda?project=389316"
,"images": [ "TEST_IMAGE" ]
,"source": {"repoSource": {"repoName": "TEST_REPO", "branchName": "TEST_BRANCH"}}
}
`)

func TestBuild(t *testing.T) {
	var context context.Context
	os.Setenv("EMAIL_FROM", "tom@foo.com")
	os.Setenv("EMAIL_TO", "jerry@foo.com")

	// replace send function with a mock
	var actualEmail *Email
	send = func(config Config, email *Email) error {
		actualEmail = email
		return nil
	}

	input := PubSubExample

	var message = PubSubMessage{Data: input}

	err := NotifyBuild(context, message)
	if err != nil {
		t.Errorf("Notification failed %s", err)
	}

	expectedEmail := Email{
		From:    "tom@foo.com",
		To:      "jerry@foo.com",
		Subject: "Build finished TEST_PROJECTID SUCCESS TEST_REPO/TEST_BRANCH (TEST_ID)",
		Text:    "Cloud Build TEST_ID finished with SUCCESS, in 3m57.576294946s.",
		Html: `<p>Cloud Build TEST_ID finished with SUCCESS, in 3m57.576294946s.</p>
<p><a href="https://console.cloud.google.com/cloud-build/builds/134da164-b72a-4395-9c34-9eeeb7fe7cda?project=389316">Build logs</a></p>
<p>Images: TEST_IMAGE</p>
`,
	}

	if *actualEmail != expectedEmail {
		t.Errorf("generated email was %s expected=%s", *actualEmail, expectedEmail)
	}
}

func TestEmailGeneration(t *testing.T) {

	config := Config{EmailFrom: "fred@foo.com", EmailTo: "wilma@foo.com"}
	build := Build{
		Id:         "TEST_ID",
		ProjectId:  "TEST_PROJECTID",
		Status:     "SUCCESS",
		StartTime:  "2000-11-29T11:00:00Z",
		FinishTime: "2000-11-29T11:01:00Z",
		LogUrl:     "https://foo.com/url1",
		Images:     []string{"TEST_IMAGE"},
		Source: Source{RepoSource: RepoSource{
			RepoName:   "TEST_REPO",
			BranchName: "TEST_BRANCH",
		}}}

	var email *Email
	var err error
	if email, err = CreateEmail(config, build); err != nil {
		t.Errorf("Failed generating email %s", err)
	}

	expectedEmail := Email{
		From:    "fred@foo.com",
		To:      "wilma@foo.com",
		Subject: "Build finished TEST_PROJECTID SUCCESS TEST_REPO/TEST_BRANCH (TEST_ID)",
		Text:    "Cloud Build TEST_ID finished with SUCCESS, in 1m0s.",
		Html: `<p>Cloud Build TEST_ID finished with SUCCESS, in 1m0s.</p>
<p><a href="https://foo.com/url1">Build logs</a></p>
<p>Images: TEST_IMAGE</p>
`,
	}

	if *email != expectedEmail {
		t.Errorf("generated email was %s expected=%s", *email, expectedEmail)
	}
}
