# Правильное проектирование на микро-уровне

## Пример 1 — [Kubernetes Scheduler](https://github.com/kubernetes/kubernetes/blob/master/pkg/scheduler/backend/queue/active_queue.go)

Исходный код:
~~~go
// addEventsIfPodInFlight adds events to inFlightEvents if the newPod is in inFlightPods.
// It returns true if pushed the event to the inFlightEvents.
func (aq *activeQueue) addEventsIfPodInFlight(oldPod, newPod *v1.Pod, events []fwk.ClusterEvent) bool {
	aq.lock.Lock()
	defer aq.lock.Unlock()

	return aq.unlockedQueue.addEventsIfPodInFlight(oldPod, newPod, events)
}
~~~

Исправленный код:
~~~go
type PodTransition struct {
    Old *v1.Pod
    New *v1.Pod
}

func (aq *activeQueue) addEventsIfPodInFlight(
    transition PodTransition,
    events []fwk.ClusterEvent,
)
~~~

## Пример 2 — [Prometheus](https://github.com/prometheus/prometheus/blob/main/tsdb/compact.go)

Исходный код:
~~~go
// Compactor provides compaction against an underlying storage
// of time series data.
type Compactor interface {
	// Write persists one or more Blocks into a directory.
	// No Block is written when resulting Block has 0 samples and returns an empty slice.
	// Prometheus always return one or no block. The interface allows returning more than one
	// block for downstream users to experiment with compactor.
	Write(dest string, b BlockReader, mint, maxt int64, base *BlockMeta) ([]ulid.ULID, error)
}
~~~

Исправленный код:
~~~go
type timeRange struct{
    min int64
    max int64
}

// Compactor provides compaction against an underlying storage
// of time series data.
type Compactor interface {
	// Write persists one or more Blocks into a directory.
	// No Block is written when resulting Block has 0 samples and returns an empty slice.
	// Prometheus always return one or no block. The interface allows returning more than one
	// block for downstream users to experiment with compactor.
	Write(dest string, b BlockReader, range timeRange, base *BlockMeta) ([]ulid.ULID, error)
}
~~~