package usecase

import "testing"

func TestGetMockDashboardDataShape(t *testing.T) {
	data := GetMockDashboardData()

	if data.DataMode != "mock" {
		t.Fatalf("expected data mode mock, got %s", data.DataMode)
	}
	if len(data.Cards) != 3 {
		t.Fatalf("expected 3 dashboard cards, got %d", len(data.Cards))
	}
	for _, card := range data.Cards {
		if len(card.ChartData) != 7 {
			t.Fatalf("expected 7 chart points for card %s, got %d", card.Key, len(card.ChartData))
		}
	}
	if len(data.RevenueMix) != 4 {
		t.Fatalf("expected 4 pie slices, got %d", len(data.RevenueMix))
	}
	if len(data.CashflowChart.XAxis) != 6 {
		t.Fatalf("expected 6 x-axis points, got %d", len(data.CashflowChart.XAxis))
	}
	if len(data.CashflowChart.Series) != 3 {
		t.Fatalf("expected 3 line series, got %d", len(data.CashflowChart.Series))
	}
	if len(data.RunwayProjection) != 6 {
		t.Fatalf("expected 6 runway projection points, got %d", len(data.RunwayProjection))
	}
}
