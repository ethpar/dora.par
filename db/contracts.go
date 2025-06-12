package db

import "github.com/ethpandaops/dora/dbtypes"

func GetContracts() map[string]*dbtypes.Contract {
	contracts := []*dbtypes.Contract{}
	err := ReaderDb.Select(&contracts, `
	SELECT
		address, owner, is_erc20, name, symbol, created_at, body
	FROM contracts
	ORDER BY address
	`)
	if err != nil {
		logger.Errorf("Error while fetching Contracts: %v", err)
		return nil
	}

	result := make(map[string]*dbtypes.Contract, len(contracts))
	for _, item := range contracts {
		result[item.Address] = item
	}
	return result
}
