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
		expected int
	}

	failed := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsNotReady",
	    },
	    {
	      "reason": "ContinuousArchivingSucceeded",
	    },
	    {
	      "reason": "LastBackupSucceeded",
	    }
	  ]
	}
    }`)

	degraded := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsReady",
	    },
	    {
	      "reason": "ContinuousArchivingSucceeded",
	    },
	    {
	      "reason": "LastBackupFailed",
	    }
	  ]
	}
    }`)

	succeeded := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "reason": "ClusterIsReady",
	    },
	    {
	      "reason": "ContinuousArchivingSucceeded",
	    },
	    {
	      "reason": "LastBackupSucceeded",
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

		cluster.Name = fmt.Sprintf("test-%v", i)

		err := json.Unmarshal(statusResponses[i], &cluster.Status)
		if err != nil {
			t.Errorf("failed unmarshal %s", err)
		}

		cases = append(cases, testCase{cluster.Name, cluster, i})
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.cluster.IsHappy()

			if tc.cluster.Happy < 1 {
				t.Errorf("%s polled x%v, cluster happy x%v", tc.cluster.Name,
					tc.cluster.Polled,
					tc.cluster.Happy)
			}
		})
	}
}
