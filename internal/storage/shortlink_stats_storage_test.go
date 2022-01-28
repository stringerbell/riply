package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShortLinkStorage_RecordHit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	sto := NewShortLinkStorage(db)
	t.Run("recording_hits_works", func(t *testing.T) {
		suffix := uniqueHash()
		res, err := sto.Save(ShortLinkRequest{
			Link:         uniqueHash(), // needs to be different
			CustomSuffix: &suffix,
		})
		if err != nil {
			t.Fatalf("got unexpected err: %s", err)
		}
		sto.RecordHit(nil, res)
		sto.RecordHit(nil, res)
		sto.RecordHit(nil, res)
		stats, err := sto.GetStatsForShortLink(suffix)
		assert.NoError(t, err)
		now := time.Now().UTC()
		assert.Equal(t, ShortLinkStats{
			Total: 3,
			DateTotals: []DateTotal{
				{
					Date:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
					Total: 3,
				},
			},
		}, stats)
	})
}
