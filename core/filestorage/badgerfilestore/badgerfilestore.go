package badgerfilestore

import (
	"context"
	"github.com/anytypeio/any-sync/commonfile/fileblockstore"
	"github.com/dgraph-io/badger/v3"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
)

const keyPrefix = "files/blocks/"

func NewBadgerStorage(db *badger.DB) fileblockstore.BlockStoreLocal {
	return &badgerStorage{db: db}
}

type badgerStorage struct {
	db *badger.DB
}

func (f *badgerStorage) Get(ctx context.Context, k cid.Cid) (b blocks.Block, err error) {
	err = f.db.View(func(txn *badger.Txn) (e error) {
		it, e := txn.Get(key(k))
		if e != nil {
			return e
		}
		if b, e = blockFromItem(it); e != nil {
			return e
		}
		return
	})
	if err == badger.ErrKeyNotFound {
		err = &format.ErrNotFound{Cid: k}
	}
	return
}

func (f *badgerStorage) GetMany(ctx context.Context, ks []cid.Cid) <-chan blocks.Block {
	var res = make(chan blocks.Block)
	go func() {
		defer close(res)
		_ = f.db.View(func(txn *badger.Txn) error {
			// TODO: log errors
			for _, k := range ks {
				it, gerr := txn.Get(key(k))
				if gerr != nil {
					return gerr
				}
				b, berr := blockFromItem(it)
				if berr != nil {
					return berr
				}
				res <- b
			}
			return nil
		})
	}()
	return res
}

func (f *badgerStorage) Add(ctx context.Context, bs []blocks.Block) error {
	return f.db.Update(func(txn *badger.Txn) error {
		for _, b := range bs {
			if err := txn.Set(key(b.Cid()), b.RawData()); err != nil {
				return err
			}
		}
		return nil
	})
}

func (f *badgerStorage) Delete(ctx context.Context, c cid.Cid) error {
	return f.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key(c))
	})
}

func (f *badgerStorage) ExistsCids(ctx context.Context, ks []cid.Cid) (exists []cid.Cid, err error) {
	err = f.db.View(func(txn *badger.Txn) error {
		for _, k := range ks {
			_, e := txn.Get(key(k))
			if e == nil {
				exists = append(exists, k)
			} else if e != badger.ErrKeyNotFound {
				return e
			}
		}
		return nil
	})
	return
}

func (f *badgerStorage) NotExistsBlocks(ctx context.Context, bs []blocks.Block) (notExists []blocks.Block, err error) {
	notExists = bs[:0]
	err = f.db.View(func(txn *badger.Txn) error {
		for _, b := range bs {
			_, e := txn.Get(key(b.Cid()))
			if e == badger.ErrKeyNotFound {
				notExists = append(notExists, b)
			} else if e != nil {
				return e
			}
		}
		return nil
	})
	return
}

func key(c cid.Cid) []byte {
	return []byte(keyPrefix + c.String())
}

func parseCID(key []byte) (cid.Cid, error) {
	if len(key) <= len(keyPrefix) {
		return cid.Cid{}, errInvalidKey
	}
	return cid.Decode(string(key[len(keyPrefix):]))
}

func blockFromItem(it *badger.Item) (b blocks.Block, err error) {
	c, err := parseCID(it.Key())
	if err != nil {
		return nil, err
	}
	if err = it.Value(func(val []byte) error {
		if b, err = blocks.NewBlockWithCid(val, c); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return
	}
	return
}