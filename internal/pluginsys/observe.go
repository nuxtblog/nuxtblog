package pluginsys

import (
	"encoding/json"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// ─── PluginStats ──────────────────────────────────────────────────────────────

// PluginStats tracks per-plugin execution metrics.
// Numeric counters are atomic; string fields are guarded by lastErrMu.
type PluginStats struct {
	Invocations     atomic.Int64
	Errors          atomic.Int64
	TotalDurationNs atomic.Int64 // nanoseconds; divide by Invocations for average

	lastErrMu   sync.Mutex
	LastError   string
	LastErrorAt time.Time
}

func (s *PluginStats) record(d time.Duration, err error) {
	s.Invocations.Add(1)
	s.TotalDurationNs.Add(int64(d))
	if err != nil {
		s.Errors.Add(1)
		s.lastErrMu.Lock()
		s.LastError = err.Error()
		s.LastErrorAt = time.Now()
		s.lastErrMu.Unlock()
	}
}

// AvgDurationMs returns the mean execution time in milliseconds.
func (s *PluginStats) AvgDurationMs() float64 {
	n := s.Invocations.Load()
	if n == 0 {
		return 0
	}
	return float64(s.TotalDurationNs.Load()) / float64(n) / 1e6
}

// StatsSnapshot is a JSON-serialisable view of PluginStats for API responses.
type StatsSnapshot struct {
	PluginID      string    `json:"plugin_id"`
	Invocations   int64     `json:"invocations"`
	Errors        int64     `json:"errors"`
	AvgDurationMs float64   `json:"avg_duration_ms"`
	LastError     string    `json:"last_error,omitempty"`
	LastErrorAt   time.Time `json:"last_error_at,omitempty"`
}

// ─── ErrorRingBuffer ─────────────────────────────────────────────────────────

// ErrorEntry records one plugin execution error for post-hoc debugging.
type ErrorEntry struct {
	At        time.Time     `json:"at"`
	EventName string        `json:"event"`
	Phase     string        `json:"phase"`              // "filter" | "handler" | "pipeline" | "route"
	Message   string        `json:"message"`
	Stack     string        `json:"stack,omitempty"`    // JS stack trace when available
	Duration  time.Duration `json:"duration_ms"`        // execution time in milliseconds
	// InputDiff is a JSON representation of changes made to ctx.Data during
	// the filter chain (diff of Input→Data). Empty for dispatch (on) errors.
	InputDiff string `json:"input_diff,omitempty"`
}

// MarshalJSON customises serialisation so Duration is rendered as milliseconds.
func (e ErrorEntry) MarshalJSON() ([]byte, error) {
	type alias struct {
		At        time.Time `json:"at"`
		EventName string    `json:"event"`
		Phase     string    `json:"phase,omitempty"`
		Message   string    `json:"message"`
		Stack     string    `json:"stack,omitempty"`
		DurationMs float64  `json:"duration_ms,omitempty"`
		InputDiff string    `json:"input_diff,omitempty"`
	}
	return json.Marshal(alias{
		At:         e.At,
		EventName:  e.EventName,
		Phase:      e.Phase,
		Message:    e.Message,
		Stack:      e.Stack,
		DurationMs: float64(e.Duration.Nanoseconds()) / 1e6,
		InputDiff:  e.InputDiff,
	})
}

const ringBufCap = 100

// ErrorRingBuffer is a fixed-capacity (100), thread-safe circular buffer.
// When full, the oldest entry is silently overwritten.
type ErrorRingBuffer struct {
	mu   sync.Mutex
	buf  [ringBufCap]ErrorEntry
	head int // next write position
	size int // valid entries, 0..ringBufCap
}

// Add appends an entry, overwriting the oldest when at capacity.
func (r *ErrorRingBuffer) Add(e ErrorEntry) {
	r.mu.Lock()
	r.buf[r.head] = e
	r.head = (r.head + 1) % ringBufCap
	if r.size < ringBufCap {
		r.size++
	}
	r.mu.Unlock()
}

// GetAll returns a copy of all entries in chronological order (oldest first).
func (r *ErrorRingBuffer) GetAll() []ErrorEntry {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.size == 0 {
		return nil
	}
	out := make([]ErrorEntry, r.size)
	start := (r.head - r.size + ringBufCap) % ringBufCap
	for i := 0; i < r.size; i++ {
		out[i] = r.buf[(start+i)%ringBufCap]
	}
	return out
}

// recordExec records a single execution (route/filter/event) into the plugin's
// observability counters.
func (lp *loadedPlugin) recordExec(pluginID, eventName string, dur time.Duration, err error) {
	if lp.stats != nil {
		lp.stats.record(dur, err)
	}
	if lp.window != nil {
		lp.window.record(err != nil)
	}
	if err != nil && lp.errors != nil {
		lp.errors.Add(ErrorEntry{
			At:        time.Now(),
			EventName: eventName,
			Phase:     "route",
			Message:   err.Error(),
			Duration:  dur,
		})
	}
}

// getPluginObs returns the observability data for a plugin by ID.
func (m *Manager) getPluginObs(id string) *loadedPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.plugins[id]
}

// ─── Public query functions ───────────────────────────────────────────────────

// GetStats returns execution metrics for a plugin.
func GetStats(id string) *StatsSnapshot {
	m := GetManager()
	if m == nil {
		return nil
	}
	lp := m.getPluginObs(id)
	if lp == nil || lp.stats == nil {
		return nil
	}
	s := lp.stats
	s.lastErrMu.Lock()
	snap := &StatsSnapshot{
		PluginID:      id,
		Invocations:   s.Invocations.Load(),
		Errors:        s.Errors.Load(),
		AvgDurationMs: s.AvgDurationMs(),
		LastError:     s.LastError,
		LastErrorAt:   s.LastErrorAt,
	}
	s.lastErrMu.Unlock()
	return snap
}

// GetHistory returns the 60-minute sliding window history for a plugin.
func GetHistory(id string) []WindowBucket {
	m := GetManager()
	if m == nil {
		return nil
	}
	lp := m.getPluginObs(id)
	if lp == nil || lp.window == nil {
		return nil
	}
	return lp.window.GetHistory()
}

// GetErrors returns the recent error ring buffer contents for a plugin.
// Returns nil if the plugin is not loaded; returns an empty (non-nil) slice
// if loaded but no errors have been recorded.
func GetErrors(id string) []ErrorEntry {
	m := GetManager()
	if m == nil {
		return nil
	}
	lp := m.getPluginObs(id)
	if lp == nil || lp.errors == nil {
		return nil
	}
	entries := lp.errors.GetAll()
	if entries == nil {
		return []ErrorEntry{}
	}
	return entries
}

// ─── SlidingWindow ────────────────────────────────────────────────────────────

const windowBuckets = 60 // 60 × 1 min = 1 h of history

// WindowBucket is one minute's counters in the sliding window.
type WindowBucket struct {
	Minute      time.Time `json:"minute"`
	Invocations int64     `json:"invocations"`
	Errors      int64     `json:"errors"`
}

type windowSlot struct {
	minute      int64 // unix timestamp / 60 (minute number)
	invocations int64
	errors      int64
}

// SlidingWindow maintains 60 one-minute buckets.
// Reads and writes are guarded by a single mutex.
type SlidingWindow struct {
	mu      sync.Mutex
	buckets [windowBuckets]windowSlot
}

func (w *SlidingWindow) record(isError bool) {
	min := time.Now().Unix() / 60
	idx := int(min % windowBuckets)
	w.mu.Lock()
	if w.buckets[idx].minute != min {
		w.buckets[idx] = windowSlot{minute: min}
	}
	w.buckets[idx].invocations++
	if isError {
		w.buckets[idx].errors++
	}
	w.mu.Unlock()
}

// GetHistory returns the last 60 minutes of data, oldest bucket first.
// Buckets with no activity are returned with zero counters so the caller
// always gets exactly 60 data points.
func (w *SlidingWindow) GetHistory() []WindowBucket {
	now := time.Now().Unix() / 60
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]WindowBucket, windowBuckets)
	for i := 0; i < windowBuckets; i++ {
		min := now - int64(windowBuckets-1) + int64(i)
		idx := int(((min % windowBuckets) + windowBuckets) % windowBuckets)
		out[i].Minute = time.Unix(min*60, 0)
		if w.buckets[idx].minute == min {
			out[i].Invocations = w.buckets[idx].invocations
			out[i].Errors = w.buckets[idx].errors
		}
	}
	return out
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

// diffMaps produces a concise JSON diff between two maps.
// Keys prefixed with "+" are new, "~" are changed (with before/after), "-" are removed.
// Returns "" when maps are identical.
func diffMaps(before, after map[string]any) string {
	diff := map[string]any{}
	for k, v := range after {
		if bv, ok := before[k]; !ok {
			diff["+"+k] = v
		} else if !reflect.DeepEqual(bv, v) {
			diff["~"+k] = map[string]any{"before": bv, "after": v}
		}
	}
	for k, v := range before {
		if _, ok := after[k]; !ok {
			diff["-"+k] = v
		}
	}
	if len(diff) == 0 {
		return ""
	}
	b, _ := json.Marshal(diff)
	return string(b)
}
