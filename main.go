package main

import (
    "fmt"
    "github.com/jbckmn/gopod"
    "time"
    // "net/http"
)

func main() {
    c := gopod.ChannelFactory("Daniel's Channel", "http://RubyDeveloper.com/", "My Blog", "http://narro.co/images/narro-icon.png")
    c.SetPubDate(time.Now().UTC())
    c.SetiTunesExplicit("No")

    c.AddItem(&gopod.Item{
        Title:       "Ruby Developer",
        Link:        "http://RubyDeveloper.com/",
        Description: "Ruby Developer",
        PubDate:     time.Now().UTC().Format(time.RFC1123),
    })

    // Example: Using an item's SetPubDate method
    i := &gopod.Item{
        Title:       "LinkedIn",
        Link:        "http://www.linkedin.com/in/dangogh",
        Description: "My LinkedIn",
    }
    i.SetPubDate(time.Now().Unix())
    c.AddItem(i)

    fmt.Printf("%s\n\n", c.Publish())
    fmt.Printf("%s\n\n", c.PublishIndent())
}