package json_util

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

type StrOrUint64 struct {
	Value uint64
}

func (r *StrOrUint64) UnmarshalJSON(data []byte) error {
	sysLpStr := string(data)
	if strings.Contains(sysLpStr, "-nan") {
		r.Value = uint64(math.Inf(-1))
		return nil
	}
	if strings.Contains(sysLpStr, "nan") {
		r.Value = uint64(math.Inf(1))
		return nil
	}
	if strings.Contains(sysLpStr, "\"") {
		sysLpStr = strings.ReplaceAll(sysLpStr, "\"", "")
	}

	return json.Unmarshal([]byte(sysLpStr), &r.Value)
}

type StrOrFloat64 struct {
	Value float64
}

func (r *StrOrFloat64) UnmarshalJSON(data []byte) error {
	sysLpStr := string(data)
	if strings.Contains(sysLpStr, "-nan") {
		r.Value = math.Inf(-1)
		return nil
	}
	if strings.Contains(sysLpStr, "nan") {
		r.Value = math.Inf(1)
		return nil
	}
	if strings.Contains(sysLpStr, "\"") {
		sysLpStr = strings.ReplaceAll(sysLpStr, "\"", "")
		//if sysLpStr == "-" {
		//	r.Value = math.NaN()
		//	return nil
		//}
		var err error
		r.Value, err = strconv.ParseFloat(sysLpStr, 64)
		if err != nil {
			r.Value = math.NaN()
		    return nil
		}
	}

	return json.Unmarshal([]byte(sysLpStr), &r.Value)
}
