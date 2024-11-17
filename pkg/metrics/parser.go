package metrics

import (
	"errors"
	"strings"

	"github.com/elastic/go-grok"
	log "github.com/sirupsen/logrus"
)

var (
	MetricsPatterns = map[string]string{
		"METRIC_HEAD": "^%{WORD:metric_name}{",
		"LAST_LABEL":  `(%{DATA:last_label_name}="%{DATA:last_label_value}")`,
		"NUMBER_SCI":  "%{NUMBER}(e%{NUMBER})?",
		"METRIC_TAIL": `}%{SPACE}%{NUMBER_SCI:metric_value}$`,
		"NORMAL_EXPR": `^%{WORD:metric_name}{((namespace="%{DATA:namespace}")|(%{DATA:last_label_name}=\"%{DATA:last_label_value}\")(,|))*}%{SPACE}%{NUMBER}(e%{NUMBER})?$`,
	}
	COMMON_MATCH_GROK = grok.New()
)

const COMMON_MATCH = `^%{WORD:name}{((namespace="%{DATA:namespace}")|(%{DATA:last_label_name}="%{DATA:last_label_value}")(,|))*}%{SPACE}%{NUMBER}(e%{NUMBER})?$`

type MetricMatcher struct {
	MtName      string
	header      string
	neededLabel []string
	grok        *grok.Grok
}

func NewMetricMatcher(label_name string) *MetricMatcher {
	return &MetricMatcher{
		MtName:      label_name,
		grok:        grok.New(),
		neededLabel: make([]string, 0),
	}
}

// example: ^%{WORD:name}{((namespace="%{DATA:ns}")|(%{DATA:last_label_name}="%{DATA:last_label_value}")(,|))*}%{SPACE}%{NUMBER}(e%{NUMBER})?$
//          ^%{WORD:name}{
//                        (
//                         (namespace="%{DATA:ns}")
//                                                 |
//                                                  (%{DATA:last_label_name}="%{DATA:last_label_value}")
//                                                                                                      (,|)
//                                                                                                          )
//                                                                                                           *
//                                                                                                            }%{SPACE}%{NUMBER}(e%{NUMBER})?$

func (mt *MetricMatcher) Compile() error {
	if err := mt.grok.AddPatterns(MetricsPatterns); err != nil {
		return err
	}
	header := mt.header
	if mt.header == "" {
		header = `^` + mt.MtName + `{` // custom head
	}
	var body []string
	for _, label := range mt.neededLabel {
		body = append(body, "("+label+`=`+`"%{DATA:`+strings.ToUpper(label)+`}")`)
	}
	tail := `${METRIC_TAIL}`
	pattern := header + "(" + strings.Join(body, "|") + "|%{LAST_LABEL}(,|))*" + tail
	log.Debugln("generated pattern: ", pattern)
	return mt.grok.Compile(pattern, true)
}

func (mt *MetricMatcher) Match(target string) (res map[string]string, err error) {
	if !mt.grok.MatchString(target) {
		if len(target) > 100 {
			target = target[:99]
		}
		log.Debugf("not match: %s", target)
		if COMMON_MATCH_GROK.MatchString(target) {
			res, err = COMMON_MATCH_GROK.ParseString(target)
			if err != nil {
				return nil, errors.Join(
					errors.New("match failed, in common expr"),
					err,
				)
			}
			return res, errors.Join(
				errors.New("match failed, in custom expr, but success in common expr, check ret result to get more detail"),
				err,
			)
		}
	}
	return mt.grok.ParseString(target)
}

func (mt *MetricMatcher) SetHeader(header string) {
	mt.header = `^` + header + `{`
}

func (mt *MetricMatcher) AddLabel(label string) {
	mt.neededLabel = append(mt.neededLabel, label)
}

func init() {
	err := COMMON_MATCH_GROK.AddPatterns(MetricsPatterns)
	if err != nil {
		log.Fatalf("add patterns failed: %v", err)
	}
	err = COMMON_MATCH_GROK.Compile(COMMON_MATCH, true)
	if err != nil {
		log.Fatalf("compile pattern failed: %v", err)
	}
}
