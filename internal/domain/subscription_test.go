package domain

import (
	"testing"
	"time"
)

func TestSubscription_IsActive(t *testing.T) {
	tests := []struct {
		name      string
		startDate time.Time
		endDate   *time.Time
		want      bool
	}{
		{
			name:      "активная подписка",
			startDate: time.Now().Add(-24 * time.Hour),
			endDate:   &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
			want:      true,
		},
		{
			name:      "просроченная подписка",
			startDate: time.Now().Add(-48 * time.Hour),
			endDate:   &[]time.Time{time.Now().Add(-24 * time.Hour)}[0],
			want:      false,
		},
		{
			name:      "бессрочная подписка",
			startDate: time.Now().Add(-24 * time.Hour),
			endDate:   nil,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				StartDate: tt.startDate,
				EndDate:   tt.endDate,
			}
			if got := sub.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
