package main

import (
    "fmt"
    "log"
    "time"
    "strings"
    "strconv"
    "net/http"
    "gopkg.in/mgo.v2"
    "labix.org/v2/mgo/bson"
    "github.com/jbckmn/gopod"
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
)


func main() {

    ctl, err := NewController()
    if err != nil {
        log.Fatal(err)
    }
    defer ctl.session.Close()

    goji.Get("/:vanity", ctl.renderPodcast)
    goji.Get("/:vanity/", ctl.renderPodcast)
    goji.Serve()

}

func (ctl *Controller) renderPodcast(c web.C, w http.ResponseWriter, r *http.Request) {
    db := ctl.session.Clone().DB("audrey")
    result := Account{}
    err := db.C("accounts").Find(bson.M{"vanity": c.URLParams["vanity"]}).One(&result)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    } else {
        iter := db.C("articles").Find(bson.M{"accountId": result.Id.Hex()}).Sort("-created").Limit(100).Iter()
        s := buildPodcast(iter, result)
        w.Header().Set("Content-Type", "application/rss+xml")
        fmt.Fprintf(w, "%s", s.Publish())
    }
}

func buildPodcast(iter *mgo.Iter, acct Account) *gopod.Channel {
    var result Article
    title := []string{acct.Username, " on Narro"}
    desc := []string{acct.Username, " uses Narro to create a podcast of articles transcribed to audio."}
    link := []string{"http://on.narro.co/", acct.Vanity}
    c := gopod.ChannelFactory(strings.Join(title, ""), strings.Join(link, ""), strings.Join(desc, ""), "http://narro.co/images/narro-icon.png")

    c.SetTTL("20") // in minutes
    c.SetPubDate(time.Now().UTC().Format(time.RFC822))
    c.SetLanguage("en")
    c.SetWebMaster("josh@narro.co")
    c.SetManagingEditor(acct.Email)
    c.SetCopyright("All article content copyright of respective source authors.")
    c.SetiTunesExplicit("no")
    c.SetiTunesAuthor(acct.Email)
    c.SetiTunesSubtitle(strings.Join(desc, ""))
    c.SetiTunesSummary(strings.Join(desc, ""))
    c.SetiTunesOwner(acct.Username, acct.Email)

    for iter.Next(&result) {
        fmt.Println(result)
        resultLink := []string{"http://narro.co/article/", result.Id.Hex()}
        resultDesc := []string{result.Description, result.Url}
        i := &gopod.Item{
            Title:          result.Title,
            Link:           strings.Join(resultLink, ""),
            Description:    strings.Join(resultDesc, " ... "),
            PubDate:        result.Created.UTC().Format(time.RFC822),
            Author:         acct.Email,
            Guid:           strings.Join(resultLink, ""),
            TunesDuration:  strconv.FormatFloat(result.Mp3Length, 'f', 0, 64),
            TunesAuthor:    acct.Email,
            TunesSubtitle:  result.Title,
            TunesSummary:   strings.Join(resultDesc, " ... "),
            TunesExplicit:  "no",
        }
        i.SetEnclosure(result.Mp3URL, strconv.FormatFloat(result.Mp3Length, 'f', 0, 64), "audio/mpeg")
        c.AddItem(i)
    }
    if err := iter.Close(); err != nil {
        fmt.Printf(err.Error())
    }

    return c
}