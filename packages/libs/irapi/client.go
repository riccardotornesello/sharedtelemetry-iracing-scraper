package irapi

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"
)

type IRacingApiClient struct {
	client     *http.Client
	retryAfter time.Time
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
		Jar: jar,
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
	var resp *http.Response
	var err error

	for {
		if c.retryAfter.After(time.Now()) {
			slog.Info(fmt.Sprintf("Rate limit exceeded, waiting until %v", c.retryAfter.Format(time.RFC3339)))
			time.Sleep(time.Until(c.retryAfter))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // TODO: allow timeout customization
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://members-ng.iracing.com"+path, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request for %s: %w", path, err)
		}

		resp, err = c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error getting %s: %w", path, err)
		}

		if resp.StatusCode != 429 {
			break
		}

		slog.Info(fmt.Sprintf("Rate limit exceeded for %s, retrying in a bit", path))

		// TODO: allow to skip retrying
		// TODO: allow max retry count
		rateLimitReset := resp.Header.Get("X-RateLimit-Reset")
		if rateLimitReset == "" {
			break
		}

		rateLimitResetInt, err := strconv.ParseInt(rateLimitReset, 10, 64)
		if err != nil {
			break
		}

		// Not atomic, but we don't care
		c.retryAfter = time.Unix(rateLimitResetInt, 0).Add(2 * time.Second)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error getting %s: %s", path, resp.Status)
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
