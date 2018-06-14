package titles

import (
	"net/http"
	"html/template"
	"time"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Page struct {
	Title string
	IsRowLight bool
	Suggestions []*TitleSuggestionDisplay
}

type TitleSuggestionDisplay struct{
	ID int
	Title string
	Author string
	TimeAgo string
	TotalVotes int
	IsRowLight bool
}

const tmplRoot = "web/templates/plugins/titles/"
var templates = template.Must(
	template.ParseFiles(
		tmplRoot + "form.html", tmplRoot + "results.html"))

var voteValidPath = regexp.MustCompile("^/(titles)/(vote)/([0-9]+)")
var ipsThatVoted = make(map[string]bool)

func WebInit() {
	http.HandleFunc("/titles", mainHandler)
	http.HandleFunc("/titles/vote/", voteHandler)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var tmpl string
	ip := getIpAddr(r)
	p := &Page{
		Title: showTitle,
		IsRowLight: false,
		Suggestions: getDisplayTitles(),
	}

	if ipsThatVoted[ip] {
		tmpl = "results"
	} else {
		tmpl = "form"
	}

	renderTemplate(w, tmpl, p)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIpAddr(r)
	println("IP: " + ip)
	if ipsThatVoted[ip] {
		http.Redirect(w, r, "/titles", http.StatusFound)
		return
	}

	m := voteValidPath.FindStringSubmatch(r.URL.Path)
	id, _ := strconv.Atoi(m[2]);
	addVoteToTitle(id);
	ipsThatVoted[ip] = true
	http.Redirect(w, r, "/titles", http.StatusFound)
}

func getIpAddr(r *http.Request) string {
	ipForwards := r.Header.Get("x-forwarded-for")
	if len(ipForwards) > 0 {
		ips := strings.Split(ipForwards, ", ")
		return ips[0]
	}
	return r.RemoteAddr
}

func getDisplayTitles() (ret []*TitleSuggestionDisplay) {
	now := time.Now()
	suggestions := GetTitles()
	for i := 0; i < len(suggestions); i++ {
		ret = append(ret, titleSuggestionToTitleSuggestionDisplay(suggestions[i], now))
	}
	return
}

func renderTemplate(w http.ResponseWriter, tmpl string, page *Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", page)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func titleSuggestionToTitleSuggestionDisplay(
	sug *TitleSuggestion, now time.Time) (*TitleSuggestionDisplay) {
	return &TitleSuggestionDisplay{
		ID: sug.ID,
		Title: sug.Title,
		Author: sug.Author,
		TimeAgo: calcTimeAgo(now, sug.CreatedOn),
		TotalVotes: sug.TotalVotes,
		IsRowLight: isRowLight(),
	}
}

func calcTimeAgo(now time.Time, past time.Time) (string) {
	secondsAgo := now.Unix() - past.Unix()
	if (60 > secondsAgo) {
		return fmt.Sprintf("%d second(s) ago.", secondsAgo)
	}

	minutesAgo := secondsAgo / 60
	if (60 > minutesAgo) {
		return fmt.Sprintf("%d minute(s) ago.", minutesAgo)
	}

	hoursAgo := minutesAgo / 60
	if (24 > hoursAgo) {
		return fmt.Sprint("%d hour(s) ago.", hoursAgo)
	}

	daysAgo := hoursAgo / 24
	return fmt.Sprint("%d day(s) ago.", daysAgo)
}

var rowLight = false
func isRowLight() (bool){
	rowLight = !rowLight
	return rowLight
}