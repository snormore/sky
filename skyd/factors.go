package skyd

import (
	"errors"
	"fmt"
	"github.com/benbjohnson/gomdb"
	"os"
	"strconv"
	"sync"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Factors object manages the factorization and defactorization of values.
type Factors struct {
	env   *mdb.Env
	path  string
	mutex sync.Mutex
}

//------------------------------------------------------------------------------
//
// Errors
//
//------------------------------------------------------------------------------

//--------------------------------------
// Factor Not Found
//--------------------------------------

func NewFactorNotFound(text string) error {
	return &FactorNotFound{text}
}

type FactorNotFound struct {
	s string
}

func (e *FactorNotFound) Error() string {
	return e.s
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// NewFactors returns a new Factors object.
func NewFactors(path string) *Factors {
	return &Factors{path: path}
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// The path to the database on disk.
func (f *Factors) Path() string {
	return f.path
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// State
//--------------------------------------

// Opens the factors databse.
func (f *Factors) Open() error {
	var err error
	if f.IsOpen() {
		return errors.New("skyd.Factors: Factors database is already open.")
	}

	// Create the factors directory.
	if err = os.MkdirAll(f.path, 0700); err != nil {
		return err
	}

	// Create the database environment.
	if f.env, err = mdb.NewEnv(); err != nil {
		return fmt.Errorf(fmt.Sprintf("skyd: Unable to create factors environment: %v", err))
	}
	// Setup max dbs.
	if err = f.env.SetMaxDBs(1024); err != nil {
		f.Close()
		return fmt.Errorf("skyd: Unable to set factors max dbs: %v", err)
	}
	// Open the database.
	if err = f.env.Open(f.path, 0, 0664); err != nil {
		f.Close()
		return fmt.Errorf("skyd: Cannot open factors database (%s): %s", f.path, err)
	}

	return nil
}

// Closes the factors database.
func (f *Factors) Close() {
	if f.env != nil {
		f.env.Close()
		f.env = nil
	}
}

// Returns whether the factors database is open.
func (f *Factors) IsOpen() bool {
	return f.env != nil
}

//--------------------------------------
// DB
//--------------------------------------

// Retrieves the value from the database for a given key.
func (f *Factors) get(namespace string, key string) (string, bool, error) {
	txn, err := f.env.BeginTxn(nil, 0)
	if err != nil {
		return "", false, fmt.Errorf("skyd: Unable to start factors get txn: %s", err)
	}
	dbi, err := txn.DBIOpen(&namespace, mdb.CREATE)
	if err != nil {
		return "", false, fmt.Errorf("skyd: Unable to open factors DBI [get]: %s", err)
	}

	// Retrieve byte array.
	data, err := txn.Get(dbi, []byte(key))
	if err != nil && err != mdb.NotFound {
		txn.Abort()
		return "", false, fmt.Errorf("skyd: Unable to get factor: %s", err)
	}
	txn.Abort()

	return string(data), (data != nil), nil
}

// Sets the value for a given key in the database.
func (f *Factors) put(namespace string, key string, value string) error {
	txn, err := f.env.BeginTxn(nil, 0)
	if err != nil {
		return fmt.Errorf("skyd: Unable to start factors put txn: %s", err)
	}
	dbi, err := txn.DBIOpen(&namespace, mdb.CREATE)
	if err != nil {
		return fmt.Errorf("skyd: Unable to open factors DBI [put]: %s", err)
	}

	// Set value for key.
	if err = txn.Put(dbi, []byte(key), []byte(value), mdb.NODUPDATA); err != nil {
		txn.Abort()
		return fmt.Errorf("skyd: Unable to put factor: %s", err)
	}
	if err = txn.Commit(); err != nil {
		txn.Abort()
		return fmt.Errorf("skyd: Unable to commit factor: %s", err)
	}

	return nil
}

//--------------------------------------
// Keys
//--------------------------------------

// The key for a given namespace/id/value.
func (f *Factors) key(id string, value string) string {
	return fmt.Sprintf("%x:%s>%s", len(id), id, value)
}

// The reverse key for a given namespace/id/value.
func (f *Factors) revkey(id string, value uint64) string {
	return fmt.Sprintf("%x:%s<%d", len(id), id, value)
}

// The sequence key for a given namespace/id.
func (f *Factors) seqkey(id string) string {
	return fmt.Sprintf("%x:%s!", len(id), id)
}

//--------------------------------------
// Factorization
//--------------------------------------

// Converts the defactorized value for a given id in a given namespace to its internal representation.
func (f *Factors) Factorize(namespace string, id string, value string, createIfMissing bool) (uint64, error) {
	// Blank is always zero.
	if value == "" {
		return 0, nil
	}

	// Otherwise find it in the database.
	data, exists, err := f.get(namespace, f.key(id, value))
	if err != nil {
		return 0, err
	}
	// If key does exist then parse and return it.
	if exists {
		return strconv.ParseUint(string(data), 10, 64)
	}

	// Create a new factor if requested.
	if createIfMissing {
		return f.add(namespace, id, value)
	}

	err = NewFactorNotFound(fmt.Sprintf("skyd.Factors: Factor not found: %v", f.key(id, value)))
	return 0, err
}

// Adds a new factor to the database if it doesn't exist.
func (f *Factors) add(namespace string, id string, value string) (uint64, error) {
	// Lock while adding a new value.
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Retry factorize within the context of the lock.
	sequence, err := f.Factorize(namespace, id, value, false)
	if err == nil {
		return sequence, nil
	} else if _, ok := err.(*FactorNotFound); !ok {
		return 0, err
	}

	// Retrieve next id in sequence.
	sequence, err = f.inc(namespace, id)
	if err != nil {
		return 0, err
	}

	// Save lookup and reverse lookup.
	if err = f.put(namespace, f.key(id, value), strconv.FormatUint(sequence, 10)); err != nil {
		return 0, err
	}
	if err = f.put(namespace, f.revkey(id, sequence), value); err != nil {
		return 0, err
	}

	return sequence, nil
}

// Converts the factorized value for a given id in a given namespace to its internal representation.
func (f *Factors) Defactorize(namespace string, id string, value uint64) (string, error) {
	// Blank is always zero.
	if value == 0 {
		return "", nil
	}

	// Find it in the database.
	data, exists, err := f.get(namespace, f.revkey(id, value))
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("skyd: Factor value does not exist: %v", f.revkey(id, value))
	}
	return string(data), nil
}

// Retrieves the next available sequence number within a namespace for an id.
func (f *Factors) inc(namespace string, id string) (uint64, error) {
	data, exists, err := f.get(namespace, f.seqkey(id))
	if err != nil {
		return 0, err
	}

	// Initialize key if it doesn't exist. Otherwise increment it.
	if !exists {
		if err := f.put(namespace, f.seqkey(id), "1"); err != nil {
			return 0, err
		}
		return 1, nil
	}

	// Parse existing sequence.
	sequence, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("skyd.Factors: Unable to parse sequence: %v", data)
	}

	// Increment and save the new value.
	sequence += 1
	if err = f.put(namespace, f.seqkey(id), strconv.FormatUint(sequence, 10)); err != nil {
		return 0, err
	}
	return sequence, nil
}
