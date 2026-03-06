package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/C4triplezero/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	feed := &RSSFeed{}
	err = xml.Unmarshal(data, feed)
	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return feed, nil
}

func scrapeFeeds(s *state) {
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")

	_, err = s.DB.MarkedFeedFetch(context.Background(), feed.ID)
	if err != nil {
		log.Println("Couldn't mark feed as fetched", err)
		return
	}
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Println("Couldn't fetch from feed url", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		validDesc := item.Description != ""
		date, err := time.Parse(time.RFC1123Z, item.PubDate)
		validPub := err == nil
		_, err = s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  validDesc,
			},
			PublishedAt: sql.NullTime{
				Time:  date,
				Valid: validPub,
			},
			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
