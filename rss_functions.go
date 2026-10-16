package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/Solaire89/gator/internal/database"
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
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	// Common practice to identify your software to the server
	req.Header.Set("User-Agent", "gator")
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return &RSSFeed{}, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("error getting feed: ", err)
	}
	feedData, err := fetchFeed(context.Background(), feed.FeedUrl)
	if err != nil {
		log.Printf("couldn't collect feed %s: %v\n", feed.FeedName, err)
		return
	}

	for _, rssItem := range feedData.Channel.Item {
		var nullTime sql.NullTime
		parseTime, err := time.Parse(time.RFC1123, rssItem.PubDate)
		if err != nil {
			log.Printf("error parsing time: %v\n", err)
		} else {
			nullTime.Time = parseTime
			nullTime.Valid = true
		}
		now := time.Now()
		recordTime := sql.NullTime{
			Time:  now,
			Valid: true,
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   recordTime,
			UpdatedAt:   recordTime,
			Title:       rssItem.Title,
			Url:         rssItem.Link,
			Description: rssItem.Description,
			PublishedAt: nullTime,
			FeedID:      feed.ID,
		})
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" {
				continue
			} else {
				log.Printf("error creating post: %v\n", err)
				continue
			}
		}
		if err != nil {
			log.Printf("error creating post: %v\n", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found\n", feed.FeedName, len(feedData.Channel.Item))
}
