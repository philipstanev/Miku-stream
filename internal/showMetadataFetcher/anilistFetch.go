package showmetadatafetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	Episodes    int32    `json:"episodes"`
	Duration    int32    `json:"duration"`
	Season      string   `json:"season"`
	SeasonYear  int32    `json:"seasonYear"`
	Format      string   `json:"format"`
	Status      string   `json:"status"`
	Genres      []string `json:"genres"`
	CoverImage  struct {
		//		Large      string `json:"large"`
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

type rateAniListLimiting struct {
	requestsLeft int
	retryAfter   int
}

func callAnilist(queryString string) (aniListResponse, error, rateAniListLimiting) {
	graphQL_URL := os.Getenv("ANILIST_ENDPOINT")
	var output aniListResponse
	rateLimit := rateAniListLimiting{-1, -1}
	payload := payload{Query: query, Variables: map[string]any{
		"search": queryString,
	}}

	body, err := json.Marshal(payload)
	if err != nil {
		return aniListResponse{}, err, rateAniListLimiting{-1, -1}
	}

	ctx := context.Background()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		graphQL_URL,
		bytes.NewReader(body),
	)

	if err != nil {
		return aniListResponse{}, err, rateAniListLimiting{-1, -1}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 1000 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return aniListResponse{}, err, rateLimit
	}
	defer resp.Body.Close()

	temp1 := resp.Header.Get("X-RateLimit-Remaining")
	if temp1 != "" {
		rateLimit.requestsLeft, err = strconv.Atoi(temp1)
		if err != nil {
			return aniListResponse{}, err, rateAniListLimiting{-1, -1}
		}
	}

	temp2 := resp.Header.Get("X-RateLimit-Reset")
	if temp2 != "" {
		rateLimit.retryAfter, err = strconv.Atoi(temp2)

		if err != nil {
			return aniListResponse{}, err, rateAniListLimiting{-1, -1}
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return aniListResponse{}, fmt.Errorf("HTTP error: %s", resp.Status), rateLimit
	}

	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return aniListResponse{}, err, rateLimit
	}

	if len(output.Errors) > 0 {
		return aniListResponse{}, errors.New(output.Errors[0].Message), rateLimit
	}

	return output, nil, rateLimit

}

func getTopShow(data aniListResponse) (perShow, error) {
	if len(data.Data.Page.Media) != 0 {
		return data.Data.Page.Media[0], nil
	}
	return perShow{}, errors.New("No show found")
}
