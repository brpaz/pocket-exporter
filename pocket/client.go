package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/phayes/freeport"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// APIURL The base url for the Pocket API.
const APIURL = "https://getpocket.com/v3"

// AuthorizeURL The endpoint to redirect the user to authorize the application.
const AuthorizeURL = "https://getpocket.com/auth/authorize?request_token=%request_token%&redirect_uri=%redirect_uri%"

// AuthenticationError Error code when authenticating on Pocket API
const AuthenticationError = "AUTHENTICATION_ERROR"

// SystemError Error code for generic error from Pocket API
const SystemError = "SYSTEM_ERROR"

// Client Client struct
type Client struct {
	ConsumerKey string
	AccessToken string
	HTTPClient  *http.Client
}

// AuthRequestResponse Response for the "Auth Request"
type AuthRequestResponse struct {
	Code string `json:"code"`
}

// AuthorizeRequestResponse Model of the "Authorize request response"
type AuthorizeRequestResponse struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

// ArticleItem ArticleItem Model
type ArticleItem struct {
	ItemID     string `json:"item_id"`
	GivenURL   string `json:"given_url"`
	GivenTitle string `json:"given_title"`
	Favorite   string `json:"favorite"`
	Status     string `json:"status"`
}

// RetrieveResponse Retrieve response result
type RetrieveResponse struct {
	Status int                    `json:"status"`
	List   map[string]ArticleItem `json:"list"`
}

// APIError Api error struct
type APIError struct {
	Code    string
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Pocket API Error %s: %s", e.Code, e.Message)
}

// NewClient Create a new Pocket client
func NewClient(consumerKey string, client *http.Client) *Client {
	return &Client{
		HTTPClient:  client,
		ConsumerKey: consumerKey,
	}
}

// make a request to Pocket API
func (c *Client) doRequest(uri string, params map[string]string, returnStruct interface{}) (interface{}, error) {
	url := APIURL + uri

	jsonParams, _ := json.Marshal(params)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case 200:
		err = json.Unmarshal(body, &returnStruct)

		if err != nil {
			return "", err
		}

		return returnStruct, nil

	case 400:
		return returnStruct, errors.New("ERROR: Bad request:" + string(body))
	case 401:
		return returnStruct, errors.New("ERROR: It was not possible to authenticate on the Pocket API")
	default:
		return returnStruct, errors.New("ERROR: Unkwon error wile connecitng to the Pocket API")
	}
}

func (c *Client) getRequestToken(redirectURL string) (string, error) {

	var response AuthRequestResponse

	params := map[string]string{"consumer_key": c.ConsumerKey, "redirect_uri": redirectURL}
	resp, err := c.doRequest("/oauth/request", params, &response)

	if err != nil {
		return "", err
	}

	return resp.(*AuthRequestResponse).Code, nil
}

func (c *Client) getAccessToken(code string) (string, error) {

	var response AuthorizeRequestResponse

	params := map[string]string{"consumer_key": c.ConsumerKey, "code": code}
	resp, err := c.doRequest("/oauth/authorize", params, &response)

	if err != nil {
		return "", err
	}

	return resp.(*AuthorizeRequestResponse).AccessToken, nil
}

// Authenticate Initliazes the Oauth Authentication flow and store the access token obtained from token
func (c *Client) Authenticate() (string, error) {
	ch := make(chan struct{})

	port, _ := freeport.GetFreePort()
	srv := &http.Server{
		Addr: ":" + strconv.Itoa(port),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Authorized. You can close this page.")
		ch <- struct{}{}
	})

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// NOTE: there is a chance that next line won't have time to run,
			// as main() doesn't wait for this goroutine to stop. don't use
			// code with race conditions like these for production. see post
			// comments below on more discussion on how to handle this.
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	defer srv.Close()

	redirectURL := "http://localhost:" + strconv.Itoa(port)

	code, err := c.getRequestToken(redirectURL)

	if err != nil {
		return "", err
	}

	authorizeURL := strings.Replace(AuthorizeURL, "%request_token%", code, 1)
	authorizeURL = strings.Replace(authorizeURL, "%redirect_uri%", redirectURL, 1)
	fmt.Println("Please open the following url in your browser to authorize the application:")
	fmt.Println(authorizeURL)

	<-ch

	accessToken, err := c.getAccessToken(code)

	if err != nil {
		return "", err
	}

	c.AccessToken = accessToken

	return accessToken, nil
}

// Retrieve Gets articles saved in Pocket
func (c *Client) Retrieve() ([]MappedArticle, error) {

	var response RetrieveResponse
	var listItems []MappedArticle

	params := map[string]string{
		"consumer_key": c.ConsumerKey,
		"access_token": c.AccessToken,
		"state":        "all",
		"sort":         "newest",
		"detailType":   "complete",
	}

	resp, err := c.doRequest("/get", params, &response)

	if err != nil {
		return listItems, err
	}

	articles := resp.(*RetrieveResponse)

	for key, item := range articles.List {
		listItems = append(listItems, MappedArticle{
			ID:    key,
			Title: item.GivenTitle,
			URL:   item.GivenURL,
		})
	}

	return listItems, nil
}
