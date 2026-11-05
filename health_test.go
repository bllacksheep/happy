package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIsHappy(t *testing.T) {

	// dummy case for now
	case1 := []byte(`{
	  "status": {
	  "conditions": [
	    {
	      "lastTransitionTime": "2026-10-29T09:40:46Z",
	      "message": "Cluster is Ready",
	      "reason": "ClusterIsNotReady",
	      "status": "True",
	      "type": "Ready"
	    },
	    {
	      "lastTransitionTime": "2026-10-29T09:39:52Z",
	      "message": "Continuous archiving is working",
	      "reason": "ContinuousArchivingSucceeded",
	      "status": "True",
	      "type": "ContinuousArchiving"
	    },
	    {
	      "lastTransitionTime": "2026-10-29T18:12:06Z",
	      "message": "exit status 1",
	      "reason": "LastBackupSucceeded",
	      "status": "False",
	      "type": "LastBackupSucceeded"
	    }
	  ]
	}
        }`)

	// named cluster can test when status is nil
	cluster := Cluster{
		Name: "test-case1-cluster",
	}

	err := json.Unmarshal(case1, &cluster.Status)
	if err != nil {
		fmt.Println("failed unmarshal", err)
	}

	cluster.IsHappy()

	if cluster.Happy < 1 {
		t.Errorf("%s polled x%v, cluster happy x%v", cluster.Name,
			cluster.Polled,
			cluster.Happy)

	}
}
