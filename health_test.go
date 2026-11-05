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
		expected float64
	}

	failed := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsNotReady"
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

	degraded := []byte(`{
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

	succeeded := []byte(`{
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

	statusResponses := [][]byte{
		failed,
		succeeded,
		degraded,
	}

	var cases []testCase

	for i := range len(statusResponses) {
		var cluster Cluster

		cluster.Name = fmt.Sprintf("test-cluster-%v", i)

		err := json.Unmarshal(statusResponses[i], &cluster.Status)
		if err != nil {
			t.Errorf("failed unmarshal %s", err)
		}

		cases = append(cases, testCase{cluster.Name, cluster, 3.0})
	}

	for _, tc := range cases {
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
