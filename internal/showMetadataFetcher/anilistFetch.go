package showmetadatafetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var query string = `query SearchAnime($search: String!) {
  Page(page: 1, perPage: 10) {
    media(search: $search, type: ANIME) {
      id
      title {
        romaji
        english
      }
      description(asHtml: false)
      episodes
      duration
      season
      seasonYear
      format
      status
      genres
      coverImage {
        large
        extraLarge
      }
    }
  }
}`

type payload struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type perShow struct {
	Id    int `json:"id"`
	Title struct {
		Romanji string `json:"romanji"`
		English string `json:"english"`
	} `json:"title"`
	Description string   `json:"description"`
	Episodes    int      `json:"episodes"`
	Duration    int      `json:"duration"`
	Season      string   `json:"season"`
	SeasonYear  int      `json:"seasonYear"`
	Format      string   `json:"format"`
	Status      string   `json:"status"`
	Genres      []string `json:"genres"`
	CoverImage  struct {
		Large      string `json:"large"`
		ExtraLarge string `json:"extraLarge"`
	} `json:"coverImage"`
}

type aniListResponse struct {
	Data struct {
		Page struct {
			Media []perShow `json:"media"`
		} `json:"Page"`
	} `json:"data"`

	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type rateAnyListLimiting struct {
}

func callAnilist(queryString string) (aniListResponse, error) {
	graphQL_URL := os.Getenv("ANILIST_ENDPOINT")
	var output aniListResponse
	payload := payload{Query: query, Variables: map[string]any{
		"search": queryString,
	}}

	body, err := json.Marshal(payload)
	if err != nil {
		return aniListResponse{}, err
	}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		graphQL_URL,
		bytes.NewReader(body),
	)

	if err != nil {
		return aniListResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return aniListResponse{}, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return aniListResponse{}, err
	}

	if len(output.Errors) > 0 {
		return aniListResponse{}, errors.New(output.Errors[0].Message)
	}

	return output, nil

}

func getPerShow(data aniListResponse) []perShow {
	return data.Data.Page.Media
}
