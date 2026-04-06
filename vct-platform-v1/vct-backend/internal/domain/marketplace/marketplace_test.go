package marketplace

import "testing"

func TestNormalizeProductStatus(t *testing.T) {
	tests := []struct {
		name   string
		status ProductStatus
		stock  int
		want   ProductStatus
	}{
		{
			name:   "draft remains draft when stock is zero",
			status: ProductStatusDraft,
			stock:  0,
			want:   ProductStatusDraft,
		},
		{
			name:   "archived remains archived when stock is zero",
			status: ProductStatusArchived,
			stock:  0,
			want:   ProductStatusArchived,
		},
		{
			name:   "active becomes out of stock when inventory is empty",
			status: ProductStatusActive,
			stock:  0,
			want:   ProductStatusOutOfStock,
		},
		{
			name:   "out of stock returns to active when inventory is refilled",
			status: ProductStatusOutOfStock,
			stock:  3,
			want:   ProductStatusActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeProductStatus(tt.status, tt.stock)
			if got != tt.want {
				t.Fatalf("normalizeProductStatus(%q, %d) = %q, want %q", tt.status, tt.stock, got, tt.want)
			}
		})
	}
}
