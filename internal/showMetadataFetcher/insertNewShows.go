package showmetadatafetcher

import (
	"context"

	metadatadb "github.com/philipstanev/Miku-stream/internal/showMetadataFetcher/db"
)

//query: path map

func stringSliceToHashMap(arr []*string) map[string]struct{} {
	output := make(map[string]struct{})

	for _, str := range arr {
		if str == nil {
			continue
		}

		output[*str] = struct{}{}
	}

	return output
}

func (s *Service) filterNewShowsQuery(ctx context.Context, all map[string]string) (map[string]string, error) {
	output := make(map[string]string)
	oldArr, err := s.Q.GetFilePaths(ctx)
	if err != nil {
		return map[string]string{}, err
	}

	oldSet := stringSliceToHashMap(oldArr)

	for query, path := range all {
		_, exists := oldSet[path]
		if exists {
			continue
		}

		output[query] = path
	}
	return output, nil
}

func (s *Service) callFetcher(queries map[string]string) []Response {
	respChan := make(chan Response, len(queries))
	resps := []Response{}
	n := len(queries)
	for query := range queries {
		req := Request{query: query, sendTo: respChan, ttl: -1}
		s.ReqChan <- req
	}

	if n == 0 {
		return []Response{}
	}
	for resp := range respChan {
		if resp.err != nil {
			continue
		}
		resps = append(resps, resp)
		n--
		if n == 0 {
			break
		}
	}

	return resps
}

func (s *Service) DBinsert(ctx context.Context, queries map[string]string, resps []Response) (int64, error) {
	rows := []metadatadb.InsertShowsParams{}
	for _, resp := range resps {
		row := metadatadb.InsertShowsParams{}
		path := queries[resp.query]
		row.Showpath = &path
		row.Titleenglish = &resp.body.Title.English
		row.Titleromanji = &resp.body.Title.Romanji
		row.ShowDescription = &resp.body.Description
		row.Episodes = &resp.body.Episodes
		row.Duration = &resp.body.Duration
		row.Season = &resp.body.Season
		row.Seasonyear = &resp.body.SeasonYear
		row.Showformat = &resp.body.Format
		row.Showstatus = &resp.body.Status
		row.Genres = resp.body.Genres
		row.Coverimageurl = &resp.body.CoverImage.ExtraLarge
		rows = append(rows, row)
	}
	return s.Q.InsertShows(ctx, rows)
}

func (s *Service) InsertForDirectory(ctx context.Context, path string) (int64, error) {
	queries, err := GetQueriesForDirectory(path)
	if err != nil {
		return 0, err
	}
	queries, err = s.filterNewShowsQuery(ctx, queries)
	if err != nil {
		return 0, err
	}
	resps := s.callFetcher(queries)
	inserted, err := s.DBinsert(ctx, queries, resps)
	return inserted, err
}
