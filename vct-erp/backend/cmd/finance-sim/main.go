package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"vct-platform/backend/internal/config"
	infraPostgres "vct-platform/backend/internal/infra/postgres"
	analyticspg "vct-platform/backend/internal/modules/analytics/adapter/postgres"
	analyticsusecase "vct-platform/backend/internal/modules/analytics/usecase"
	financepg "vct-platform/backend/internal/modules/finance/adapter/postgres"
	financedomain "vct-platform/backend/internal/modules/finance/domain"
	financeusecase "vct-platform/backend/internal/modules/finance/usecase"
	ledgerpg "vct-platform/backend/internal/modules/ledger/adapter/postgres"
	ledgerdomain "vct-platform/backend/internal/modules/ledger/domain"
	ledgerusecase "vct-platform/backend/internal/modules/ledger/usecase"
	"vct-platform/backend/internal/shared/id"
)

const (
	simulationCompany = "VCT_SIM"
	baseCompany       = "VCT_GROUP"
	reportDateLayout  = "2006-01-02"
)

var (
	reportFrom = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	reportTo   = time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC)
)

type simulator struct {
	db               *sql.DB
	ledgerStore      *ledgerpg.Store
	financeStore     *financepg.Store
	postEntryUC      *ledgerusecase.PostEntryUseCase
	saasService      *financeusecase.SaaSService
	dojoService      *financeusecase.DojoService
	retailService    *financeusecase.RetailService
	rentalService    *financeusecase.RentalService
	captureUC        *financeusecase.CaptureUseCase
	voidUC           *financeusecase.VoidTransactionUseCase
	analyticsService *analyticsusecase.Service
	accounts         map[string]string
	rng              *rand.Rand
	stats            seedStats
}

type seedStats struct {
	baseOps             int
	duplicatesAttempted int
	duplicatesBlocked   int
	saasContracts       int
	saasRecognitions    int
	dojoAssessments     int
	dojoPayments        int
	retailSales         int
	retailRefunds       int
	rentalCaptures      int
	rentalSettlements   int
	expenseEntries      int
	voids               int
	trialBalance        []trialBalanceRow
	pnlRows             []reportRow
	grossRows           []grossProfitRow
	revenueStream       []reportAmountRow
	cashRunway          analyticsSnapshot
	partitionCounts     []partitionRow
	consistency         []string
}

type trialBalanceRow struct {
	AccountCode    string
	AccountName    string
	NormalSide     string
	OpeningBalance string
	PeriodDebit    string
	PeriodCredit   string
	ClosingBalance string
}

type reportRow struct {
	LineCode string
	LineName string
	Amount   string
}

type grossProfitRow struct {
	CostCenter        string
	GrossRevenue      string
	RevenueDeductions string
	OtherIncome       string
	CostOfGoodsSold   string
	GrossProfit       string
}

type reportAmountRow struct {
	Label  string
	Amount string
}

type partitionRow struct {
	PartitionName string
	RowCount      int
}

type analyticsSnapshot struct {
	CurrentCash        string
	AverageMonthlyBurn string
	Months             []analyticsMonth
}

type analyticsMonth struct {
	MonthLabel       string
	OpeningCash      string
	ContractedInflow string
	ProjectedBurn    string
	ProjectedEnding  string
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := infraPostgres.Open(ctx, cfg.DatabaseDSN, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	sim := newSimulator(db)
	if err := sim.prepareDatabase(ctx); err != nil {
		log.Fatalf("prepare database: %v", err)
	}
	if err := sim.seedAll(ctx); err != nil {
		log.Fatalf("seed simulation: %v", err)
	}
	if err := sim.collectSnapshots(ctx); err != nil {
		log.Fatalf("collect snapshots: %v", err)
	}
	if err := sim.writeArtifacts(); err != nil {
		log.Fatalf("write artifacts: %v", err)
	}

	fmt.Println(sim.renderMarkdown())
}

func newSimulator(db *sql.DB) *simulator {
	ledgerStore := ledgerpg.NewStore(db)
	financeStore := financepg.NewStore(db)
	postEntryUC := ledgerusecase.NewPostEntryUseCase(
		ledgerStore,
		ledgerStore,
		ledgerStore,
		ledgerStore,
		ledgerStore,
		ledgerStore,
		nil,
		0,
		nil,
		"ledger.events",
	)
	saasService := financeusecase.NewSaaSService(ledgerStore, postEntryUC, financeStore)
	dojoService := financeusecase.NewDojoService(ledgerStore, postEntryUC, financeStore)
	retailService := financeusecase.NewRetailService(postEntryUC)
	rentalService := financeusecase.NewRentalService(ledgerStore, postEntryUC, financeStore)

	return &simulator{
		db:               db,
		ledgerStore:      ledgerStore,
		financeStore:     financeStore,
		postEntryUC:      postEntryUC,
		saasService:      saasService,
		dojoService:      dojoService,
		retailService:    retailService,
		rentalService:    rentalService,
		captureUC:        financeusecase.NewCaptureUseCase(financeStore, saasService, dojoService, retailService, rentalService),
		voidUC:           financeusecase.NewVoidTransactionUseCase(ledgerStore, ledgerStore, postEntryUC, nil),
		analyticsService: analyticsusecase.NewService(analyticspg.NewRepository(db)),
		rng:              rand.New(rand.NewSource(20260321)),
	}
}

func (s *simulator) prepareDatabase(ctx context.Context) error {
	exists, err := tableExists(ctx, s.db, "accounts")
	if err != nil {
		return err
	}
	if !exists {
		if err := executeSQLFile(ctx, s.db, "schema.sql"); err != nil {
			return fmt.Errorf("apply schema.sql: %w", err)
		}
	}
	if err := executeSQLFile(ctx, s.db, filepath.Join("migrations", "202603210001_vas_tt200_reverse_reporting.sql")); err != nil {
		return fmt.Errorf("apply migration: %w", err)
	}
	if err := s.ensureSimulationAccounts(ctx); err != nil {
		return fmt.Errorf("ensure simulation chart of accounts: %w", err)
	}
	if err := s.resetSimulationData(ctx); err != nil {
		return fmt.Errorf("reset simulation data: %w", err)
	}
	if err := s.loadAccounts(ctx); err != nil {
		return fmt.Errorf("load simulation accounts: %w", err)
	}
	return nil
}

func (s *simulator) seedAll(ctx context.Context) error {
	if err := s.seedSaaS(ctx); err != nil {
		return err
	}
	if err := s.seedDojo(ctx); err != nil {
		return err
	}
	retailSaleEntries, err := s.seedRetailAndRental(ctx)
	if err != nil {
		return err
	}
	if err := s.seedOperatingExpenses(ctx); err != nil {
		return err
	}
	if err := s.seedSaaSRecognitions(ctx); err != nil {
		return err
	}
	if err := s.seedVoidCases(ctx, retailSaleEntries); err != nil {
		return err
	}
	return nil
}

func (s *simulator) seedSaaS(ctx context.Context) error {
	for index := 0; index < 300; index++ {
		startDate := monthStart(randomDateInRange(s.rng, reportFrom, reportTo))
		monthlyFee := moneyFromInt(900_000 + s.rng.Intn(3_500_000))
		totalAmount := percentageOf(monthlyFee, 1200)
		payload := financeusecase.CaptureAnnualContractRequest{
			CompanyCode:                simulationCompany,
			ContractNo:                 fmt.Sprintf("SAAS-%04d", index+1),
			CustomerRef:                fmt.Sprintf("SAAS-CUST-%04d", index+1),
			CashAccountID:              s.pickCashAccount(),
			DeferredRevenueAccountID:   s.account("3387"),
			RecognizedRevenueAccountID: s.account("5113"),
			CurrencyCode:               "VND",
			ServiceStartDate:           startDate,
			CapturedAt:                 startDate,
			TermMonths:                 12,
			TotalAmount:                totalAmount,
			SourceRef:                  fmt.Sprintf("SIM-SAAS-%04d", index+1),
		}

		if _, err := s.submitCapture(ctx, financedomain.BusinessLineSaaS, financedomain.OperationSaaSCaptureAnnualContract, payload); err != nil {
			return fmt.Errorf("seed saas contract %d: %w", index+1, err)
		}
		s.stats.saasContracts++

		cost := percentageOf(totalAmount, 18)
		if err := s.postDirectCost(ctx, startDate, "saas", s.account("632"), s.pickBankAccount(), cost, fmt.Sprintf("Chi phi trien khai hop dong %s", payload.ContractNo)); err != nil {
			return fmt.Errorf("seed saas direct cost %d: %w", index+1, err)
		}
	}
	return nil
}

func (s *simulator) seedSaaSRecognitions(ctx context.Context) error {
	monthCursor := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for !monthCursor.After(reportTo) {
		payload := financeusecase.RecognizeDueRevenueRequest{
			CompanyCode: simulationCompany,
			UpTo:        monthCursor,
			Limit:       5000,
		}
		result, err := s.submitCapture(ctx, financedomain.BusinessLineSaaS, financedomain.OperationSaaSRecognizeDueRevenue, payload)
		if err != nil {
			return fmt.Errorf("recognize saas revenue for %s: %w", monthCursor.Format("2006-01"), err)
		}

		var typed financeusecase.RecognizeDueRevenueResult
		if err := json.Unmarshal(result.Payload, &typed); err != nil {
			return fmt.Errorf("decode saas recognition result: %w", err)
		}
		s.stats.saasRecognitions += typed.RecognizedCount
		monthCursor = monthCursor.AddDate(0, 1, 0)
	}
	return nil
}

func (s *simulator) seedDojo(ctx context.Context) error {
	type receivableSeed struct {
		studentRef   string
		billingMonth time.Time
		amount       ledgerdomain.Money
	}

	receivables := make([]receivableSeed, 0, 220)
	for index := 0; index < 220; index++ {
		billingMonth := monthStart(randomDateInRange(s.rng, reportFrom, reportTo))
		amount := moneyFromInt(1_200_000 + s.rng.Intn(1_600_000))
		studentRef := fmt.Sprintf("DOJO-STU-%04d", index+1)
		payload := financeusecase.AssessMonthlyTuitionRequest{
			CompanyCode:         simulationCompany,
			StudentRef:          studentRef,
			BillingMonth:        billingMonth,
			DueDate:             billingMonth.AddDate(0, 0, 7),
			ReceivableAccountID: s.account("131"),
			RevenueAccountID:    s.account("5113"),
			CurrencyCode:        "VND",
			Amount:              amount,
			SourceRef:           fmt.Sprintf("DOJO-TUITION-%04d", index+1),
		}

		if _, err := s.submitCapture(ctx, financedomain.BusinessLineDojo, financedomain.OperationDojoAssessMonthlyTuition, payload); err != nil {
			return fmt.Errorf("seed dojo tuition assessment %d: %w", index+1, err)
		}
		s.stats.dojoAssessments++
		receivables = append(receivables, receivableSeed{
			studentRef:   studentRef,
			billingMonth: billingMonth,
			amount:       amount,
		})

		cost := percentageOf(amount, 32)
		if err := s.postDirectCost(ctx, billingMonth, "dojo", s.account("632"), s.pickCashAccount(), cost, fmt.Sprintf("Chi phi dao tao hoc vien %s", studentRef)); err != nil {
			return fmt.Errorf("seed dojo direct cost %d: %w", index+1, err)
		}
	}

	for index := 0; index < 120; index++ {
		rec := receivables[index]
		paidAt := rec.billingMonth.AddDate(0, 0, 3+s.rng.Intn(20))
		if paidAt.After(reportTo) {
			paidAt = reportTo
		}
		payload := financeusecase.CaptureDojoPaymentRequest{
			CompanyCode:   simulationCompany,
			StudentRef:    rec.studentRef,
			BillingMonth:  rec.billingMonth,
			PaidAt:        paidAt,
			CashAccountID: s.pickCashAccount(),
			CurrencyCode:  "VND",
			PaymentAmount: rec.amount,
			SourceRef:     fmt.Sprintf("DOJO-PAY-%04d", index+1),
		}
		if _, err := s.submitCapture(ctx, financedomain.BusinessLineDojo, financedomain.OperationDojoCapturePayment, payload); err != nil {
			return fmt.Errorf("seed dojo payment %d: %w", index+1, err)
		}
		s.stats.dojoPayments++
	}

	for index := 0; index < 30; index++ {
		examMonth := monthStart(randomDateInRange(s.rng, reportFrom, reportTo))
		studentRef := fmt.Sprintf("DOJO-EXAM-%04d", index+1)
		amount := moneyFromInt(350_000 + s.rng.Intn(450_000))

		assessment := financeusecase.AssessMonthlyTuitionRequest{
			CompanyCode:         simulationCompany,
			StudentRef:          studentRef,
			BillingMonth:        examMonth,
			DueDate:             examMonth,
			ReceivableAccountID: s.account("131"),
			RevenueAccountID:    s.account("5113"),
			CurrencyCode:        "VND",
			Amount:              amount,
			SourceRef:           fmt.Sprintf("DOJO-EXAM-FEE-%04d", index+1),
		}
		if _, err := s.submitCapture(ctx, financedomain.BusinessLineDojo, financedomain.OperationDojoAssessMonthlyTuition, assessment); err != nil {
			return fmt.Errorf("seed dojo exam assessment %d: %w", index+1, err)
		}
		s.stats.dojoAssessments++

		payment := financeusecase.CaptureDojoPaymentRequest{
			CompanyCode:   simulationCompany,
			StudentRef:    studentRef,
			BillingMonth:  examMonth,
			PaidAt:        examMonth,
			CashAccountID: s.pickCashAccount(),
			CurrencyCode:  "VND",
			PaymentAmount: amount,
			SourceRef:     fmt.Sprintf("DOJO-EXAM-PAY-%04d", index+1),
		}
		if _, err := s.submitCapture(ctx, financedomain.BusinessLineDojo, financedomain.OperationDojoCapturePayment, payment); err != nil {
			return fmt.Errorf("seed dojo exam payment %d: %w", index+1, err)
		}
		s.stats.dojoPayments++
	}
	return nil
}

func (s *simulator) seedRetailAndRental(ctx context.Context) ([]string, error) {
	saleEntries := make([]string, 0, 150)
	refundableOrders := make([]string, 0, 90)
	saleDates := make(map[string]time.Time, 90)

	for index := 0; index < 150; index++ {
		orderDate := randomDateInRange(s.rng, reportFrom, reportTo)
		amount := moneyFromInt(250_000 + s.rng.Intn(2_750_000))
		payload := financeusecase.CaptureRetailSaleRequest{
			CompanyCode:      simulationCompany,
			OrderNo:          fmt.Sprintf("POS-%04d", index+1),
			CustomerRef:      fmt.Sprintf("RTL-CUST-%04d", 1+s.rng.Intn(80)),
			OccurredAt:       orderDate,
			CashAccountID:    s.pickCashAccount(),
			RevenueAccountID: s.account("5111"),
			CurrencyCode:     "VND",
			Amount:           amount,
			SourceRef:        fmt.Sprintf("POS-SALE-%04d", index+1),
		}

		result, err := s.submitCapture(ctx, financedomain.BusinessLineRetail, financedomain.OperationRetailCaptureSale, payload)
		if err != nil {
			return nil, fmt.Errorf("seed retail sale %d: %w", index+1, err)
		}
		var typed financeusecase.CaptureRetailSaleResult
		if err := json.Unmarshal(result.Payload, &typed); err != nil {
			return nil, fmt.Errorf("decode retail sale result: %w", err)
		}
		saleEntries = append(saleEntries, typed.JournalEntryID)
		if index < 90 {
			refundableOrders = append(refundableOrders, payload.OrderNo)
			saleDates[payload.OrderNo] = orderDate
		}
		s.stats.retailSales++

		cost := percentageOf(amount, 58)
		if err := s.postDirectCost(ctx, orderDate, "retail", s.account("632"), s.pickBankAccount(), cost, fmt.Sprintf("Gia von POS %s", payload.OrderNo)); err != nil {
			return nil, fmt.Errorf("seed retail cost %d: %w", index+1, err)
		}
	}

	for index := 0; index < 30; index++ {
		refundDate := saleDates[refundableOrders[index]].AddDate(0, 0, 1+s.rng.Intn(21))
		if refundDate.After(reportTo) {
			refundDate = reportTo
		}
		amount := moneyFromInt(120_000 + s.rng.Intn(550_000))
		payload := financeusecase.CaptureRetailRefundRequest{
			CompanyCode:            simulationCompany,
			OrderNo:                refundableOrders[index],
			CustomerRef:            fmt.Sprintf("RTL-CUST-%04d", 1+s.rng.Intn(80)),
			RefundedAt:             refundDate,
			CashAccountID:          s.pickCashAccount(),
			ContraRevenueAccountID: s.account("5211"),
			CurrencyCode:           "VND",
			Amount:                 amount,
			SourceRef:              fmt.Sprintf("POS-REFUND-%04d", index+1),
		}
		if _, err := s.submitCapture(ctx, financedomain.BusinessLineRetail, financedomain.OperationRetailCaptureRefund, payload); err != nil {
			return nil, fmt.Errorf("seed retail refund %d: %w", index+1, err)
		}
		s.stats.retailRefunds++
	}

	for index := 0; index < 60; index++ {
		heldAt := randomDateInRange(s.rng, reportFrom, reportTo.AddDate(0, 0, -10))
		amount := moneyFromInt(300_000 + s.rng.Intn(900_000))
		orderNo := fmt.Sprintf("RENT-%04d", index+1)
		payload := financeusecase.CaptureRentalDepositRequest{
			CompanyCode:      simulationCompany,
			RentalOrderID:    orderNo,
			CustomerRef:      fmt.Sprintf("RENT-CUST-%04d", 1+s.rng.Intn(60)),
			HeldAt:           heldAt,
			CashAccountID:    s.pickCashAccount(),
			HoldingAccountID: s.account("3388"),
			CurrencyCode:     "VND",
			Amount:           amount,
			SourceRef:        fmt.Sprintf("RENT-HOLD-%04d", index+1),
		}
		if _, err := s.submitCapture(ctx, financedomain.BusinessLineRental, financedomain.OperationRentalCaptureDeposit, payload); err != nil {
			return nil, fmt.Errorf("seed rental capture %d: %w", index+1, err)
		}
		s.stats.rentalCaptures++

		damageAmount := ledgerdomain.MustParseMoney("0.0000")
		if index < 20 {
			damageAmount = percentageOf(amount, 25+s.rng.Intn(50))
			if damageAmount.Cmp(amount) > 0 {
				damageAmount = amount
			}
		}
		releaseAt := heldAt.AddDate(0, 0, 3+s.rng.Intn(20))
		if releaseAt.After(reportTo) {
			releaseAt = reportTo
		}
		release := financeusecase.ReleaseRentalDepositRequest{
			CompanyCode:            simulationCompany,
			RentalOrderID:          orderNo,
			ReleasedAt:             releaseAt,
			CashAccountID:          payload.CashAccountID,
			DamageAmount:           damageAmount,
			DamageRevenueAccountID: s.account("711"),
			SourceRef:              fmt.Sprintf("RENT-SETTLE-%04d", index+1),
		}
		if _, err := s.submitCapture(ctx, financedomain.BusinessLineRental, financedomain.OperationRentalReleaseDeposit, release); err != nil {
			return nil, fmt.Errorf("seed rental settlement %d: %w", index+1, err)
		}
		s.stats.rentalSettlements++
	}

	return saleEntries, nil
}

func (s *simulator) seedOperatingExpenses(ctx context.Context) error {
	monthCursor := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for !monthCursor.After(reportTo) {
		expenses := []struct {
			costCenter string
			accountID  string
			amount     ledgerdomain.Money
			label      string
		}{
			{costCenter: "saas", accountID: s.account("6422"), amount: moneyFromInt(4_500_000 + s.rng.Intn(2_200_000)), label: "Van hanh SaaS"},
			{costCenter: "dojo", accountID: s.account("6421"), amount: moneyFromInt(6_000_000 + s.rng.Intn(3_000_000)), label: "Quan ly vo duong"},
			{costCenter: "retail", accountID: s.account("6422"), amount: moneyFromInt(3_200_000 + s.rng.Intn(1_800_000)), label: "Van hanh retail"},
		}
		for _, expense := range expenses {
			if err := s.postOperatingExpense(ctx, monthCursor, expense.costCenter, expense.accountID, expense.amount, expense.label); err != nil {
				return fmt.Errorf("seed operating expense %s %s: %w", expense.costCenter, monthCursor.Format("2006-01"), err)
			}
			s.stats.expenseEntries++
		}
		monthCursor = monthCursor.AddDate(0, 1, 0)
	}
	return nil
}

func (s *simulator) seedVoidCases(ctx context.Context, saleEntries []string) error {
	for index := 0; index < 5 && index < len(saleEntries); index++ {
		if _, err := s.voidUC.VoidEntry(ctx, saleEntries[index], fmt.Sprintf("Simulation void %d", index+1), "finance-sim"); err != nil {
			return fmt.Errorf("void retail sale %d: %w", index+1, err)
		}
		s.stats.voids++
	}
	return nil
}

func (s *simulator) submitCapture(ctx context.Context, line financedomain.BusinessLine, operation financedomain.CaptureOperation, payload any) (*financedomain.CaptureResult, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	s.stats.baseOps++
	key := fmt.Sprintf("sim-%04d", s.stats.baseOps)
	result, err := s.captureUC.Capture(ctx, financedomain.CaptureRequest{
		IdempotencyKey: key,
		BusinessLine:   line,
		Operation:      operation,
		Payload:        raw,
	})
	if err != nil {
		return nil, err
	}

	if s.stats.baseOps%10 == 0 {
		s.stats.duplicatesAttempted++
		replayed, err := s.captureUC.Capture(ctx, financedomain.CaptureRequest{
			IdempotencyKey: key,
			BusinessLine:   line,
			Operation:      operation,
			Payload:        raw,
		})
		if err != nil {
			return nil, fmt.Errorf("duplicate idempotency replay failed: %w", err)
		}
		if replayed.Replay {
			s.stats.duplicatesBlocked++
		}
	}
	return result, nil
}

func (s *simulator) postDirectCost(ctx context.Context, postingDate time.Time, costCenter string, expenseAccountID string, cashAccountID string, amount ledgerdomain.Money, description string) error {
	if !amount.IsPositive() {
		return nil
	}
	_, err := s.postEntryUC.PostEntry(ctx, ledgerusecase.PostEntryRequest{
		VoucherType:  "PC",
		CompanyCode:  simulationCompany,
		SourceModule: "simulation",
		Description:  description,
		CurrencyCode: "VND",
		PostingDate:  postingDate,
		Metadata: map[string]any{
			"business_line": costCenter,
			"cost_center":   costCenter,
			"flow_type":     "direct_cost",
		},
		Items: []ledgerusecase.PostEntryItemRequest{
			{AccountID: expenseAccountID, Side: "debit", Amount: amount},
			{AccountID: cashAccountID, Side: "credit", Amount: amount},
		},
	})
	return err
}

func (s *simulator) postOperatingExpense(ctx context.Context, postingDate time.Time, costCenter string, expenseAccountID string, amount ledgerdomain.Money, label string) error {
	_, err := s.postEntryUC.PostEntry(ctx, ledgerusecase.PostEntryRequest{
		VoucherType:  "PC",
		CompanyCode:  simulationCompany,
		SourceModule: "simulation",
		Description:  label,
		CurrencyCode: "VND",
		PostingDate:  postingDate,
		Metadata: map[string]any{
			"business_line": costCenter,
			"cost_center":   costCenter,
			"flow_type":     "operating_expense",
		},
		Items: []ledgerusecase.PostEntryItemRequest{
			{AccountID: expenseAccountID, Side: "debit", Amount: amount},
			{AccountID: s.pickBankAccount(), Side: "credit", Amount: amount},
		},
	})
	return err
}

func (s *simulator) collectSnapshots(ctx context.Context) error {
	if err := s.collectTrialBalance(ctx); err != nil {
		return err
	}
	if err := s.collectPL(ctx); err != nil {
		return err
	}
	if err := s.collectGrossProfit(ctx); err != nil {
		return err
	}
	if err := s.collectRevenueStream(ctx); err != nil {
		return err
	}
	if err := s.collectCashRunway(ctx); err != nil {
		return err
	}
	if err := s.collectPartitions(ctx); err != nil {
		return err
	}
	s.buildConsistencyLog()
	return nil
}

func (s *simulator) collectTrialBalance(ctx context.Context) error {
	query, err := os.ReadFile(filepath.Join("sql", "reports", "trial_balance.sql"))
	if err != nil {
		return fmt.Errorf("read trial balance sql: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, string(query), simulationCompany, reportFrom.Format(reportDateLayout), reportTo.Format(reportDateLayout))
	if err != nil {
		return fmt.Errorf("query trial balance: %w", err)
	}
	defer rows.Close()

	tbRows := make([]trialBalanceRow, 0, 32)
	summaryRows := make([]ledgerusecase.TrialBalanceRow, 0, 32)
	for rows.Next() {
		var (
			row        trialBalanceRow
			openingRaw string
			debitRaw   string
			creditRaw  string
			closingRaw string
		)
		if err := rows.Scan(&row.AccountCode, &row.AccountName, &row.NormalSide, &openingRaw, &debitRaw, &creditRaw, &closingRaw); err != nil {
			return fmt.Errorf("scan trial balance row: %w", err)
		}
		row.OpeningBalance = openingRaw
		row.PeriodDebit = debitRaw
		row.PeriodCredit = creditRaw
		row.ClosingBalance = closingRaw
		tbRows = append(tbRows, row)

		opening, _ := ledgerdomain.ParseMoney(openingRaw)
		debit, _ := ledgerdomain.ParseMoney(debitRaw)
		credit, _ := ledgerdomain.ParseMoney(creditRaw)
		closing, _ := ledgerdomain.ParseMoney(closingRaw)
		summaryRows = append(summaryRows, ledgerusecase.TrialBalanceRow{
			AccountCode:    row.AccountCode,
			AccountName:    row.AccountName,
			NormalSide:     ledgerdomain.Side(row.NormalSide),
			OpeningBalance: opening,
			PeriodDebit:    debit,
			PeriodCredit:   credit,
			ClosingBalance: closing,
		})
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate trial balance rows: %w", err)
	}
	if _, err := ledgerusecase.SummarizeTrialBalance(summaryRows); err != nil {
		return fmt.Errorf("validate trial balance: %w", err)
	}
	sort.Slice(tbRows, func(i, j int) bool { return tbRows[i].AccountCode < tbRows[j].AccountCode })
	s.stats.trialBalance = pickTrialBalanceSnapshot(tbRows)
	return nil
}

func (s *simulator) collectPL(ctx context.Context) error {
	query, err := os.ReadFile(filepath.Join("sql", "reports", "profit_and_loss_b02_dn.sql"))
	if err != nil {
		return fmt.Errorf("read p&l sql: %w", err)
	}
	rows, err := s.db.QueryContext(ctx, string(query), simulationCompany, reportFrom.Format(reportDateLayout), reportTo.Format(reportDateLayout))
	if err != nil {
		return fmt.Errorf("query p&l: %w", err)
	}
	defer rows.Close()

	result := make([]reportRow, 0, 16)
	for rows.Next() {
		var row reportRow
		if err := rows.Scan(&row.LineCode, &row.LineName, &row.Amount); err != nil {
			return fmt.Errorf("scan p&l row: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate p&l rows: %w", err)
	}
	s.stats.pnlRows = result
	return nil
}

func (s *simulator) collectGrossProfit(ctx context.Context) error {
	query, err := os.ReadFile(filepath.Join("sql", "reports", "gross_profit_by_cost_center.sql"))
	if err != nil {
		return fmt.Errorf("read gross profit sql: %w", err)
	}
	rows, err := s.db.QueryContext(ctx, string(query), simulationCompany, reportFrom.Format(reportDateLayout), reportTo.Format(reportDateLayout))
	if err != nil {
		return fmt.Errorf("query gross profit by cost center: %w", err)
	}
	defer rows.Close()

	result := make([]grossProfitRow, 0, 8)
	for rows.Next() {
		var row grossProfitRow
		var financialIncome string
		if err := rows.Scan(&row.CostCenter, &row.GrossRevenue, &row.RevenueDeductions, &financialIncome, &row.OtherIncome, &row.CostOfGoodsSold, &row.GrossProfit); err != nil {
			return fmt.Errorf("scan gross profit row: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate gross profit rows: %w", err)
	}
	s.stats.grossRows = filterGrossProfitSnapshot(result)
	return nil
}

func (s *simulator) collectRevenueStream(ctx context.Context) error {
	points, err := s.analyticsService.RevenueStream(ctx, simulationCompany, reportFrom, reportTo)
	if err != nil {
		return fmt.Errorf("load revenue stream: %w", err)
	}
	rows := make([]reportAmountRow, 0, len(points))
	for _, point := range points {
		rows = append(rows, reportAmountRow{
			Label:  point.CostCenter,
			Amount: point.Amount.String(),
		})
	}
	s.stats.revenueStream = rows
	return nil
}

func (s *simulator) collectCashRunway(ctx context.Context) error {
	runway, err := s.analyticsService.CashRunway(ctx, simulationCompany, time.Date(2026, 3, 21, 0, 0, 0, 0, time.UTC), 3)
	if err != nil {
		return fmt.Errorf("load cash runway: %w", err)
	}
	snapshot := analyticsSnapshot{
		CurrentCash:        runway.CurrentCash.String(),
		AverageMonthlyBurn: runway.AverageMonthlyBurn.String(),
		Months:             make([]analyticsMonth, 0, len(runway.Months)),
	}
	for _, month := range runway.Months {
		snapshot.Months = append(snapshot.Months, analyticsMonth{
			MonthLabel:       month.MonthLabel,
			OpeningCash:      month.OpeningCash.String(),
			ContractedInflow: month.ContractedInflow.String(),
			ProjectedBurn:    month.ProjectedBurn.String(),
			ProjectedEnding:  month.ProjectedEnding.String(),
		})
	}
	s.stats.cashRunway = snapshot
	return nil
}

func (s *simulator) collectPartitions(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `
SELECT tableoid::regclass::text AS partition_name, COUNT(*)::int
FROM journal_items
WHERE company_code = $1
GROUP BY tableoid::regclass::text
ORDER BY partition_name`,
		simulationCompany,
	)
	if err != nil {
		return fmt.Errorf("query partition counts: %w", err)
	}
	defer rows.Close()

	result := make([]partitionRow, 0, 8)
	for rows.Next() {
		var row partitionRow
		if err := rows.Scan(&row.PartitionName, &row.RowCount); err != nil {
			return fmt.Errorf("scan partition row: %w", err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate partition rows: %w", err)
	}
	s.stats.partitionCounts = result
	return nil
}

func (s *simulator) buildConsistencyLog() {
	s.stats.consistency = []string{
		fmt.Sprintf("Base capture operations seeded: %d", s.stats.baseOps),
		fmt.Sprintf("Idempotency duplicates attempted: %d", s.stats.duplicatesAttempted),
		fmt.Sprintf("Idempotency duplicates blocked as replay: %d", s.stats.duplicatesBlocked),
		fmt.Sprintf("Void cases executed: %d", s.stats.voids),
		fmt.Sprintf("Refund cases executed: %d", s.stats.retailRefunds),
		fmt.Sprintf("SaaS recognitions posted through 2026-03-31: %d", s.stats.saasRecognitions),
		"Trial balance validation: total debit equals total credit and closing balance formula passed.",
	}
}

func (s *simulator) writeArtifacts() error {
	return os.WriteFile(filepath.Join("docs", "finance-simulation-snapshot.md"), []byte(s.renderMarkdown()), 0o644)
}

func (s *simulator) renderMarkdown() string {
	var b strings.Builder
	b.WriteString("# Finance Simulation Snapshot\n\n")
	b.WriteString(fmt.Sprintf("Company: `%s`\n\n", simulationCompany))
	b.WriteString("## Trial Balance Snapshot\n\n")
	renderTable(&b, []string{"Account", "Name", "Normal", "Opening", "Debit", "Credit", "Closing"}, func() [][]string {
		rows := make([][]string, 0, len(s.stats.trialBalance))
		for _, row := range s.stats.trialBalance {
			rows = append(rows, []string{row.AccountCode, row.AccountName, row.NormalSide, row.OpeningBalance, row.PeriodDebit, row.PeriodCredit, row.ClosingBalance})
		}
		return rows
	}())
	b.WriteString("\n## P&L B02-DN Snapshot\n\n")
	renderTable(&b, []string{"Line", "Description", "Amount"}, func() [][]string {
		rows := make([][]string, 0, len(s.stats.pnlRows))
		for _, row := range s.stats.pnlRows {
			rows = append(rows, []string{row.LineCode, row.LineName, row.Amount})
		}
		return rows
	}())
	b.WriteString("\n## Gross Profit By Cost Center\n\n")
	renderTable(&b, []string{"Cost Center", "Gross Revenue", "Deductions", "Other Income", "COGS", "Gross Profit"}, func() [][]string {
		rows := make([][]string, 0, len(s.stats.grossRows))
		for _, row := range s.stats.grossRows {
			rows = append(rows, []string{row.CostCenter, row.GrossRevenue, row.RevenueDeductions, row.OtherIncome, row.CostOfGoodsSold, row.GrossProfit})
		}
		return rows
	}())
	b.WriteString("\n## Revenue Stream JSON Snapshot\n\n")
	renderTable(&b, []string{"Cost Center", "Net Revenue"}, func() [][]string {
		rows := make([][]string, 0, len(s.stats.revenueStream))
		for _, row := range s.stats.revenueStream {
			rows = append(rows, []string{row.Label, row.Amount})
		}
		return rows
	}())
	b.WriteString("\n## Cash Runway Snapshot\n\n")
	b.WriteString(fmt.Sprintf("Current cash: `%s`\n\n", s.stats.cashRunway.CurrentCash))
	b.WriteString(fmt.Sprintf("Average monthly burn: `%s`\n\n", s.stats.cashRunway.AverageMonthlyBurn))
	renderTable(&b, []string{"Month", "Opening Cash", "Contracted Inflow", "Projected Burn", "Projected Ending"}, func() [][]string {
		rows := make([][]string, 0, len(s.stats.cashRunway.Months))
		for _, month := range s.stats.cashRunway.Months {
			rows = append(rows, []string{month.MonthLabel, month.OpeningCash, month.ContractedInflow, month.ProjectedBurn, month.ProjectedEnding})
		}
		return rows
	}())
	b.WriteString("\n## Partition Distribution\n\n")
	renderTable(&b, []string{"Partition", "Rows"}, func() [][]string {
		rows := make([][]string, 0, len(s.stats.partitionCounts))
		for _, row := range s.stats.partitionCounts {
			rows = append(rows, []string{row.PartitionName, fmt.Sprintf("%d", row.RowCount)})
		}
		return rows
	}())
	b.WriteString("\n## Consistency Log\n\n")
	for _, line := range s.stats.consistency {
		b.WriteString("- ")
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String()
}

func (s *simulator) ensureSimulationAccounts(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM accounts WHERE company_code = $1`, simulationCompany).Scan(&count); err != nil {
		return fmt.Errorf("count simulation accounts: %w", err)
	}
	if count > 0 {
		return nil
	}

	rows, err := s.db.QueryContext(ctx, `
SELECT
    id,
    code,
    name,
    parent_id,
    account_type,
    normal_side,
    is_postable,
    is_active,
    description,
    COALESCE(metadata::text, '{}')
FROM accounts
WHERE company_code = $1
ORDER BY depth, code`, baseCompany)
	if err != nil {
		return fmt.Errorf("load base accounts: %w", err)
	}
	defer rows.Close()

	type baseAccount struct {
		oldID       string
		code        string
		name        string
		parentID    sql.NullString
		accountType string
		normalSide  string
		isPostable  bool
		isActive    bool
		description sql.NullString
		metadata    string
	}

	baseAccounts := make([]baseAccount, 0, 32)
	for rows.Next() {
		var item baseAccount
		if err := rows.Scan(
			&item.oldID,
			&item.code,
			&item.name,
			&item.parentID,
			&item.accountType,
			&item.normalSide,
			&item.isPostable,
			&item.isActive,
			&item.description,
			&item.metadata,
		); err != nil {
			return fmt.Errorf("scan base account: %w", err)
		}
		baseAccounts = append(baseAccounts, item)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate base accounts: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin clone accounts tx: %w", err)
	}
	defer tx.Rollback()

	idMap := make(map[string]string, len(baseAccounts))
	now := time.Now().UTC()
	for _, account := range baseAccounts {
		newID := id.NewUUID()
		var parentID any
		if account.parentID.Valid {
			parentID = idMap[account.parentID.String]
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO accounts (
    id,
    company_code,
    code,
    name,
    parent_id,
    account_type,
    normal_side,
    is_postable,
    is_active,
    description,
    metadata,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, CAST($6 AS account_category), CAST($7 AS normal_side), $8, $9, NULLIF($10, ''), CAST($11 AS JSONB), $12, $13)`,
			newID,
			simulationCompany,
			account.code,
			account.name,
			parentID,
			account.accountType,
			account.normalSide,
			account.isPostable,
			account.isActive,
			account.description.String,
			account.metadata,
			now,
			now,
		); err != nil {
			return fmt.Errorf("insert simulation account %s: %w", account.code, err)
		}
		idMap[account.oldID] = newID
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit clone accounts: %w", err)
	}
	return nil
}

func (s *simulator) resetSimulationData(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := []struct {
		query string
		args  []any
	}{
		{`DELETE FROM bank_statement_lines WHERE company_code = $1`, []any{simulationCompany}},
		{`DELETE FROM rental_deposits WHERE company_code = $1`, []any{simulationCompany}},
		{`DELETE FROM dojo_receivables WHERE company_code = $1`, []any{simulationCompany}},
		{`DELETE FROM saas_revenue_schedules WHERE contract_id IN (SELECT id FROM saas_contracts WHERE company_code = $1)`, []any{simulationCompany}},
		{`DELETE FROM saas_contracts WHERE company_code = $1`, []any{simulationCompany}},
		{`DELETE FROM outbox_events WHERE aggregate_id IN (SELECT id FROM journal_entries WHERE company_code = $1)`, []any{simulationCompany}},
		{`DELETE FROM account_balances WHERE company_code = $1`, []any{simulationCompany}},
		{`DELETE FROM voucher_sequences WHERE company_code = $1`, []any{simulationCompany}},
		{`DELETE FROM journal_entries WHERE company_code = $1`, []any{simulationCompany}},
		{`DELETE FROM idempotency_keys WHERE scope = 'finance_capture' AND idempotency_key LIKE 'sim-%'`, nil},
	}

	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement.query, statement.args...); err != nil {
			return fmt.Errorf("reset simulation data with %q: %w", statement.query, err)
		}
	}
	return tx.Commit()
}

func (s *simulator) loadAccounts(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `SELECT code, id FROM accounts WHERE company_code = $1`, simulationCompany)
	if err != nil {
		return err
	}
	defer rows.Close()

	s.accounts = make(map[string]string, 32)
	for rows.Next() {
		var code string
		var accountID string
		if err := rows.Scan(&code, &accountID); err != nil {
			return err
		}
		s.accounts[code] = accountID
	}
	return rows.Err()
}

func (s *simulator) account(code string) string {
	accountID, ok := s.accounts[code]
	if !ok {
		panic("missing simulation account code: " + code)
	}
	return accountID
}

func (s *simulator) pickCashAccount() string {
	if s.rng.Intn(100) < 72 {
		return s.account("1121")
	}
	return s.account("1111")
}

func (s *simulator) pickBankAccount() string {
	return s.account("1121")
}

func executeSQLFile(ctx context.Context, db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(content)) == "" {
		return nil
	}
	_, err = db.ExecContext(ctx, string(content))
	return err
}

func tableExists(ctx context.Context, db *sql.DB, tableName string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, `
SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = 'public'
      AND table_name = $1
)`, tableName).Scan(&exists)
	return exists, err
}

func monthStart(value time.Time) time.Time {
	return time.Date(value.UTC().Year(), value.UTC().Month(), 1, 0, 0, 0, 0, time.UTC)
}

func randomDateInRange(rng *rand.Rand, from time.Time, to time.Time) time.Time {
	span := to.Sub(from)
	if span <= 0 {
		return from
	}
	return from.Add(time.Duration(rng.Int63n(int64(span))))
}

func moneyFromInt(value int) ledgerdomain.Money {
	return ledgerdomain.MustParseMoney(fmt.Sprintf("%d.0000", value))
}

func percentageOf(amount ledgerdomain.Money, percentage int) ledgerdomain.Money {
	return amount.MustMul(percentage).MustDiv(100)
}

func pickTrialBalanceSnapshot(rows []trialBalanceRow) []trialBalanceRow {
	wanted := map[string]struct{}{
		"1111": {},
		"1121": {},
		"131":  {},
		"3387": {},
		"3388": {},
		"5111": {},
		"5113": {},
		"5211": {},
		"632":  {},
		"6421": {},
		"6422": {},
		"711":  {},
	}
	result := make([]trialBalanceRow, 0, len(wanted))
	for _, row := range rows {
		if _, ok := wanted[row.AccountCode]; ok {
			result = append(result, row)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].AccountCode < result[j].AccountCode })
	return result
}

func filterGrossProfitSnapshot(rows []grossProfitRow) []grossProfitRow {
	filtered := make([]grossProfitRow, 0, len(rows))
	for _, row := range rows {
		if row.CostCenter == "saas" || row.CostCenter == "dojo" || row.CostCenter == "retail" || row.CostCenter == "rental" {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

func renderTable(builder *strings.Builder, headers []string, rows [][]string) {
	builder.WriteString("| ")
	builder.WriteString(strings.Join(headers, " | "))
	builder.WriteString(" |\n| ")
	separators := make([]string, len(headers))
	for index := range separators {
		separators[index] = "---"
	}
	builder.WriteString(strings.Join(separators, " | "))
	builder.WriteString(" |\n")
	for _, row := range rows {
		builder.WriteString("| ")
		builder.WriteString(strings.Join(row, " | "))
		builder.WriteString(" |\n")
	}
}
