package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/htr/feedparser"
	mm "github.com/sloonz/go-mime-message"
	qprintable "github.com/sloonz/go-qprintable"
)

type Feed struct {
	Url    string
	Folder string
}

func (f Feed) pull(itemProcessor func(*feedparser.Feed, *feedparser.FeedItem)) error {

	rdr, err := http.Get(f.Url)
	if err != nil {
		return err
	}

	feed, err := feedparser.NewFeed(rdr.Body)
	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		itemProcessor(feed, item)
	}
	return nil
}

func newFeed(url, folder string) Feed {
	return Feed{
		Url:    url,
		Folder: folder,
	}
}

func loadFeeds() []Feed {
	feedsList := []Feed{}

	for folder, feeds := range ctx.Folders {
		for _, feed := range feeds {
			feedsList = append(feedsList, newFeed(feed, ctx.FolderPrefix+folder))
		}
	}

	return feedsList
}

func itemToBody(feed *feedparser.Feed, item *feedparser.FeedItem) io.Reader {

	tmpl := template.Must(template.New("message").Parse(ctx.Template))

	type tmplData struct {
		Link    string
		Title   string
		Author  string
		Content template.HTML
	}

	data := tmplData{
		Link:  item.Link,
		Title: item.Title,
	}

	if len(item.Content) > 0 {
		data.Content = template.HTML(item.Content)
	} else {
		data.Content = template.HTML(item.Description)
	}

	var doc bytes.Buffer
	tmpl.Execute(&doc, data)
	return &doc
}

func itemToMsg(feed *feedparser.Feed,
	item *feedparser.FeedItem, body io.Reader) *mm.Message {

	msg := mm.NewTextMessage(qprintable.UnixTextEncoding, body)
	msg.SetHeader("Date", item.When.Format(time.RFC822Z))
	msg.SetHeader("Subject", mm.EncodeWord(item.Title))
	msg.SetHeader("Content-Type", "text/html; charset=utf-8")
	var author string
	if len(item.Author) > 0 {
		author = item.Author
	} else {
		author = feed.Title
	}
	msg.SetHeader("From", mm.EncodeWord(author)+" "+ctx.From)
	msg.SetHeader("To", ctx.To)
	return msg
}

func pullFeeds(feeds []Feed) {
	imapSession := ctx.Imap

	for _, feed := range feeds {

		err := imapSession.CreateFolder(feed.Folder)
		if err != nil {
			fmt.Println("unable to create folder", feed.Folder, " => ", err)
			continue
		}

		fmt.Println("checking...", feed.Url)

		msgDropper := func(f *feedparser.Feed, item *feedparser.FeedItem) {
			if appendLog.Exists(item.Id) {
				return
			}
			println("appending", item.Id, "to", feed.Folder)
			body := itemToBody(f, item)
			msg := itemToMsg(f, item, body)
			imapSession.Append(feed.Folder, ctx.Labels, msg)
			appendLog.Add(item.Id)
		}
		feed.pull(msgDropper)
		appendLog.Save()
	}
}
