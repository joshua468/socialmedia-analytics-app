package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/huandu/facebook"
)

type SocialMediaClient struct {
	APIToken       string
	APIEndpoint    string
	TwitterClient  *twitter.Client
	FacebookClient *facebook.Session
}

type Post struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Likes     int       `json:"likes"`
	Shares    int       `json:"shares"`
	TimeStamp time.Time `json:"timestamp"`
}

type AnalyticsData struct {
	Posts        []Post         `json:"posts"`
	Followers    int            `json:"followers"`
	HashtagUsage map[string]int `json:"hashtag_usage"`
}

func main() {
	twitterClient := initTwitterClient()
	facebookClient := initFacebookClient()

	client := &SocialMediaClient{
		APIToken:       "api_key",
		APIEndpoint:    "api",
		TwitterClient:  twitterClient,
		FacebookClient: facebookClient,
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", client.dashboardHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("server listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func initTwitterClient() *twitter.Client {
	config := oauth1.NewConfig("jRFbfnEhMQgdXF9dTtmHVmFlV", "RxchVsPlv7LPlKHuk88CD6IgSateVVKI2glwdvLS5tSOpJF51O")
	token := oauth1.NewToken("1533464023064268800-qWYBJHxRlsD3bICvWCEwQfwa4SSwkK", "wq6L2MsGieCWJJLMMDHuB04VfOKnnGOjykh07nZglXOAh")
	httpClient := config.Client(oauth1.NoContext, token)
	return twitter.NewClient(httpClient)
}
func initFacebookClient() *facebook.Session {
	app := facebook.New("912899540328104", "912899540328104")
	session := app.Session("41a9426f921e5221315a3db5fe7b5d80")
	return session
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (client *SocialMediaClient) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	analyticsData := client.fetchAnalyticsData()

	renderTemplate(w, "dashboard.html", analyticsData)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (client *SocialMediaClient) fetchAnalyticsData() *AnalyticsData {
	return &AnalyticsData{
		Posts: []Post{
			{ID: "1", Message: "Hello World!", Likes: 10, Shares: 5, TimeStamp: time.Now()},
			{ID: "2", Message: "Good Morning", Likes: 15, Shares: 7, TimeStamp: time.Now()},
		},
		Followers:    1000,
		HashtagUsage: map[string]int{"#golang": 20, "#analytics": 15, "#social": 10},
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl = fmt.Sprintf("templates/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
