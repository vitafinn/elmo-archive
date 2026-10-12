package handler

import "testing"

func TestSelectorUnknownAndExplicitOutcomes(t *testing.T) {
	h := make(selectorGuarantees)
	cases := []struct {
		banner           int64
		annotation, want string
		eligible         bool
	}{
		{1, "", "initial", false}, {1, "", "unknown", false}, {1, "loss", "loss", true},
		{1, "", "guaranteed", false}, {1, "", "unknown", false}, {1, "win", "win", true},
		{2, "", "initial", false}, {3, "unknown", "unknown", false}, {3, "", "unknown", false},
		{1, "auto", "unknown", false},
	}
	for _, tc := range cases {
		got, eligible := h.observe(tc.banner, tc.annotation)
		if got != tc.want || eligible != tc.eligible {
			t.Fatalf("%+v: %s %v", tc, got, eligible)
		}
	}
}
