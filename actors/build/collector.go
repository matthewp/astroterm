package build

import (
	"astroterm/info"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type BuildDataCollector struct {
	complete bool
	inData   bool
	forward  io.Writer
	raw      string
}

var buildDataStartStr = "ASTROTERM_BUILD_DATA:START"
var buildDataEndStr = "ASTROTERM_BUILD_DATA:END"

func NewBuildDataCollector(forward io.Writer) *BuildDataCollector {
	return &BuildDataCollector{
		complete: false,
		forward:  forward,
		inData:   false,
		raw:      "",
	}
}

func (bd *BuildDataCollector) Write(p []byte) (n int, err error) {
	if bd.complete {
		return bd.forward.Write(p)
	}

	str := string(p)

	if bd.inData {
		idx := strings.Index(str, buildDataEndStr)
		if idx != -1 {
			bd.raw += str[:idx]
			bd.complete = true
		} else {
			bd.raw += str
		}
	} else {
		idx := strings.Index(str, buildDataStartStr)
		if idx != -1 {
			bd.raw = str[idx:]
			bd.inData = true
		} else {
			return bd.forward.Write(p)
		}
	}

	return len(p), nil
}

func (bd *BuildDataCollector) Raw() string {
	return bd.raw
}

func (bd *BuildDataCollector) Json() (*info.BuildData, error) {
	if bd.raw == "" {
		return nil, errors.New("no build data collected")
	}

	var data info.BuildData
	err := json.Unmarshal([]byte(bd.raw), &data)
	return &data, err
}
