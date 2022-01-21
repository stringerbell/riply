package storage

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestShortLinkStorage_Save(t *testing.T) {
  if testing.Short() {
    t.Skip("skipping test in short mode..")
  }
  sto := NewShortLinkStorage(db)
  t.Run("saving_the_same_link_twice_is_a_conflict", func(t *testing.T) {
    _, err := sto.Save(ShortLinkRequest{
      Link:         "a link",
    })
    assert.NoError(t, err)

    _, err = sto.Save(ShortLinkRequest{
      Link:         "a link",
    })
    var conflict ErrConflict
    assert.ErrorAs(t, err, &conflict)
  })
}
