package transaction

import (
	"context"
	"encoding/json"

	"github.com/srchain/srcd/database"
	"github.com/srchain/srcd/errors"
)

//Tracker filter tracker object.
type TxFeedManager struct {
	DB       database.LDBDatabase
	TxFeeds  []*TxFeed
	//chain    *blockchain.BlockChain
	txfeedCh chan *Tx
}

//TxFeed describe a filter
type TxFeed struct {
	ID     string `json:"id,omitempty"`
	Alias  string `json:"alias"`
	Filter string `json:"filter,omitempty"`
	Param  filter `json:"param,omitempty"`
}

type filter struct {
	AssetID          string `json:"assetid,omitempty"`
	AmountLowerLimit uint64 `json:"lowerlimit,omitempty"`
	AmountUpperLimit uint64 `json:"upperlimit,omitempty"`
	TransType        string `json:"transtype,omitempty"`
}

func NewTxFeedManager(db database.LDBDatabase) *TxFeedManager {
	s := &TxFeedManager{
		DB:       db,
		TxFeeds:  make([]*TxFeed, 0, 10),
		//chain:    chain,
		txfeedCh: make(chan *Tx, 100),
	}

	return s
}

//new or load txfeed from db
func (t *TxFeedManager) Prepare(ctx context.Context) error {
	var err error
	t.TxFeeds, err = loadTxFeed(t.DB, t.TxFeeds)
	return err
}

func loadTxFeed(db database.LDBDatabase, txFeeds []*TxFeed) ([]*TxFeed, error) {
	iter := db.NewIterator()
	defer iter.Release()

	for iter.Next() {
		txFeed := &TxFeed{}
		if err := json.Unmarshal(iter.Value(), &txFeed); err != nil {
			return nil, err
		}
		//filter, err := parseFilter(txFeed.Filter)
		//if err != nil {
		//	return nil, err
		//}
		//txFeed.Param = filter
		txFeeds = append(txFeeds, txFeed)
	}
	return txFeeds, nil
}

func (t *TxFeedManager)AddTxFeed(ctx context.Context,alias string ,filter string) error  {
	for _ ,txfeed := range t.TxFeeds {
		if txfeed.Alias == alias{
			return errors.New("alias must be unique")
		}
	}
	feed := &TxFeed{
		Alias:  alias,
		Filter: filter,
	}
	t.TxFeeds = append(t.TxFeeds, feed)

	key, err := json.Marshal(feed.Alias)
	if err != nil {
		return errors.New("alias marshal err")
	}
	value, err := json.Marshal(feed)
	if err != nil {
		return errors.New("feed marshal err")
	}
	t.DB.Put(key, value)
	return nil
}

func (t *TxFeedManager)DeleteTxFeed(ctx context.Context,alias string) error  {
	for i ,txfeed := range t.TxFeeds {
		if txfeed.Alias == alias{
			t.TxFeeds = append(t.TxFeeds[:i], t.TxFeeds[i+1:]...)
			key, err := json.Marshal(alias)
			if err != nil {
				return errors.New("alias marshal err")
			}
			t.DB.Delete(key)
		}
	}
	return nil
}

func (t *TxFeedManager)GetTxFeed(ctx context.Context,alias string)(*TxFeed,error)  {
	for i, txfeed := range t.TxFeeds {
		if txfeed.Alias == alias{
			return t.TxFeeds[i],nil
		}
	}
	return nil,nil
}
