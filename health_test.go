package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIsHappy(t *testing.T) {

	type testCase struct {
		name          string
		cluster       Cluster
		data          []byte
		thresholdLow  float64
		thresholdHigh float64
	}

	// cluster not ready should dominate 0.0 failure
	failed_no_backup := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsNotReady"
	    },
	    {
	      "reason": "ContinuousArchivingFailing"
	    },
	    {
	      "reason": "LastBackupFailed"
	    }
	  ]
	}
    }`)

	// cluster not ready should dominate 0.0 failure
	// backup succeeded doesn't show under conditions if cluster down archiving down
	failed_with_backup := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsNotReady"
	    },
	    {
	      "reason": "ContinuousArchivingFailing"
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

	// cluster ready should dominate but degraded
	succeeded_no_backup_degraded := []byte(`{
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

	// cluster ready should dominate but degraded
	succeeded_no_archiving_degraded := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsReady"
	    },
	    {
	      "reason": "ContinuousArchivingFailing"
	    },
	    {
	      "reason": "LastBackupSucceeded"
	    }
	  ]
	}
    }`)

	cases := []testCase{
		{name: "failed-no-backup",
			data:          failed_no_backup,
			thresholdLow:  0.0,
			thresholdHigh: 1.5,
		}, // failed
		{name: "failed-with-backup",
			data:          failed_with_backup,
			thresholdLow:  0.0,
			thresholdHigh: 1.5,
		}, // failed
		{name: "succeeded-with-backup",
			data:          succeeded_with_backup,
			thresholdLow:  2.5,
			thresholdHigh: 3.0,
		}, // succeeded
		{name: "succeeded-no-backup",
			data:          succeeded_no_backup_degraded,
			thresholdLow:  1.5,
			thresholdHigh: 2.5,
		}, // degraded
		{name: "succeeded-no-archiving",
			data:          succeeded_no_archiving_degraded,
			thresholdLow:  1.5,
			thresholdHigh: 2.5,
		}, // degraded
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

			if tc.cluster.OverallHealth <= tc.thresholdHigh ||
				tc.cluster.OverallHealth >= tc.thresholdLow {
				t.Errorf("expected: %v < %v <= %v", tc.thresholdLow,
					tc.cluster.OverallHealth,
					tc.thresholdHigh)
			}
		})

	}
}
