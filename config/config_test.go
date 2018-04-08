package config

import (
	"io/ioutil"
	"regexp"
	"testing"
)

func TestParseConfig(t *testing.T) {

	b, err := ioutil.ReadFile("../samples/config.json")
	if err != nil {
		t.Error(err)
	}

	configTests := []struct {
		name            string
		configBytes     []byte
		ignoreMissingDB bool
		correct         bool
		errMsgRegexp    *regexp.Regexp
		nbColl          int
	}{
		{
			name:            "samples/config.json",
			configBytes:     b,
			ignoreMissingDB: false,
			correct:         true,
			errMsgRegexp:    nil,
			nbColl:          2,
		},
		{
			name: "invalid content",
			configBytes: []byte(`[{
				"database": "datagen_it_test", 
				"collection": "test",
				"count": 1000,
				"content": { "k": invalid }
				}]`),
			ignoreMissingDB: false,
			correct:         false,
			errMsgRegexp:    regexp.MustCompile("^Error in configuration file: object / array / Date badly formatted: \n\n\t\t.*"),
			nbColl:          0,
		},
		{
			name: "missing database field",
			configBytes: []byte(`[{
				"collection": "test",
				"count": 1000,
				"content": {}
				}]`),
			ignoreMissingDB: false,
			correct:         false,
			errMsgRegexp:    regexp.MustCompile("^Error in configuration file: \n\t'collection' and 'database' fields can't be empty.*"),
			nbColl:          0,
		},
		{
			name: "count > 0",
			configBytes: []byte(`[{
				"database": "datagen_it_test", 
				"collection": "test",
				"count": 0,
				"content": {}
				}]`),
			ignoreMissingDB: false,
			correct:         false,
			errMsgRegexp:    regexp.MustCompile("^Error in configuration file: \n\tfor collection.*"),
			nbColl:          0,
		},
	}

	for _, tt := range configTests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseConfig(tt.configBytes, tt.ignoreMissingDB)
			if tt.correct {
				if err != nil {
					t.Errorf("expected no error for config %s: %v", tt.configBytes, err)
				}
				if len(c) != tt.nbColl {
					t.Errorf("expected %d coll but got %d", tt.nbColl, len(c))
				}
			} else {
				if err == nil {
					t.Errorf("expected an error for config %s", tt.configBytes)
				}
				if !tt.errMsgRegexp.MatchString(err.Error()) {
					t.Errorf("error message should match %s, but was %v", tt.errMsgRegexp.String(), err)
				}
			}
		})
	}
}