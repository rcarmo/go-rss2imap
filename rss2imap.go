package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/user"

	"github.com/htr/feedparser"
)

var appendLog *AppendLog
var ctx *Context

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}

func maybeConnect() {
	if ctx.Imap == nil || !ctx.Imap.IsConnected() {
		var err error
		ctx.Imap, err = imapConnect(ctx.ImapServer, ctx.Username, ctx.Password)
		maybePanic(err)
	}
}

func maybeDisconnect() {
	if ctx.Imap != nil && ctx.Imap.IsConnected() {
		ctx.Imap.Close()
	}
}

func pull() {
	usr, err := user.Current()
	maybePanic(err)
	homeDir := usr.HomeDir

	appendLog, err = loadAppendLog(homeDir + "/.rss2imap.log")
	maybePanic(err)

	defer appendLog.Save()

	maybeConnect()

	feeds := loadFeeds()
	pullFeeds(feeds)

	if ctx.Cleanup {
		cleanup()
	}
}

func cleanup() {
	foldersList := []string{}

	for folder, _ := range ctx.Folders {
		foldersList = append(foldersList, ctx.FolderPrefix+folder)
	}

	maybeConnect()
	ctx.Imap.imapCleanup(foldersList)
}

func testFeed(url string) {
	fmt.Println("url", url)
	f := newFeed(url, "")
	itemsLeft := 1
	f.pull(func(feed *feedparser.Feed, item *feedparser.FeedItem) {
		if itemsLeft == 0 {
			return
		}
		itemsLeft--

		body := itemToBody(feed, item)
		b := new(bytes.Buffer)
		b.ReadFrom(body)
		fmt.Println("\nID:", item.Id)
		fmt.Println("\nAuthor:", item.Author)
		fmt.Println("\nLink:", item.Link)
		fmt.Println("\nSubtitle:", item.Title)
		fmt.Println("\nDescription:", item.Description)
		fmt.Println("\nContent:", item.Content)
		fmt.Println("\n=================")
		fmt.Println(b.String())
		fmt.Println("\n")
	})

}

func main() {

	var pullFlag = flag.Bool("pull", false, "pulls the registered feeds")
	var cleanupFlag = flag.Bool("cleanup", false, "clean the known folders")
	var testFeedFlag = flag.String("test-feed", "", "test the feed in the given url")

	flag.Parse()

	usr, err := user.Current()
	maybePanic(err)
	homeDir := usr.HomeDir
	ctx, err = loadConfig(homeDir + "/.rss2imap.json")
	maybePanic(err)

	defer maybeDisconnect()

	if *pullFlag {
		pull()
	}

	if *cleanupFlag {
		cleanup()
	}

	if len(*testFeedFlag) > 0 {
		testFeed(*testFeedFlag)
	}
}
