package metrics

type MatchRules = []MetricMatcher

func GenerateMatchRules() MatchRules {
	return MatchRules{
		{
			header:      "",
			neededLabel: []string{},
		},
	}
}
