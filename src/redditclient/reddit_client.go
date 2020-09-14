package redditclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	redditAPIAccessTokenEndpoint = "https://www.reddit.com/api/v1/access_token"
	redditBaseAPIEndpoint        = "https://oauth.reddit.com"
)

type Config struct {
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	UserAgent    string
}

type AuthToken struct {
	ExpirationTimestamp int64
	Token               string
}

type Client struct {
	Config    Config
	AuthToken *AuthToken
}

type Submission struct {
	ID          string
	URL         string
	FallbackURL string
	Title       string
	Permalink   string
	Author      string
	IsSelf      bool
	CreatedUTC  int64
}

type Comment struct {
	ID         string
	IsStickied bool
	Body       string
	Author     string
	Permalink  string
	CreatedUTC int64
}

func (authToken *AuthToken) isReadyToExpire() bool {
	return time.Now().Add(time.Minute).Unix() > authToken.ExpirationTimestamp
}

func (client *Client) apiGetRequest(endpoint string) ([]byte, error) {
	if client.AuthToken.isReadyToExpire() {
		err := client.renewAuthToken()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(http.MethodGet, redditBaseAPIEndpoint+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", client.AuthToken.Token))
	req.Header.Add("User-Agent", client.Config.UserAgent)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *Client) GetLatestSoccerSubmissions() ([]Submission, error) {
	response, err := client.apiGetRequest("/r/soccer/new.json?limit=100")
	if err != nil {
		return nil, err
	}
	submissionsJSON := gjson.GetBytes(response, "data.children.#.data").Array()
	submissions := make([]Submission, len(submissionsJSON))
	for idx, submissionJSON := range submissionsJSON {
		submissionJSONMap := submissionJSON.Map()

		fallbackURLResult := submissionJSON.Get("media.reddit_video.fallback_url")
		fallbackURL := ""
		if fallbackURLResult.Exists() {
			fallbackURL = fallbackURLResult.String()
		}

		submissions[idx] = Submission{
			submissionJSONMap["id"].String(),
			submissionJSONMap["url"].String(),
			fallbackURL,
			submissionJSONMap["title"].String(),
			submissionJSONMap["permalink"].String(),
			submissionJSONMap["author"].String(),
			submissionJSONMap["is_self"].Bool(),
			int64(submissionJSONMap["created_utc"].Float()),
		}
	}
	return submissions, nil
}

func (client *Client) getComments(permalink, path string) ([]Comment, error) {
	response, err := client.apiGetRequest(permalink)
	if err != nil {
		return nil, err
	}

	commentsJSON := gjson.GetBytes(response, path).Array()
	comments := []Comment{}
	for _, commentJSON := range commentsJSON {
		commentJSONMap := commentJSON.Map()
		comments = append(comments, Comment{
			commentJSONMap["id"].String(),
			commentJSONMap["stickied"].Bool(),
			commentJSONMap["body"].String(),
			commentJSONMap["author"].String(),
			commentJSONMap["permalink"].String(),
			commentJSONMap["created_utc"].Int(),
		})
	}

	return comments, err
}

func (client *Client) GetStickiedComment(permalink string) (*Comment, error) {
	topLevelComments, err := client.getComments(permalink, "1.data.children.#.data")
	if err != nil {
		return nil, err
	}

	for _, comment := range topLevelComments {
		if comment.IsStickied {
			return &comment, nil
		}
	}
	return nil, nil
}

func (client *Client) GetCommentReplies(permalink string) ([]Comment, error) {
	return client.getComments(permalink, "1.data.children.0.data.replies.data.children.#.data")
}

func newAuthToken(config Config) (*AuthToken, error) {
	formData := fmt.Sprintf("grant_type=password&username=%s&password=%s", config.Username, config.Password)
	body := bytes.NewBufferString(formData)

	req, err := http.NewRequest(http.MethodPost, redditAPIAccessTokenEndpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", config.UserAgent)
	req.SetBasicAuth(config.ClientID, config.ClientSecret)

	authTime := time.Now().Unix()
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responseJSON := gjson.ParseBytes(responseBody).Map()
	token, expiresIn := responseJSON["access_token"].String(), responseJSON["expires_in"].Int()

	expirationTimestamp := authTime + expiresIn
	return &AuthToken{expirationTimestamp, token}, nil
}

func (client *Client) renewAuthToken() error {
	authToken, err := newAuthToken(client.Config)
	if err != nil {
		log.Warn("Failed to renew Reddit auth token", err)
		return err
	}
	client.AuthToken = authToken
	return nil
}

func NewRedditClient(config Config) *Client {
	client := &Client{Config: config}
	err := client.renewAuthToken()
	if err != nil {
		log.Fatal("Failed to create the Reddit client", err)
	}
	return client
}
