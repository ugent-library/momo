package es6

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ugent-library/momo/internal/engine"
)

type recCursor struct {
	store    *store
	args     engine.SearchArgs
	scrollID *string
	err      error
	recs     []*engine.Rec
	recsIdx  int
}

func (c *recCursor) Next() bool {
	if c.scrollID == nil {
		c.err = c.initSearch()
		if c.err == nil && len(c.recs) > 0 {
			return true
		}
	} else if c.recsIdx < len(c.recs) {
		return true
	} else {
		c.err = c.scrollSearch()
		c.recsIdx = 0
		if c.err == nil && len(c.recs) > 0 {
			return true
		}
	}
	return false
}

func (c *recCursor) Value() *engine.Rec {
	rec := c.recs[c.recsIdx]
	c.recsIdx++
	return rec
}

func (c *recCursor) Error() error {
	return c.err
}

func (c *recCursor) Close() {
}

func (c *recCursor) initSearch() error {
	client := c.store.client
	query, _, _ := buildQuery(c.args)
	query["size"] = 100

	r, err := c.store.search(query,
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(c.store.indexName("rec")),
		client.Search.WithSort("_doc"),
		client.Search.WithScroll(time.Minute),
	)
	if err != nil {
		return err
	}

	c.scrollID = &r.ScrollID

	for _, h := range r.Hits.Hits {
		var rec engine.Rec

		if err := json.Unmarshal(h.Source, &rec); err != nil {
			return err
		}

		c.recs = append(c.recs, &rec)
	}

	return nil
}

func (c *recCursor) scrollSearch() error {
	client := c.store.client

	res, err := client.Scroll(
		client.Scroll.WithScrollID(*c.scrollID),
		client.Scroll.WithScroll(time.Minute),
	)
	if err != nil {
		return err
	}

	r, err := decodeRes(res)
	if err != nil {
		return err
	}

	c.scrollID = &r.ScrollID

	// clear recs but keep allocated array
	c.recs = c.recs[:0]

	for _, h := range r.Hits.Hits {
		var rec engine.Rec

		if err := json.Unmarshal(h.Source, &rec); err != nil {
			return err
		}

		c.recs = append(c.recs, &rec)
	}

	return nil
}
