package rss

import (
	"testing"
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
			t.Error(
				"expected", expectedPubDate,
				"got", item.PubDate,
			)
		}
	}
}
