package model

import (
	"sort"
	"sync"
	"time"

	"github.com/ChristianNorbertBraun/seaweed-banking/seaweed-banking-account-updater/config"
)

// AccountInfo holds all information for an account from the oldest to
// the latest transaction
type AccountInfo struct {
	Name              string       `json:"name"`
	BIC               string       `json:"bic"`
	IBAN              string       `json:"iban"`
	Balance           int32        `json:"balance"`
	Predecessor       string       `json:"predecessor"`
	OldestTransaction string       `json:"oldestTransaction"`
	LatestTransaction string       `json:"latestTransaction"`
	Transactions      Transactions `json:"transactions"`
	mutex             sync.Mutex
}

// NewAccountInfo creates a new accountInfo
func NewAccountInfo(name string, bic string, iban string, balance int32, predecessor string) *AccountInfo {
	accountInfo := AccountInfo{Name: name, BIC: bic, IBAN: iban, Balance: balance, Predecessor: predecessor}
	accountInfo.Transactions = []*Transaction{}

	return &accountInfo
}

// AddTransaction adds a Transaction to the AccountInfo and updates the balance and
// the oldest and latest transaction date
func (ai *AccountInfo) AddTransaction(transaction *Transaction) (bool, *AccountInfo) {
	ai.mutex.Lock()
	defer ai.mutex.Unlock()
	if ai.Transactions.Len() < config.Configuration.Seaweed.MaxTransactionsPerAccountInfo {
		ai.Transactions = append(ai.Transactions, transaction)
		ai.Balance += transaction.ValueInSmallestUnit
		sort.Sort(ai.Transactions)

		ai.OldestTransaction = ai.Transactions[0].
			BookingDate.
			Format(time.RFC3339Nano)
		ai.LatestTransaction = ai.Transactions[ai.Transactions.Len()-1].
			BookingDate.
			Format(time.RFC3339Nano)

		return true, nil
	}

	accountInfo := NewAccountInfo(ai.Name, ai.BIC, ai.IBAN, ai.Balance, ai.OldestTransaction)
	accountInfo.AddTransaction(transaction)

	return false, accountInfo
}
