package storage

import (
	"fmt"
	"net/http"
	"time"
)

type DateTotal struct {
	Date  time.Time `json:"date"`
	Total int       `json:"total"`
}

type ShortLinkStats struct {
	Total      int         `json:"total"`
	DateTotals []DateTotal `json:"stats,omitempty"`
}

func (s *ShortLinkStorage) RecordHit(r *http.Request, resp ShortLinkResponse) error {
	// TODO do something interesting with r
	_, err := s.db.Exec("insert into stats (shortlink_id) values ((select id from shortlinks where link = $1))", resp.Link)
	return err
}

func (s *ShortLinkStorage) GetStatsForShortLink(shortlink string) (ShortLinkStats, error) {
	r, err := s.GetOneShortLink(shortlink)
	if err != nil {
		return ShortLinkStats{}, fmt.Errorf("s.GetOneShortLink err: %w", err)
	}
	rows, err := s.db.Query(`
        SELECT EXTRACT(year FROM stats.created_at)  as year,
               EXTRACT(month FROM stats.created_at) as month,
               EXTRACT(day FROM stats.created_at)   as day,
               sum(1)                               as total
        FROM stats
        where shortlink_id = $1
        GROUP BY year, month, day;
        `, r.id)
	if err != nil {
		return ShortLinkStats{}, fmt.Errorf("s.db.Query err: %w", err)
	}
	var sts ShortLinkStats
	var totalAcross int
	for rows.Next() {
		var year, month, day, total int
		err := rows.Scan(&year, &month, &day, &total)
		if err != nil {
			return ShortLinkStats{}, err
		}
		sts.DateTotals = append(sts.DateTotals, DateTotal{Total: total, Date: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)})
		totalAcross += total
	}
	sts.Total = totalAcross
	return sts, rows.Err()
}
