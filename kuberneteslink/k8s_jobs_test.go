package kuberneteslink

import (
	"github.com/fredex42/gowait/filescanner"
	"testing"
	"time"
)

func TestEnvVarFromMap(t *testing.T) {
	valuesMap := map[string]string{"firstkey": "firstvalue", "secondkey": "secondvalue"}
	rec := filescanner.WatchRecord{Path: "/path/to/a", Filename: "file", StableIterations: 3, LastMtime: time.Now()}
	result := EnvVarFromMap(&valuesMap, &rec)
	if len(result) != 4 {
		t.Errorf("Expected 4 elements, got %d: %s", len(result), result)
	}

	if result[0].Name != "firstkey" || result[0].Value != "firstvalue" || result[0].ValueFrom != nil {
		t.Errorf("Unexpected result 1: %s", result[0])
	}
	if result[1].Name != "secondkey" || result[1].Value != "secondvalue" || result[1].ValueFrom != nil {
		t.Errorf("Unexpected result 2: %s", result[1])
	}
	if result[2].Name != "filename" || result[2].Value != "file" || result[2].ValueFrom != nil {
		t.Errorf("Unexpected result 4: %s", result[2])
	}
	if result[3].Name != "path" || result[3].Value != "/path/to/a" || result[3].ValueFrom != nil {
		t.Errorf("Unexpected result 3: %s", result[3])
	}

}
