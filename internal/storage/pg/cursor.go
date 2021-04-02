package pg

import (
	"database/sql"

	"github.com/ugent-library/momo/internal/engine"
	"gorm.io/gorm"
)

type recCursor struct {
	db   *gorm.DB
	rows *sql.Rows
	err  error
	val  *engine.Rec
}

func (c *recCursor) Next() bool {
	if c.rows.Next() {
		rec := Rec{}
		if err := c.db.ScanRows(c.rows, &rec); err != nil {
			c.err = err
			c.val = nil
		} else {
			c.err = nil
			c.val = reifyRec(&rec)
		}
		return true
	}
	return false
}

func (c *recCursor) Value() *engine.Rec {
	return c.val
}

func (c *recCursor) Error() error {
	return c.err
}

func (c *recCursor) Close() {
	c.rows.Close()
}
