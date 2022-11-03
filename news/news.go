package news

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	collectionInterval = 1 * time.Minute

	finnhubNewsURL     = "https://finnhub.io/api/v1/company-news"
	finnHubTokenHeader = "X-Finnhub-Token"
	apiToken           = "cdgr05qad3i2r375f74gcdgr05qad3i2r375f750"
	apiDateFormat      = "2006-01-02"
	snowflakeTicker    = "AAPL"

	newsDatetimeFormat = "January, 2 2006 15:04PM"
)

// News represents a news story in the Finnhub format,
// with only the fields we're interested in.
// See https://finnhub.io/docs/api/company-news for details
type News struct {
	Datetime int64  `json:"datetime"`
	Headline string `json:"headline"`
	Id       int64  `json:"id"`
	Source   string `json:"source"`
	Summary  string `json:"summary"`
	Url      string `json:"url"`
}

// NewsCollector regularly collects news via HTTP requests
// and caches the responses. The cache can be read by clients
// to fetch recent news.
type NewsCollector struct {
	client *http.Client

	cache []News
	sync.RWMutex
}

func NewNewsCollector() *NewsCollector {
	return &NewsCollector{
		client: http.DefaultClient,
	}
}

func (c *NewsCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(collectionInterval)

	c.collect() // run once upon init
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.collect()
			}
		}
	}()
}

// RecentNews returns the n most recent news stories
// stored in the cache. If the cache has less than n
// entries, the entire cache contents is returned.
func (c *NewsCollector) RecentNews(n int) []News {
	c.RLock()
	defer c.RUnlock()

	newsCount := min(len(c.cache), n)

	// assume items stored in cache are in chronological order (latest news first)
	return c.cache[:newsCount]
}

func (c *NewsCollector) collect() {
	req, err := buildNewsRequest()
	if err != nil {
		fmt.Printf("Failed to build request: %s\n", err.Error())
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("Failed to execute request %v: %v\n", req, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Got unexpected response status code: %d: %v \n", resp.StatusCode, resp)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	var newsStories []News
	err = json.Unmarshal(body, &newsStories)
	if err != nil {
		fmt.Printf("Failed to unmarshal response body %s: %v\n", string(body), err)
		return
	}

	c.setCache(newsStories)
}

func buildNewsRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", finnhubNewsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new request: %w", err)
	}

	req.Header.Add(finnHubTokenHeader, apiToken)

	from, to := lastThreeDayRange()
	req.URL.RawQuery = url.Values{
		"from":   {from},
		"to":     {to},
		"symbol": {snowflakeTicker},
	}.Encode()

	return req, nil
}

func (c *NewsCollector) setCache(newsStories []News) {
	c.Lock()
	defer c.Unlock()
	c.cache = newsStories
}

func lastThreeDayRange() (string, string) {
	from := time.Now().Add(-time.Hour * 24 * 3).Format(apiDateFormat)
	to := time.Now().Format(apiDateFormat)

	return from, to
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func FormatNews(news []News) string {
	var sb strings.Builder
	sb.WriteString("Recent News Stories\n")
	for _, n := range news {
		parsedDatetime := parseDatetime(n.Datetime)

		sb.WriteString("===============================\n\n")
		sb.WriteString(fmt.Sprintf("%s -- %s\n", n.Headline, parsedDatetime))
		sb.WriteString("-------------------------------\n")
		sb.WriteString(fmt.Sprintf("Summary: %s\n", n.Summary))
		sb.WriteString(fmt.Sprintf("Read more: %s\n\n", n.Url))
	}

	return sb.String()
}

func parseDatetime(unixTime int64) string {
	t := time.Unix(unixTime, 0)

	return t.Format(newsDatetimeFormat)
}
