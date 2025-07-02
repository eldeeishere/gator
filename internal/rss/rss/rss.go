package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// This function would typically fetch the RSS feed from the given URL.

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Gator RSS Client/1.0")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch feed: %s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	xmlFeed := &RSSFeed{}

	if err := xml.Unmarshal(data, xmlFeed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RSS feed: %w", err)
	}
	xmlFeed.Channel.Title = html.UnescapeString(xmlFeed.Channel.Title)
	xmlFeed.Channel.Description = html.UnescapeString(xmlFeed.Channel.Description)
	for i, item := range xmlFeed.Channel.Items {
		xmlFeed.Channel.Items[i].Title = html.UnescapeString(item.Title)
		xmlFeed.Channel.Items[i].Description = html.UnescapeString(item.Description)
	}

	return xmlFeed, nil

}
