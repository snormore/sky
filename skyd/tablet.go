package skyd

import (
	"github.com/jmhodges/levigo"
)

// A Tablet is a small wrapper for the underlying data storage
// which is contained in LevelDB.
type Tablet struct {
	db   *levigo.DB
	path string
}

// NewTablet returns a new Tablet that is stored at a given path.
func NewTablet(path string) Tablet {
  return Tablet{
    path: path,
  }
}

// Opens the underlying LevelDB database.
func (t *Tablet) Open() error {
	opts := levigo.NewOptions()
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(t.path, opts)
	if err != nil {
		return err
	}
	t.db = db
	return nil
}

// Closes the underlying LevelDB database.
func (t *Tablet) Close() {
  if t.db != nil {
	  t.db.Close()
  }
}

