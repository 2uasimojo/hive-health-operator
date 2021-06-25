package clustersync

import (
	"fmt"

	"github.com/go-logr/logr"
)

// processAlerts uses the current state of the globals in failure_counter.go to determine whether
// we need to processAlerts on any failing syncsets, and do so.
func processAlerts(logger logr.Logger) error {
	failures.mutex.RLock()
	defer failures.mutex.RUnlock()

	for ns, keys := range failures.byNamespace {
		message := fmt.Sprintf("ClusterDeployment in Namespace %s has %d failing [Selector]SyncSet(s):\n", ns, len(keys))
		for _, key := range keys {
			status := failures.statuses[key]
			// TODO: Filter alerts based on the age of the status.
			// ...which I think might not be possible at the moment, because hive
			// does not store a time stamp for the initial failure.
			// Perhaps we should save such information in a CR somewhere.
			// Problem is, neither we nor hive can (re)create that information idempotently.
			message += fmt.Sprintf("\t%s: %s\n", status.Name, status.FailureMessage)
		}
		alert(message, logger)
	}
	return nil
}

func alert(message string, logger logr.Logger) {
	// TODO: Make this actually alert
	logger.Info("!!ALERT!!", "Message", message)
}
