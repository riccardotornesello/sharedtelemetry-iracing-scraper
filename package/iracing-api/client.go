package irapi

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

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

func NewIRacingApiClient(email string, password string) (*IRacingApiClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 60 * time.Second, // TODO: make this configurable
	}

	tokenIn := []byte(password + strings.ToLower(email))
	hasher := sha256.New()
	hasher.Write(tokenIn)
	tokenHash := hasher.Sum(nil)
	tokenB64 := base64.StdEncoding.EncodeToString(tokenHash)

	resp, err := client.Post("https://members-ng.iracing.com/auth", "application/json", strings.NewReader(`{"email":"`+email+`","password":"`+tokenB64+`"}`))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error authenticating: %s", resp.Status)
	}

	defer resp.Body.Close()

	authResponse := &IRacingAuthResponse{}
	err = json.NewDecoder(resp.Body).Decode(authResponse)
	if err != nil {
		return nil, err
	}

	return &IRacingApiClient{
		client: client,
	}, nil
}

func (c *IRacingApiClient) get(path string) (io.ReadCloser, error) {
	resp, err := c.client.Get("https://members-ng.iracing.com" + path)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", path, err)
	}

	// Retry once more if the request failed
	// TODO: allow for more retries
	// TODO: allow to skip retrying
	if resp.StatusCode == 429 {
		log.Printf("Rate limit exceeded for %s, retrying in a bit", path)

		rateLimitReset := resp.Header.Get("X-RateLimit-Reset")
		if rateLimitReset == "" {
			return nil, fmt.Errorf("Rate limit exceeded but can't find rate limit reset time")
		}

		rateLimitResetInt, err := strconv.ParseInt(rateLimitReset, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing rate limit reset time: %w", err)
		}

		time.Sleep(time.Until(time.Unix(rateLimitResetInt, 0)) + 2*time.Second)

		resp, err = c.client.Get("https://members-ng.iracing.com" + path)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting %s: %s - %s", path, resp.Status)
	}

	defer resp.Body.Close()

	response := &IRacingResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	payloadResp, err := c.client.Get(response.Link)
	if err != nil {
		return nil, err
	}

	return payloadResp.Body, nil
}

func (c *IRacingApiClient) getChunks(chunkInfo *IRacingChunkInfo) ([]io.ReadCloser, error) {
	out := make([]io.ReadCloser, len(chunkInfo.ChunkFileNames))

	for i, chunkFileName := range chunkInfo.ChunkFileNames {
		resp, err := c.client.Get(chunkInfo.BaseDownloadUrl + chunkFileName)
		if err != nil {
			return nil, err
		}

		out[i] = resp.Body
	}

	return out, nil
}
