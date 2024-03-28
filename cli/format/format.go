// Copyright 2018 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package format

import (
	"io"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-openapi/strfmt"

	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/alertmanager/matcher"
)

const DefaultDateFormat = "2006-01-02 15:04:05 MST"

var dateFormat *string

func InitFormatFlags(app *kingpin.Application) {
	dateFormat = app.Flag("date.format", "Format of date output").Default(DefaultDateFormat).String()
}

// Formatter needs to be implemented for each new output formatter.
type Formatter interface {
	SetOutput(io.Writer)
	FormatSilences([]models.GettableSilence) error
	FormatAlerts([]*models.GettableAlert) error
	FormatConfig(*models.AlertmanagerStatus) error
	FormatClusterStatus(status *models.ClusterStatus) error
}

// Formatters is a map of cli argument names to formatter interface object.
var Formatters = map[string]Formatter{}

func FormatDate(input strfmt.DateTime) string {
	return time.Time(input).Format(*dateFormat)
}

func labelsMatcher(m models.Matcher) *matcher.Matcher {
	var t matcher.MatchType
	// Support for older alertmanager releases, which did not support isEqual.
	if m.IsEqual == nil {
		isEqual := true
		m.IsEqual = &isEqual
	}
	switch {
	case !*m.IsRegex && *m.IsEqual:
		t = matcher.MatchEqual
	case !*m.IsRegex && !*m.IsEqual:
		t = matcher.MatchNotEqual
	case *m.IsRegex && *m.IsEqual:
		t = matcher.MatchRegexp
	case *m.IsRegex && !*m.IsEqual:
		t = matcher.MatchNotRegexp
	}

	return &matcher.Matcher{Type: t, Name: *m.Name, Value: *m.Value}
}
