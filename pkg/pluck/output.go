package pluck

import (
	// default
	"encoding/json"

	// external
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ResultJSON returns the result, formatted as JSON.
// If their are no results, it returns an empty string.
func (p *Plucker) ResultJSON(indent ...bool) string {
	totalResults := 0
	for key := range p.result {
		b, _ := json.Marshal(p.result[key])
		totalResults += len(b)
	}
	if totalResults == len(p.result)*2 { // results == 2 because its just []
		return ""
	}
	var err error
	var resultJSON []byte
	if len(indent) > 0 && indent[0] {
		resultJSON, err = json.MarshalIndent(p.result, "", "    ")
	} else {
		resultJSON, err = json.Marshal(p.result)
	}
	if err != nil {
		log.Error(errors.Wrap(err, "result marshalling failed"))
	}
	return string(resultJSON)
}
