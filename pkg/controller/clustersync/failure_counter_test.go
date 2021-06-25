package clustersync

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	testifyassert "github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	hiveinternal "github.com/openshift/hive/apis/hiveinternal/v1alpha1"
)

func assertEqual(t *testing.T, x, y interface{}) {
	testifyassert.Empty(t, cmp.Diff(x, y), "+actual, -expected")
}

func mkSyncStatus(name string, result hiveinternal.SyncSetResult) hiveinternal.SyncStatus {
	return hiveinternal.SyncStatus{
		Name:   name,
		Result: result,
	}
}

func mkClusterSync(ssStatuses []hiveinternal.SyncStatus, sssStatuses []hiveinternal.SyncStatus) hiveinternal.ClusterSync {
	cs := hiveinternal.ClusterSync{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "foo",
		},
	}
	cs.Status.SyncSets = ssStatuses
	cs.Status.SelectorSyncSets = sssStatuses
	return cs
}

func Test_countErrors(t *testing.T) {
	successStatus1 := mkSyncStatus("good1", hiveinternal.SuccessSyncSetResult)
	successStatus2 := mkSyncStatus("good2", hiveinternal.SuccessSyncSetResult)
	failStatus1 := mkSyncStatus("bad1", hiveinternal.FailureSyncSetResult)
	failStatus2 := mkSyncStatus("bad2", hiveinternal.FailureSyncSetResult)

	type args struct {
		cs *hiveinternal.ClusterSync
	}
	tests := []struct {
		name                string
		cs                  hiveinternal.ClusterSync
		expectedStatuses    map[string]*hiveinternal.SyncStatus
		expectedByNamespace map[string][]string
		wantErr             bool
	}{
		{
			name:                "Empty ClusterSync",
			cs:                  mkClusterSync([]hiveinternal.SyncStatus{}, []hiveinternal.SyncStatus{}),
			expectedStatuses:    map[string]*hiveinternal.SyncStatus{},
			expectedByNamespace: map[string][]string{},
		},
		{
			name:                "No erroring syncsets",
			cs:                  mkClusterSync([]hiveinternal.SyncStatus{successStatus1}, []hiveinternal.SyncStatus{successStatus2}),
			expectedStatuses:    map[string]*hiveinternal.SyncStatus{},
			expectedByNamespace: map[string][]string{},
		},
		{
			name: "Failing syncsets",
			cs: mkClusterSync(
				[]hiveinternal.SyncStatus{failStatus1, failStatus2},
				[]hiveinternal.SyncStatus{successStatus1, successStatus2},
			),
			expectedStatuses: map[string]*hiveinternal.SyncStatus{
				"foo/ss/bad1": &failStatus1,
				"foo/ss/bad2": &failStatus2,
			},
			expectedByNamespace: map[string][]string{"foo": {"foo/ss/bad1", "foo/ss/bad2"}},
		},
		{
			name: "Failing selectorsyncsets",
			cs: mkClusterSync(
				[]hiveinternal.SyncStatus{successStatus1, successStatus2},
				[]hiveinternal.SyncStatus{failStatus1, failStatus2},
			),
			expectedStatuses: map[string]*hiveinternal.SyncStatus{
				"foo/sss/bad1": &failStatus1,
				"foo/sss/bad2": &failStatus2,
			},
			expectedByNamespace: map[string][]string{"foo": {"foo/sss/bad1", "foo/sss/bad2"}},
		},
		{
			name: "One of each failing",
			cs: mkClusterSync(
				[]hiveinternal.SyncStatus{successStatus1, failStatus1},
				[]hiveinternal.SyncStatus{failStatus2, successStatus2},
			),
			expectedStatuses: map[string]*hiveinternal.SyncStatus{
				"foo/ss/bad1": &failStatus1,
				"foo/sss/bad2": &failStatus2,
			},
			expectedByNamespace: map[string][]string{"foo": {"foo/sss/bad2", "foo/ss/bad1"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// NOTE: By not cleaning out the globals between tests, we're exercising
			// idempotence to some extent as well.
			// TODO: Expect synchronization errors if running multiple test cases in parallel.
			if err := countErrors(&tt.cs, log); (err != nil) != tt.wantErr {
				t.Errorf("countErrors() error = %v, wantErr %v", err, tt.wantErr)
			}
			assertEqual(t, tt.expectedStatuses, failures.statuses)
			assertEqual(t, tt.expectedByNamespace, failures.byNamespace)
		})
	}
}
