package main

import (
	"fmt"
	"github.com/jbckmn/gopod"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	ctl, err := NewController()
	if err != nil {
		log.Fatal(err)
	}

	goji.Get("/", http.RedirectHandler("https://www.narro.co", 301))
	goji.Get("/:vanity", ctl.renderPodcast)
	goji.Get("/:vanity/", ctl.renderPodcast)
	goji.Get("/:vanity/keywords/:topic", ctl.renderKeywordPodcast)
	goji.Get("/:vanity/keywords/:topic/", ctl.renderKeywordPodcast)
	goji.Serve()
}

func (ctl *Controller) renderPodcast(c web.C, w http.ResponseWriter, r *http.Request) {
	session, err := sqlx.Open("postgres", ctl.dbURI)
	if err != nil {
		log.Fatal(err)
		log.Fatal("unable to acquire DB connection")
		return
	}
	defer session.Close()
	result, accountErr := findRequestAccountByVanity(c, r, session)
	if accountErr != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	limit := 200
	rows, articlesErr := session.Queryx("SELECT _id, title, url, \"mp3URL\", \"mp3Length\", description, \"accountId\", created, links from articles WHERE active = true AND deleted = false and \"accountId\" = $1 ORDER BY created DESC LIMIT $2", result.Id, limit)
	if articlesErr != nil {
		log.Printf("unable to query DB for articles")
		http.Error(w, articlesErr.Error(), http.StatusInternalServerError)
		return
	}
	s := buildPodcast(rows, result, result.Vanity)
	w.Header().Set("Content-Type", "application/rss+xml")
	fmt.Fprintf(w, "%s", s.Publish())
}

func (ctl *Controller) renderKeywordPodcast(c web.C, w http.ResponseWriter, r *http.Request) {
	session, err := sqlx.Open("postgres", ctl.dbURI)
	if err != nil {
		log.Printf("unable to acquire DB connection")
		return
	}
	defer session.Close()
	result, accountErr := findRequestAccountByVanity(c, r, session)
	if accountErr != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	limit := 25
	topics := strings.Split(c.URLParams["topic"], "+")
	rows, articlesErr := session.Queryx("SELECT _id, title, url, \"mp3URL\", \"mp3Length\", description, \"accountId\", created, links from articles WHERE active = true AND deleted = false and \"accountId\" = $1 AND keywords ?& $2 ORDER BY created DESC LIMIT $3", result.Id, pq.Array(topics), limit)
	if articlesErr != nil {
		log.Printf("unable to query DB for articles")
		http.Error(w, articlesErr.Error(), http.StatusInternalServerError)
		return
	}
	name := []string{result.Vanity, c.URLParams["topic"]}
	s := buildPodcast(rows, result, strings.Join(name, "/"))
	w.Header().Set("Content-Type", "application/rss+xml")
	fmt.Fprintf(w, "%s", s.Publish())
}

func buildPodcast(iter *sqlx.Rows, acct Account, name string) *gopod.Channel {
	titlePieces := []string{name, " on Narro"}
	title := strings.Join(titlePieces, "")
	descPieces := []string{acct.Vanity, " uses Narro to create a podcast of articles transcribed to audio."}
	desc := strings.Join(descPieces, "")
	link := []string{"https://on.narro.co/", acct.Vanity}
	image := "https://www.narro.co/images/narro-icon-lg.png"
	if acct.Image.Valid {
		if acct.Image.String != "" {
			image = acct.Image.String
		}
	}
	if acct.PodcastTitle.Valid {
		if acct.PodcastTitle.String != "" {
			title = acct.PodcastTitle.String
		}
	}
	if acct.PodcastDescription.Valid {
		if acct.PodcastDescription.String != "" {
			desc = acct.PodcastDescription.String
		}
	}
	c := gopod.ChannelFactory(title, strings.Join(link, ""), desc, image)

	c.SetTTL("20") // in minutes
	c.SetPubDate(time.Now().UTC().Format(time.RFC1123))
	c.SetLanguage("en")
	c.SetWebMaster("josh@narro.co")
	c.SetManagingEditor(acct.Email)
	c.SetCopyright("All article content copyright of respective source authors.")
	c.SetiTunesExplicit("no")
	c.SetiTunesAuthor(acct.Vanity)
	TunesCategories := ""
	if acct.TunesCategories.Valid {
		TunesCategories = acct.TunesCategories.String
	}
	c.SetCategory(TunesCategories)
	c.SetiTunesSubtitle(desc)
	c.SetiTunesSummary(desc)
	c.SetiTunesOwner(acct.Vanity, "josh@narro.co")

	for iter.Next() {
		result := Article{}
		if err := iter.StructScan(&result); err != nil {
			log.Printf("error scanning article struct")
			log.Fatal(err)
		}
		linkList := listLinks(result.Links)
		resultLink := []string{"https://www.narro.co/article/", result.Id}
		enclosureLink := []string{"https://www.narro.co/article/", result.Id, ".mp3"}
		url := ""
		if result.Url.Valid {
			urlPieces := []string{"<a href=\"", result.Url.String, "\">", result.Url.String, "</a>"}
			url = strings.Join(urlPieces, "")
		}
		resultDesc := []string{result.Description, url, linkList}
		i := &gopod.Item{
			Title:         result.Title,
			Link:          strings.Join(resultLink, ""),
			Description:   strings.Join(resultDesc, "<br/> "),
			PubDate:       result.Created.UTC().Format(time.RFC1123),
			Author:        acct.Email,
			Guid:          strings.Join(resultLink, ""),
			TunesDuration: strconv.FormatFloat(result.Mp3Length, 'f', 0, 64),
			TunesAuthor:   acct.Email,
			TunesSubtitle: result.Title,
			TunesSummary:  strings.Join(resultDesc, " ... "),
			TunesExplicit: "no",
		}
		i.SetEnclosure(strings.Join(enclosureLink, ""), strconv.FormatFloat(result.Mp3Length, 'f', 0, 64), "audio/mpeg")
		c.AddItem(i)
	}
	if err := iter.Close(); err != nil {
		log.Printf(err.Error())
	}

	return c
}

func listLinks(LinksJSON types.JSONText) string {
	links := make([]ArticleLink, 0)
	LinksJSON.Unmarshal(&links)
	results := make([]string, len(links)+2)
	results = append(results, "<ul class=\"linkList\">")
	for _, r := range links {
		link := "<li><a href=\"first\">second</a></li>"
		link = strings.Replace(link, "first", r.Href, 1)
		link = strings.Replace(link, "second", r.Text, 1)
		results = append(results, link)
	}
	results = append(results, "</ul>")
	return strings.Join(results, "")
}
