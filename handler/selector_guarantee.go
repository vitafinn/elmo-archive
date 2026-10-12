package handler

// Without historical targets, later results must remain unknown. Explicit user
// annotations take precedence; only the initial and post-loss guarantees infer wins.
type selectorGuarantees map[int64]string

func (history selectorGuarantees) observe(bannerID int64, annotation string) (string, bool) {
	previous, seen := history[bannerID]
	outcome := annotation
	if outcome == "" || outcome == "auto" {
		switch {
		case !seen:
			outcome = "initial"
		case previous == "loss":
			outcome = "guaranteed"
		default:
			outcome = "unknown"
		}
	}
	history[bannerID] = outcome
	return outcome, outcome == "win" || outcome == "loss"
}
