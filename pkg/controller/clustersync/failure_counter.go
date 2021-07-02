package clustersync

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"

	hiveinternal "github.com/openshift/hive/apis/hiveinternal/v1alpha1"
)

var failures struct {
	mutex sync.RWMutex
	// statuses is a map, keyed by "$namespace/[s]ss/$syncsetname", of pointers to SyncStatus
	statuses map[string]*hiveinternal.SyncStatus
	// byNamespace is a map, keyed by "$namespace", of `statuses` keys
	byNamespace map[string][]string
}

func init() {
	failures.statuses = make(map[string]*hiveinternal.SyncStatus)
	failures.byNamespace = make(map[string][]string)
}

// failureKey generates a unique string used to key into the `failures.statuses` map.
func key(ns, ssOrsss string, status *hiveinternal.SyncStatus) string {
	// The ClusterSync has the same namespace/name as the corresponding ClusterDeployment
	return fmt.Sprintf("%s/%s/%s", ns, ssOrsss, status.Name)
}

// shouldRecord looks at a SyncStatus and answers whether we should mark it down as a failure we want to alert on.
func shouldRecord(status *hiveinternal.SyncStatus) bool {
	// We only care about failures
	if status.Result != hiveinternal.FailureSyncSetResult {
		return false
	}
	// ...that are at least a certain age.
	// TODO: Make this age configurable via HiveHealthConfig
	var oldEnough time.Duration = time.Hour * 4
	var statusAge time.Duration = time.Now().Sub(status.LastTransitionTime.Time)
	if statusAge < oldEnough {
		return false
	}
	// TODO: Make specific ns/ssOrsss/ssname combinations silenceable via HiveHealthConfig
	return true
}

// recordFailures modifies `failures`, removing any successes and ensuring any failures are present.
// Must be under lock.
func recordFailures(ns, ssOrsss string, status *hiveinternal.SyncStatus) error {
	if shouldRecord(status) {
		k := key(ns, ssOrsss, status)
		failures.statuses[k] = status
		failures.byNamespace[ns] = append(failures.byNamespace[ns], k)
	}
	return nil
}

func countErrors(cs *hiveinternal.ClusterSync, logger logr.Logger) error {
	// Since we're using globals, synchronize for multiple controller threads
	failures.mutex.Lock()
	defer failures.mutex.Unlock()

	ns := cs.GetNamespace()

	// Start by wiping out all entries associated with this namespace. This is as opposed
	// to removing entries that have succeeded -- that would miss removing entries for
	// syncsets that have been deleted and are no longer in the ClusterSync.
	if byNamespace, ok := failures.byNamespace[ns]; ok {
		for _, statusKey := range byNamespace {
			delete(failures.statuses, statusKey)
		}
	}
	failures.byNamespace[ns] = make([]string, 0)

	ssOrsss := "sss"
	// Process selectorSyncSets
	for _, s := range cs.Status.SelectorSyncSets {
		scopy := s.DeepCopy()
		if err := recordFailures(ns, ssOrsss, scopy); err != nil {
			logger.Error(err, "")
		}
	}
	ssOrsss = "ss"
	// Process syncSets
	for _, s := range cs.Status.SyncSets {
		scopy := s.DeepCopy()
		if err := recordFailures(ns, ssOrsss, scopy); err != nil {
			logger.Error(err, "")
		}
	}

	// No need to leak namespace keys
	// FIXME: But we still could, when a namespace is deleted.
	if len(failures.byNamespace[ns]) == 0 {
		delete(failures.byNamespace, ns)
	}

	return nil
}
