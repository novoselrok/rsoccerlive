package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/novoselrok/rsoccerlive/src/models"
	"github.com/novoselrok/rsoccerlive/src/redditclient"
	log "github.com/sirupsen/logrus"
)

const (
	firstTeamScoreRegexPattern  = "[\\[\\(]?([^0-9][0-9]{1,2})[\\]\\)]?"
	secondTeamScoreRegexPattern = "[\\[\\(]?([^0-9]?[0-9]{1,2})[\\]\\)]?"

	autoModeratorAuthor  = "AutoModerator"
	mirrorsCommentPrefix = "Mirrors / Alternate angles"
)

var supportedHighlightHosts = []string{
	"streamable.com",
	"streamja.com",
	"clippituser.tv",
}

var highlightTitleRegexps []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(fmt.Sprintf("%s\\s*-\\s*%s", firstTeamScoreRegexPattern, secondTeamScoreRegexPattern)),
	regexp.MustCompile(fmt.Sprintf("%s\\s*[Vv][Ss]\\.?\\s*%s", firstTeamScoreRegexPattern, secondTeamScoreRegexPattern)),
	regexp.MustCompile("goals?"),
	regexp.MustCompile("penalty"),
	regexp.MustCompile("(red|yellow)\\s+card"),
}

var redditMarkupUrlRegexps []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile("\\[.*?\\]\\s*\\((https?://\\S+?)\\)"),
	regexp.MustCompile("(?:[^*(]|^)(https?://\\S+?)(?:\\s|\\)|$)"),
}

func isHighlightTitle(title string) bool {
	titleLower := strings.ToLower(title)
	for _, highlightTitleRegexp := range highlightTitleRegexps {
		if highlightTitleRegexp.FindStringSubmatch(titleLower) != nil {
			return true
		}
	}
	return false
}

func isSupportedHighlightURL(highlightURL string) bool {
	parsedURL, err := url.Parse(highlightURL)
	if err != nil {
		return false
	}

	highlightHost := parsedURL.Host
	if strings.HasPrefix(highlightHost, "www.") {
		highlightHost = highlightHost[4:]
	}

	for _, supportedHighlightHost := range supportedHighlightHosts {
		if supportedHighlightHost == highlightHost {
			return true
		}
	}
	return false
}

func isHighlightSubmission(submission redditclient.Submission) bool {
	return !submission.IsSelf && isHighlightTitle(submission.Title) && isSupportedHighlightURL(submission.URL)
}

func GetLatestHighlightSubmissions(redditClient *redditclient.Client) []redditclient.Submission {
	submissions, err := redditClient.GetLatestSoccerSubmissions()
	if err != nil {
		log.Error("Failed retrieving latest soccer submissions")
		return []redditclient.Submission{}
	}

	highlightSubmissions := []redditclient.Submission{}
	for _, submission := range submissions {
		if isHighlightSubmission(submission) {
			highlightSubmissions = append(highlightSubmissions, submission)
		}
	}

	return highlightSubmissions
}

func ConvertSubmissionsToHighlights(submissions []redditclient.Submission) []models.Highlight {
	highlights := make([]models.Highlight, len(submissions))
	for idx, submission := range submissions {
		highlights[idx] = models.Highlight{
			URL:                submission.URL,
			Title:              submission.Title,
			RedditSubmissionID: submission.ID,
			RedditPermalink:    submission.Permalink,
			RedditAuthor:       submission.Author,
			RedditCreatedAt:    time.Unix(submission.CreatedUTC, 0),
		}
	}
	return highlights
}

func GetMirrorsCommentThreadReplies(redditClient *redditclient.Client, permalink string) ([]redditclient.Comment, error) {
	comment, err := redditClient.GetStickiedComment(permalink)
	if comment == nil {
		return nil, err
	}

	isMirrorsCommentThread := comment.Author == autoModeratorAuthor && strings.Contains(comment.Body, mirrorsCommentPrefix)
	if !isMirrorsCommentThread {
		return nil, nil
	}

	replies, err := redditClient.GetCommentReplies(comment.Permalink)
	if err != nil {
		return nil, err
	}
	return replies, nil
}

func extractHighlightUrls(body string) []string {
	urls := []string{}
	for _, regexp := range redditMarkupUrlRegexps {
		matches := regexp.FindAllStringSubmatch(body, -1)
		if matches == nil {
			continue
		}
		for _, match := range matches {
			urls = append(urls, match[1])
		}
	}

	highlightUrls := []string{}
	for _, url := range urls {
		if isSupportedHighlightURL(url) {
			highlightUrls = append(highlightUrls, url)
		}
	}
	return highlightUrls
}

func GetHighlightMirrors(redditClient *redditclient.Client, highlight models.Highlight) ([]models.HighlightMirror, error) {
	replies, err := GetMirrorsCommentThreadReplies(redditClient, highlight.RedditPermalink)
	if err != nil {
		return nil, err
	}

	highlightMirrors := []models.HighlightMirror{}
	for _, reply := range replies {
		highlightUrls := extractHighlightUrls(reply.Body)

		for _, highlightURL := range highlightUrls {
			highlightMirrors = append(highlightMirrors, models.HighlightMirror{
				HighlightID:     highlight.ID,
				URL:             highlightURL,
				RedditPermalink: reply.Permalink,
				RedditAuthor:    reply.Author,
				RedditCreatedAt: time.Unix(reply.CreatedUTC, 0),
			})
		}
	}
	return highlightMirrors, nil
}
