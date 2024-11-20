package metrics

type MatchRules = []MetricMatcher

func GenerateMatchRules() MatchRules {
	return MatchRules{
		{
			Header: "",
			Labels: []Label{},
		},
	}
}
