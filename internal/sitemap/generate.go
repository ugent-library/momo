package sitemap

// import (
// 	"time"

// 	"github.com/ugent-library/momo/internal/engine"
// )

// func Generate(e engine.Engine, dir string) error {
// 	now := time.Now().Unix()
// 	n := 1
// 	c := e.AllRecs()
// 	defer c.Close()
// 	for c.Next() {
// 		if err := c.Error(); err != nil {
// 			return err
// 		}
// 		rec := c.Value()
// 		n++
// 	}

// 	return nil
// }
