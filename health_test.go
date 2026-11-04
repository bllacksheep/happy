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
	      "reason": "ClusterIsReady",
	      "status": "True",
	      "type": "Ready"
	    },
	    {
	      "lastTransitionTime": "2026-10-29T09:39:52Z",
	      "message": "Continuous archiving is working",
	      "reason": "ContinuousArchivingSuccess",
	      "status": "True",
	      "type": "ContinuousArchiving"
	    },
	    {
	      "lastTransitionTime": "2026-10-29T18:12:06Z",
	      "message": "exit status 1",
	      "reason": "LastBackupFailed",
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

	if cluster.Happy == false {
		t.Errorf("not happy %s", cluster.Status)
	}

}
