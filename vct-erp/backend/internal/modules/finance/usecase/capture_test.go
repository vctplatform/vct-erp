package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	financedomain "vct-platform/backend/internal/modules/finance/domain"
)

func TestCaptureUseCaseReplaysCompletedResponse(t *testing.T) {
	expected := financedomain.CaptureResult{
		BusinessLine: financedomain.BusinessLineSaaS,
		Operation:    financedomain.OperationSaaSCaptureAnnualContract,
		ResourceID:   "contract-1",
		Payload:      json.RawMessage(`{"status":"ok"}`),
	}
	rawExpected, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("marshal expected result: %v", err)
	}

	repo := &fakeIdempotencyRepo{
		reservation: financedomain.IdempotencyReservation{
			Status:          financedomain.IdempotencyStatusReplay,
			ResponsePayload: rawExpected,
		},
	}
	uc := NewCaptureUseCase(repo, nil, nil, nil)

	result, err := uc.Capture(context.Background(), financedomain.CaptureRequest{
		IdempotencyKey: "idem-1",
		BusinessLine:   financedomain.BusinessLineSaaS,
		Operation:      financedomain.OperationSaaSCaptureAnnualContract,
		Payload:        json.RawMessage(`{"contract_no":"S-001"}`),
	})
	if err != nil {
		t.Fatalf("Capture returned error: %v", err)
	}
	if !result.Replay {
		t.Fatal("expected replay flag to be true")
	}
	if result.ResourceID != "contract-1" {
		t.Fatalf("unexpected resource id: %s", result.ResourceID)
	}
}

type fakeIdempotencyRepo struct {
	reservation financedomain.IdempotencyReservation
}

func (f *fakeIdempotencyRepo) Reserve(_ context.Context, _ string, _ string, _ string, _ time.Time) (financedomain.IdempotencyReservation, error) {
	return f.reservation, nil
}

func (f *fakeIdempotencyRepo) Complete(_ context.Context, _ string, _ string, _ []byte, _ string, _ time.Time) error {
	return nil
}

func (f *fakeIdempotencyRepo) Fail(_ context.Context, _ string, _ string, _ string, _ time.Time) error {
	return nil
}
