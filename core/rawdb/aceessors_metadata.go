package rawdb

import (
	"encoding/json"
	"github.com/srchain/srcd/common/common"
	"github.com/srchain/srcd/log"
	"github.com/srchain/srcd/params"
)

// ReadChainConfig retrieves the consensus settings based on the given genesis hash.
func ReadChainConfig(db DatabaseReader, hash common.Hash) *params.ChainConfig {
	data, _ := db.Get(configKey(hash))
	if len(data) == 0 {
		return nil
	}
	var config params.ChainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Error("Invalid chain config JSON", "hash", hash, "err", err)
		return nil
	}
	return &config
}


// WriteChainConfig writes the chain config settings to the database.
func WriteChainConfig(db DatabaseWriter, hash common.Hash, cfg *params.ChainConfig) {
	if cfg == nil {
		return
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Crit("Failed to JSON encode chain config", "err", err)
	}
	if err := db.Put(configKey(hash), data); err != nil {
		log.Crit("Failed to store chain config", "err", err)
	}
}



