package domain

import "time"

// AccountType classifies the behavior of a chart-of-accounts node.
type AccountType string

const (
	AccountTypeAsset           AccountType = "asset"
	AccountTypeLiability       AccountType = "liability"
	AccountTypeEquity          AccountType = "equity"
	AccountTypeRevenue         AccountType = "revenue"
	AccountTypeExpense         AccountType = "expense"
	AccountTypeContraAsset     AccountType = "contra_asset"
	AccountTypeContraLiability AccountType = "contra_liability"
	AccountTypeContraEquity    AccountType = "contra_equity"
	AccountTypeContraRevenue   AccountType = "contra_revenue"
	AccountTypeContraExpense   AccountType = "contra_expense"
)

// Side identifies the debit or credit direction of a journal line.
type Side string

const (
	SideDebit  Side = "debit"
	SideCredit Side = "credit"
)

// Account models a single chart-of-accounts node.
type Account struct {
	ID          string
	CompanyCode string
	Code        string
	Name        string
	ParentID    *string
	Depth       int16
	Type        AccountType
	NormalSide  Side
	IsPostable  bool
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ExpectedNormalSide returns the default normal balance direction for a chart-of-accounts type.
func ExpectedNormalSide(accountType AccountType) Side {
	switch accountType {
	case AccountTypeAsset, AccountTypeExpense, AccountTypeContraLiability, AccountTypeContraEquity, AccountTypeContraRevenue:
		return SideDebit
	case AccountTypeLiability, AccountTypeEquity, AccountTypeRevenue, AccountTypeContraAsset, AccountTypeContraExpense:
		return SideCredit
	default:
		return SideDebit
	}
}

// HasExpectedNormalSide guards the master-data setup used by automated postings and VAS reports.
func (a Account) HasExpectedNormalSide() bool {
	return a.NormalSide == ExpectedNormalSide(a.Type)
}
