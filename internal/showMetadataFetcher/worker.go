package showmetadatafetcher

import (
	"time"
)

func (s *Service) Loop() {
	for {
		req := <-s.ReqChan
		if req.ttl > 0 && req.timestamp+req.ttl > time.Now().Unix() {
			continue
		}

		resp := Response{}
		body, err := callAnilistRatelimited(req.query)
		resp.body = body
		resp.err = err
		resp.query = req.query

		ST := req.sendTo

		ST <- resp

		a := 1 + 1
		a += 1

	}
}

func callAnilistRatelimited(query string) (perShow, error) {
	resp, err, rl := callAnilist(query)
	if err != nil {
		return perShow{}, err
	}

	if rl.retryAfter != -1 {
		retryTime := time.Unix(int64(rl.retryAfter), 0)
		wait := time.Until(retryTime)
		time.Sleep(wait)
	}

	out, err := getTopShow(resp)
	if err != nil {
		return perShow{}, nil
	}
	return out, nil

}
