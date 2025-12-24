# exa-go

A robust, idiomatic Go client for the [Exa API](https://exa.ai).

## Installation

```bash
go get github.com/ahmedthabet/exa-go
```

## Usage

First, obtain an API key from the [Exa Dashboard](https://dashboard.exa.ai).

### Initialization

```go
package main

import (
	"context"
	"os"

	"github.com/ahmedthabet/exa-go"
)

func main() {
	apiKey := os.Getenv("EXA_API_KEY")
	client := exa.New(apiKey)

	// ...
}
```

### Search

Search the web using neural or keyword search.

```go
	resp, err := client.Search(context.Background(), "latest developments in AI", exa.SearchOptions{
		NumResults: 5,
		Type:       "neural",
		Contents: &exa.ContentsOptions{
			Text: &exa.TextOptions{
				MaxCharacters: 1000,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	for _, result := range resp.Results {
		fmt.Printf("Title: %s\nURL: %s\n\n", result.Title, result.URL)
	}
```

### Answer

Get a grounded answer to a question using Exa's search results.

```go
	answer, err := client.Answer(context.Background(), "Who won the super bowl in 2024?")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Answer)
	for _, citation := range answer.Citations {
		fmt.Printf("Source: %s (%s)\n", citation.Title, citation.URL)
	}
```

### Context (Exa Code)

Find practical code examples and snippets.

```go
	resp, err := client.Context(context.Background(), "React hooks for state management", "dynamic")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Response)
```

### Find Similar

Find similar pages to a given URL.

```go
	resp, err := client.FindSimilar(context.Background(), "https://example.com/interesting-article", exa.FindSimilarOptions{
		NumResults: 3,
	})
```

### Get Contents

Retrieve clean, parsed text from URLs.

```go
	resp, err := client.GetContents(context.Background(), []string{"https://example.com/article1"}, exa.GetContentsOptions{
		Text: &exa.TextOptions{
			IncludeHtmlTags: false,
		},
	})
```

## Configuration

You can customize the client with functional options:

```go
client := exa.New(apiKey, 
    exa.WithBaseURL("https://custom-proxy.com"),
    exa.WithHTTPClient(customHTTPClient),
)
```

## Engineering Principles

- **Not Over-engineered**: Simple, single-file implementation where appropriate.
- **Seasoned**: Uses idiomatic Go patterns (context, functional options, explicit error handling).
- **Tested**: Includes unit tests with HTTP mocking and an end-to-end example.
- **Robust**: Handles API errors and unmarshals complex responses correctly.

## License

MIT
