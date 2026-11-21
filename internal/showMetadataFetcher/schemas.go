package showmetadatafetcher

import (
	metadatadb "github.com/philipstanev/Miku-stream/internal/showMetadataFetcher/db"
)

type Request struct {
	query     string
	sendTo    chan Response
	timestamp int64
	ttl       int64
}

type Response struct {
	body  perShow
	query string
	err   error
}

type Service struct {
	Q       *metadatadb.Queries
	ReqChan chan Request
}
