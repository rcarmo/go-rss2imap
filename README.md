# go-rss2imap

A tweak of @htr's rss2imap to match my rss2imap feature set.

Things that need fixing:

* <strike>RSS items are always sent with the same From: address</strike>
* RSS items are inserted with the current date rather than the item's date
* Concurrency could be improved

## Quickstart

Checkout the repository (all dependencies are included) and issue the `make` command to build.

> The bundled `Makefile` also provides `pi` and `synology` targets for the Raspberry Pi and old Synology ARM5 NAS.

## Configuration

For the moment, you should use the provided template (configuration will be refactored to live in `~/.config` soon):

    cp rss2imap.json.sample ~/.rss2imap.json

edit the configuration file:

```javascript
{
    "cleanup": true,                // cleanup messages already seen and not flagged
    "folders": {                    // dictionary of folder names mapped to lists of feed urls
        "folder": [
            "feedurl1",
            "feedurl2"
        ],
        "folder2": [
            "feedurl3"
        ]
    },
    "imapserver": "imap.server.tld",  // imap server with tls support
    "from": "<email@server.tld",      // "From" header to be appended to each feeditem Author
    "folderprefix": "feeds/",         // a folder prefix. an empty string is accepted
    "to": "<email@server.tld>",       // "To" header.
    "labels": [],                     // additional labels/folders to where each new message should be appended
    "template": "<table>\n<tbody>\n<tr><td><a href=\"{{ .Link }}\">{{ .Title }}</a></td></tr>\n<hr />\n<tr><td>{{ .Author }}</td></tr>\n<tr><td>{{ .Content }}</td></tr>\n</tbody>\n</table>",     // the message body template
    "password": "password",           // imap username
    "username": "username"            // imap password
}

```


pull all feeds:

```
rss2imap -pull
```

just do a cleanup:

```
rss2imap -cleanup
```


test how a feed is processed:

```
rss2imap -test-feed=feedurl
```


Details
-------


### imap and gmail

When connected to a gmail imap server, rss2imap uses the `X-GM-EXT-1` features:

 * when *duplicating* messages to other labels/mailboxes: store X-GM-LABELS vs COPY
 * when deleting messages (cleanup): store X-GM-LABELS \Trash vs store \Deleted
 
 
TODO
----

* TODO
* ...








