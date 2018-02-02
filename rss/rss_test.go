package rss

import (
	"testing"
)
import (
	rssHttp "rss_parser/http"
	"strings"
)

func TestSortItems(t *testing.T) {
	var rssExpected, rssTesting Rss
	rssExpected.Channel.Items = []Item{
		{PubDate: "Sun, 28 Jan 2018 12:04:00 +0400"},
		{PubDate: "Mon, 29 Jan 2018 11:04:05 +0400"},
		{PubDate: "Fri, 02 Feb 2018 05:03:05 +0400"},
	}
	rssTesting.Channel.Items = []Item{
		{PubDate: "Mon, 29 Jan 2018 11:04:05 +0400"},
		{PubDate: "Sun, 28 Jan 2018 12:04:00 +0400"},
		{PubDate: "Fri, 02 Feb 2018 05:03:05 +0400"},
	}

	rssTesting.SortItems()

	for i, item := range rssTesting.Channel.Items {
		expectedPubDate := rssExpected.Channel.Items[i].PubDate
		if item.PubDate != expectedPubDate {
			t.Fatal(
				"expected", expectedPubDate,
				"got", item.PubDate,
			)
		}
	}
}

func TestCheckCreateAndSendRequest(t *testing.T) {
	var rssTesting Rss
	request := rssHttp.CreateRequest("http://bash.im/rss/")

	targetUrl := "http://bash.im/rss/"
	receivedUrl := strings.Trim(request.URL.String(), "\n")
	if receivedUrl != targetUrl {
		t.Fatal(
			"expected", targetUrl,
			"got", receivedUrl,
		)
	}

	response := rssHttp.SendRequest(request)
	if response == nil {
		t.SkipNow()
	}
	rssTesting.DecodeXmlHttpResponse(response)

	targetValue := "Bash.im"
	receivedValue := strings.Trim(rssTesting.Channel.Title, "\n")
	if receivedValue != targetValue {
		t.Fatal("expected", targetValue,
			"got", receivedValue,
		)
	}
}

func TestSetDefaultAttributes(t *testing.T) {
	var rssTesting Rss
	rssTesting.setDefaultAttributes()

	type TestValues struct {
		expected string
		got      string
	}

	var tests = []TestValues{
		{
			expected: "http://localhost",
			got:      rssTesting.Channel.Link,
		},
		{
			expected: "ru-RU",
			got:      rssTesting.Channel.Language,
		},
		{
			expected: "Local Rss",
			got:      rssTesting.Channel.Title,
		},
	}

	for _, test := range tests {
		if test.expected != test.got {
			t.Fatal(
				"expected", test.expected,
				"got", test.got,
			)
		}
	}
}
