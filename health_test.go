package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIsHappy(t *testing.T) {

	type testCase struct {
		name     string
		cluster  Cluster
		data     []byte
		expected float64
	}

	// cluster not ready should dominate 0.0 failure
	failed_no_backup := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsNotReady"
	    },
	    {
	      "reason": "ContinuousArchivingFailed"
	    },
	    {
	      "reason": "LastBackupFailed"
	    }
	  ]
	}
    }`)

	// cluster not ready should dominate 0.0 failure
	failed_with_backup := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsNotReady"
	    },
	    {
	      "reason": "ContinuousArchivingFailed"
	    },
	    {
	      "reason": "LastBackupSucceeded"
	    }
	  ]
	}
    }`)

	// cluster ready should dominate
	succeeded_with_backup := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsReady"
	    },
	    {
	      "reason": "ContinuousArchivingSucceeded"
	    },
	    {
	      "reason": "LastBackupSucceeded"
	    }
	  ]
	}
    }`)

	// cluster ready should dominate
	succeeded_no_archiving := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsReady"
	    },
	    {
	      "reason": "ContinuousArchivingFailed"
	    },
	    {
	      "reason": "LastBackupSucceeded"
	    }
	  ]
	}
    }`)

	// cluster ready should dominate
	succeeded_no_backup := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsReady"
	    },
	    {
	      "reason": "ContinuousArchivingSucceeded"
	    },
	    {
	      "reason": "LastBackupFailed"
	    }
	  ]
	}
    }`)

	cases := []testCase{
		{name: "failed-no-backup", data: failed_no_backup, expected: 0.0},             // failed
		{name: "failed-with-backup", data: failed_with_backup, expected: 0.0},         // failed
		{name: "succeeded-with-backup", data: succeeded_with_backup, expected: 3.0},   // succeeded
		{name: "succeeded-no-backup", data: succeeded_no_backup, expected: 1.0},       // degraded
		{name: "succeeded-no-archiving", data: succeeded_no_archiving, expected: 1.0}, // degraded
	}

	for id, tc := range cases {
		var cluster Cluster
		cluster.Name = fmt.Sprintf("test-cluster-%d", id)
		tc.name = fmt.Sprintf("%s-%s", cluster.Name, tc.name)

		err := json.Unmarshal(tc.data, &cluster.Status)
		if err != nil {
			t.Errorf("failed unmarshal %s", err)
		}

		tc.cluster = cluster

		t.Run(tc.name, func(t *testing.T) {
			tc.cluster.IsHappy()

			// 2.5 > success < 3
			// 1.5 > degraded < 2.5
			// 0 > failed < 1.5

			if tc.cluster.OverallHealth != tc.expected {
				t.Errorf("expected: %v, got: %v", tc.expected,
					tc.cluster.Happy)
			}
		})

	}
}
