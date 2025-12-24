package exa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.exa.ai"
)

// Client is the Exa API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option defines a functional option for the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// New creates a new Exa API client.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Common request types

// TextOptions defines options for retrieving text content from pages.
type TextOptions struct {
	// IncludeHtmlTags specifies whether to include HTML tags in the returned text.
	IncludeHtmlTags bool `json:"includeHtmlTags,omitempty"`
	// MaxCharacters is the maximum number of characters to return for each result.
	MaxCharacters int `json:"maxCharacters,omitempty"`
}

// HighlightOptions defines options for generating highlights from search results.
type HighlightOptions struct {
	// NumSentences is the number of sentences per highlight.
	NumSentences int `json:"numSentences,omitempty"`
	// HighlightsPerURL is the maximum number of highlights to return per URL.
	HighlightsPerURL int `json:"highlightsPerUrl,omitempty"`
	// Query is a specific query to use for highlights (if different from search query).
	Query string `json:"query,omitempty"`
}

// SummaryOptions defines options for generating summaries.
type SummaryOptions struct {
	// Query is a specific query to use for generating the summary.
	Query string `json:"query,omitempty"`
}

// ContentsOptions defines options for what content to retrieve for each result.
type ContentsOptions struct {
	// Text retrieves the full text content of the result.
	Text *TextOptions `json:"text,omitempty"`
	// Highlights retrieves key excerpts from the text.
	Highlights *HighlightOptions `json:"highlights,omitempty"`
	// Summary retrieves a concise summary of the page.
	Summary *SummaryOptions `json:"summary,omitempty"`
	// LiveCrawl specifies the crawling behavior ("always", "fallback", "never").
	LiveCrawl string `json:"livecrawl,omitempty"`
}

// SearchOptions defines parameters for the Search endpoint.
type SearchOptions struct {
	// NumResults is the number of results to return (default: 10).
	NumResults int `json:"numResults,omitempty"`
	// IncludeDomains is a list of domains to include in the search.
	IncludeDomains []string `json:"includeDomains,omitempty"`
	// ExcludeDomains is a list of domains to exclude from the search.
	ExcludeDomains []string `json:"excludeDomains,omitempty"`
	// StartCrawlDate filters results crawled after this date (ISO 8601).
	StartCrawlDate string `json:"startCrawlDate,omitempty"`
	// EndCrawlDate filters results crawled before this date (ISO 8601).
	EndCrawlDate string `json:"endCrawlDate,omitempty"`
	// StartPublishedDate filters results published after this date (ISO 8601).
	StartPublishedDate string `json:"startPublishedDate,omitempty"`
	// EndPublishedDate filters results published before this date (ISO 8601).
	EndPublishedDate string `json:"endPublishedDate,omitempty"`
	// IncludeText is a list of strings that must be present in the result.
	IncludeText []string `json:"includeText,omitempty"`
	// ExcludeText is a list of strings that must not be present in the result.
	ExcludeText []string `json:"excludeText,omitempty"`
	// Contents specifies what page content to include in results.
	Contents *ContentsOptions `json:"contents,omitempty"`
	// UseAutoprompt enables automatic query enhancement.
	UseAutoprompt *bool `json:"useAutoprompt,omitempty"`
	// Type is the search type ("neural" or "keyword").
	Type string `json:"type,omitempty"`
	// Category is the search category (e.g., "company", "research paper").
	Category string `json:"category,omitempty"`
}

type searchRequest struct {
	Query string `json:"query"`
	SearchOptions
}

// Result represents a single search result.
type Result struct {
	// ID is the unique identifier for the result.
	ID string `json:"id"`
	// URL is the URL of the result.
	URL string `json:"url"`
	// Title is the title of the page.
	Title string `json:"title"`
	// Author is the author of the page content, if available.
	Author string `json:"author"`
	// PublishedDate is the date the page was published, if available.
	PublishedDate string `json:"publishedDate"`
	// Text is the retrieved text content, if requested.
	Text string `json:"text"`
	// Highlights are the retrieved text highlights, if requested.
	Highlights []string `json:"highlights"`
	// Summary is the generated summary, if requested.
	Summary string `json:"summary"`
	// Score is the relevance score of the result.
	Score float64 `json:"score"`
	// Image is a representative image URL for the page.
	Image string `json:"image,omitempty"`
	// Favicon is the URL of the page's favicon.
	Favicon string `json:"favicon,omitempty"`
}

// SearchResponse represents the response from the Search endpoint.
type SearchResponse struct {
	// Results is the list of search results.
	Results []Result `json:"results"`
	// RequestID is the unique identifier for the API request.
	RequestID string `json:"requestId"`
}

// Search performs a search using the Exa API.
// It uses neural or keyword search to find the most relevant results for the query.
func (c *Client) Search(ctx context.Context, query string, opts SearchOptions) (*SearchResponse, error) {
	reqBody := searchRequest{
		Query:         query,
		SearchOptions: opts,
	}
	var resp SearchResponse
	if err := c.do(ctx, "POST", "/search", reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindSimilarOptions defines parameters for the FindSimilar endpoint.
type FindSimilarOptions struct {
	// NumResults is the number of results to return (default: 10).
	NumResults int `json:"numResults,omitempty"`
	// IncludeDomains is a list of domains to include.
	IncludeDomains []string `json:"includeDomains,omitempty"`
	// ExcludeDomains is a list of domains to exclude.
	ExcludeDomains []string `json:"excludeDomains,omitempty"`
	// StartCrawlDate filters results crawled after this date.
	StartCrawlDate string `json:"startCrawlDate,omitempty"`
	// EndCrawlDate filters results crawled before this date.
	EndCrawlDate string `json:"endCrawlDate,omitempty"`
	// StartPublishedDate filters results published after this date.
	StartPublishedDate string `json:"startPublishedDate,omitempty"`
	// EndPublishedDate filters results published before this date.
	EndPublishedDate string `json:"endPublishedDate,omitempty"`
	// Contents specifies what page content to include in results.
	Contents *ContentsOptions `json:"contents,omitempty"`
}

type findSimilarRequest struct {
	URL string `json:"url"`
	FindSimilarOptions
}

// FindSimilar finds similar links to the provided URL.
// It returns pages that are similar in meaning to the input URL.
func (c *Client) FindSimilar(ctx context.Context, url string, opts FindSimilarOptions) (*SearchResponse, error) {
	reqBody := findSimilarRequest{
		URL:                url,
		FindSimilarOptions: opts,
	}
	var resp SearchResponse
	if err := c.do(ctx, "POST", "/findSimilar", reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetContentsOptions defines parameters for the Contents endpoint.
type GetContentsOptions struct {
	// Text retrieves the full text content of the result.
	Text *TextOptions `json:"text,omitempty"`
	// Highlights retrieves key excerpts from the text.
	Highlights *HighlightOptions `json:"highlights,omitempty"`
	// Summary retrieves a concise summary of the page.
	Summary *SummaryOptions `json:"summary,omitempty"`
	// LiveCrawl specifies the crawling behavior ("always", "fallback", "never").
	LiveCrawl string `json:"livecrawl,omitempty"`
}

type getContentsRequest struct {
	IDs []string `json:"ids"`
	GetContentsOptions
}

// GetContentsResponse represents the response from the Contents endpoint.
type GetContentsResponse struct {
	// Results is the list of retrieved contents.
	Results []Result `json:"results"`
	// RequestID is the unique identifier for the API request.
	RequestID string `json:"requestId"`
}

// GetContents retrieves contents for the provided IDs (URLs).
// It returns clean, parsed content from the specified web pages.
func (c *Client) GetContents(ctx context.Context, ids []string, opts GetContentsOptions) (*GetContentsResponse, error) {
	reqBody := getContentsRequest{
		IDs:                ids,
		GetContentsOptions: opts,
	}
	var resp GetContentsResponse
	if err := c.do(ctx, "POST", "/contents", reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AnswerOptions defines parameters for the Answer endpoint.
type AnswerOptions struct {
	// Query is the question to answer.
	Query string `json:"query"`
	// Text specifies whether to return full text for citations.
	Text bool `json:"text,omitempty"`
}

type answerResponse struct {
	Answer    string   `json:"answer"`
	Citations []Result `json:"citations"`
	RequestID string   `json:"requestId"`
}

// AnswerResponse represents the response from the Answer endpoint.
type AnswerResponse struct {
	Answer    string
	Citations []Result
	RequestID string
}

// ContextOptions defines parameters for the Context (Exa Code) endpoint.
type ContextOptions struct {
	// Query is the search query to find relevant code snippets.
	Query string `json:"query"`
	// TokensNum is the token limit for the response. Can be "dynamic" or an integer.
	TokensNum interface{} `json:"tokensNum,omitempty"`
}

// ContextResponse represents the response from the Context endpoint.
type ContextResponse struct {
	// Response is the formatted code snippets and contextual examples.
	Response string `json:"response"`
	// RequestID is the unique identifier for the API request.
	RequestID string `json:"requestId"`
	// ResultsCount is the number of results found.
	ResultsCount int `json:"resultsCount"`
	// SearchTime is the time taken for the search in seconds.
	SearchTime float64 `json:"searchTime"`
	// OutputTokens is the number of tokens in the response.
	OutputTokens int `json:"outputTokens"`
	// CostDollars is the cost of the request.
	CostDollars interface{} `json:"costDollars"`
}

// Answer gets an answer to a question.
// It uses an LLM informed by Exa search results to generate a grounded answer.
func (c *Client) Answer(ctx context.Context, query string) (*AnswerResponse, error) {
	reqBody := AnswerOptions{
		Query: query,
		Text:  true, // Default to true as per common usage
	}
	var resp answerResponse
	if err := c.do(ctx, "POST", "/answer", reqBody, &resp); err != nil {
		return nil, err
	}
	return &AnswerResponse{
		Answer:    resp.Answer,
		Citations: resp.Citations,
		RequestID: resp.RequestID,
	}, nil
}

// Context gets code context for a query (Exa Code).
// It searches billions of repos and docs to find practical code examples.
// tokensNum can be "dynamic" or a specific integer (e.g., 5000).
func (c *Client) Context(ctx context.Context, query string, tokensNum interface{}) (*ContextResponse, error) {
	if tokensNum == nil {
		tokensNum = "dynamic"
	}
	reqBody := ContextOptions{
		Query:     query,
		TokensNum: tokensNum,
	}
	var resp ContextResponse
	if err := c.do(ctx, "POST", "/context", reqBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// do executes the HTTP request.
func (c *Client) do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("User-Agent", "exa-go-client/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("api request failed with status %d", resp.StatusCode)
		}
		return fmt.Errorf("api error: %s (status %d)", errResp.Message, resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
