package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

type GitHubProfile struct {
	Login             string   `json:"login"`
	ID                int      `json:"id"`
	AvatarURL         string   `json:"avatar_url"`
	Name              string   `json:"name"`
	Company           string   `json:"company"`
	Blog              string   `json:"blog"`
	Location          string   `json:"location"`
	Email             string   `json:"email"`
	Bio               string   `json:"bio"`
	TwitterUsername   string   `json:"twitter_username"`
	PublicRepos       int      `json:"public_repos"`
	PublicGists       int      `json:"public_gists"`
	Followers         int      `json:"followers"`
	Following         int      `json:"following"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}

type Repository struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
	Fork        bool      `json:"fork"`
	HTMLURL     string    `json:"html_url"`
	Language    string    `json:"language"`
	Stargazers  int       `json:"stargazers_count"`
	Watchers    int       `json:"watchers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	PushedAt    string    `json:"pushed_at"`
	Topics      []string  `json:"topics"`
}

type ProfileData struct {
	Profile      GitHubProfile `json:"profile"`
	Repositories []Repository  `json:"repositories"`
	LastUpdated  time.Time     `json:"last_updated"`
}

var cache ProfileData

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", homeHandler)
	r.GET("/api/profile", profileHandler)
	r.GET("/api/repos", reposHandler)
	r.GET("/rss", rssHandler)

	// Initialize cache
	go initializeCache()

	fmt.Println("🚀 GitHub Profile Server started on :8080")
	fmt.Println("⚡ Anime Mode: ON - Levi energy activated!")
	log.Fatal(r.Run(":8080"))
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":     "Mohit Tiwari - GitHub Profile",
		"anime_mode": "ON",
		"quote":     "If you don't take risks, you can't create a future.",
		"character": "Monkey D. Luffy",
	})
}

func profileHandler(c *gin.Context) {
	if cache.Profile.Login != "" {
		c.JSON(http.StatusOK, cache.Profile)
		return
	}

	profile, err := fetchGitHubProfile("Mohit4289")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cache.Profile = profile
	c.JSON(http.StatusOK, profile)
}

func reposHandler(c *gin.Context) {
	if len(cache.Repositories) > 0 {
		c.JSON(http.StatusOK, cache.Repositories)
		return
	}

	repos, err := fetchRepositories("Mohit4289")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cache.Repositories = repos
	cache.LastUpdated = time.Now()
	c.JSON(http.StatusOK, repos)
}

func rssHandler(c *gin.Context) {
	repos, err := fetchRepositories("Mohit4289")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching repositories")
		return
	}

	feed := &feeds.Feed{
		Title:       "Mohit Tiwari's GitHub Repositories",
		Link:        &feeds.Link{Href: "https://github.com/Mohit4289"},
		Description: "Latest repositories from Mohit Tiwari's GitHub",
		Author:      &feeds.Author{Name: "Mohit Tiwari", Email: "srttiwari4289@gmail.com"},
		Created:     time.Now(),
	}

	for _, repo := range repos {
		if !repo.Private && !repo.Fork {
			item := &feeds.Item{
				Title:       repo.Name,
				Link:        &feeds.Link{Href: repo.HTMLURL},
				Description: repo.Description,
				Created:     parseTime(repo.CreatedAt),
			}
			feed.Items = append(feed.Items, item)
		}
	}

	rss, err := feed.ToRss()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating RSS")
		return
	}

	c.Header("Content-Type", "application/rss+xml")
	c.String(http.StatusOK, rss)
}

func fetchGitHubProfile(username string) (GitHubProfile, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://api.github.com/users/%s", username))
	if err != nil {
		return GitHubProfile{}, err
	}
	defer resp.Body.Close()

	var profile GitHubProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return GitHubProfile{}, err
	}

	return profile, nil
}

func fetchRepositories(username string) ([]Repository, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&sort=updated", username))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func initializeCache() {
	profile, err := fetchGitHubProfile("Mohit4289")
	if err == nil {
		cache.Profile = profile
	}

	repos, err := fetchRepositories("Mohit4289")
	if err == nil {
		cache.Repositories = repos
		cache.LastUpdated = time.Now()
	}
}

func parseTime(timeStr string) time.Time {
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t
	}
	return time.Now()
}
