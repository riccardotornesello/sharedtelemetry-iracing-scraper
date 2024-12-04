package client

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

type IRacingError struct {
	Message string `json:"message"`
}

func (e *IRacingError) Error() string {
	return e.Message
}

type IRacingApiClient struct {
	client *http.Client
}

type IRacingAuthResponse struct {
	Authcode string `json:"authcode"`
}

type IRacingResponse struct {
	Link string `json:"link"`
}

type IRacingChunkInfo struct {
	ChunkSize       int      `json:"chunk_size"`
	NumChunks       int      `json:"num_chunks"`
	Rows            int      `json:"rows"`
	BaseDownloadUrl string   `json:"base_download_url"`
	ChunkFileNames  []string `json:"chunk_file_names"`
}

func NewIRacingApiClient(email string, password string) *IRacingApiClient {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second, // TODO: make this configurable
	}

	tokenIn := []byte(password + strings.ToLower(email))
	hasher := sha256.New()
	hasher.Write(tokenIn)
	tokenHash := hasher.Sum(nil)
	tokenB64 := base64.StdEncoding.EncodeToString(tokenHash)

	resp, err := client.Post("https://members-ng.iracing.com/auth", "application/json", strings.NewReader(`{"email":"`+email+`","password":"`+tokenB64+`"}`))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Login failed")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	authResponse := &IRacingAuthResponse{}
	err = json.Unmarshal(body, authResponse)
	if err != nil {
		log.Fatal("Login failed")
	}

	return &IRacingApiClient{
		client: client,
	}
}

func (c *IRacingApiClient) Get(path string) ([]byte, error) {
	resp, err := c.client.Get("https://members-ng.iracing.com" + path)
	if err != nil {
		return nil, err
	}

	// Retry once more if the request failed
	// TODO: allow for more retries
	// TODO: allow to skip retrying
	if resp.StatusCode == 429 {
		rateLimitReset := resp.Header.Get("X-RateLimit-Reset")
		if rateLimitReset == "" {
			return nil, &IRacingError{Message: "Rate limit exceeded. Can't find rate limit reset time"}
		}

		rateLimitResetInt, err := strconv.ParseInt(rateLimitReset, 10, 64)
		if err != nil {
			return nil, err
		}

		time.Sleep(time.Until(time.Unix(rateLimitResetInt, 0)))

		resp, err = c.client.Get("https://members-ng.iracing.com" + path)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, &IRacingError{Message: string(body)}
	}

	response := &IRacingResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	resp, err = c.client.Get(response.Link)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *IRacingApiClient) GetChunks(chunkInfo *IRacingChunkInfo) [][]byte {
	out := make([][]byte, 0)

	for _, chunkFileName := range chunkInfo.ChunkFileNames {
		resp, err := c.client.Get(chunkInfo.BaseDownloadUrl + chunkFileName)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		out = append(out, body)
	}

	return out
}
