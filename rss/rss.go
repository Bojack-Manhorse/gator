package rss

import (
	"html"

	"github.com/Bojack-Manhorse/pokedexcli/gator/utils"
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

func (item *RSSItem) processItem() {
	item.Title = html.UnescapeString(item.Title)
	item.Link = html.UnescapeString(item.Link)
	item.Description = html.UnescapeString(item.Description)
	item.PubDate = html.UnescapeString(item.PubDate)
}

func (feed *RSSFeed) processFeed() {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for idx := range feed.Channel.Item {
		feed.Channel.Item[idx].processItem()
	}

}

func FetchFeed(feedUrl string) (*RSSFeed, error) {

	data, err := utils.GetBytesFromHTML(feedUrl, "gator")

	if err != nil {
		return nil, err
	}

	parsedXml, err := utils.ParseBytesToXML[RSSFeed](data)

	if err != nil {
		return nil, err
	}

	parsedXml.processFeed()

	return &parsedXml, nil

}
