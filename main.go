package main

import (
    "os"
    "fmt"
    "time"
    "net/http"
    "gopkg.in/mgo.v2"
    "github.com/jbckmn/gopod"
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
)

var session = dbConnect()

func main() {

    goji.Get("/:vanity", renderPodcast)
    goji.Serve()

}

func renderPodcast(c web.C, w http.ResponseWriter, r *http.Request) {
    defer session.Close()
    db := session.Clone().DB("audrey")
    if db != nil {
        fmt.Println("db session")
    }
    w.Header().Set("Content-Type", "text/plain") // for testing
    // w.Header().Set("Content-Type", "application/rss+xml")
    s := buildPodcast()
    fmt.Fprintf(w, "%s", s.PublishIndent())
}

func buildPodcast() *gopod.Channel {
        c := gopod.ChannelFactory("Daniel's Channel", "http://RubyDeveloper.com/", "My Blog", "image.png")
        c.SetPubDate(time.Now().UTC())

        c.AddItem(&gopod.Item{
                Title:       "Stack Overflow",
                Link:        "http://stackoverflow.com/users/1305696/daniel",
                Description: "Stack Overflow",
                PubDate:     time.Now().UTC().Format(time.RFC1123),
        })

        i := &gopod.Item{
                Title:       "LinkedIn",
                Link:        "http://www.linkedin.com/in/dangogh",
                Description: "My LinkedIn",
        }
        i.SetPubDate(time.Now().Unix())
        c.AddItem(i)

        return c
        // return c.Publish()
}

func dbConnect() *mgo.Session {
    u := os.Getenv("MONGO_URI")
    session, err := mgo.Dial(u)
    if err != nil {
        fmt.Println("No db session")
        panic(err)
    }
    return session
}