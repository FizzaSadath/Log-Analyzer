package filter

import (
	"fmt"
	"log_analyzer/models"
	"sync"
	"time"
)

func FilterLogs(
	store models.LogStore,
	levels, components, hosts, reqIDs []string,
	startTime time.Time, endTime time.Time,
) []models.LogEntry {
	start := time.Now()

	var result []models.LogEntry
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Loop through segments concurrently
	for _, segment := range store.Segments {
		seg := segment // capture loop variable
		wg.Add(1)
		go func(seg models.Segment) {
			defer wg.Done()

			totalFilters := 0
			if !startTime.IsZero() && seg.EndTime.Before(startTime) {
				return
			}
			if !endTime.IsZero() && seg.StartTime.After(endTime) {
				return
			}

			matchedIndex := make(map[int]bool)

			if len(levels) > 0 {
				totalFilters++
				for _, level := range levels {
					for _, idx := range seg.Index.ByLevel[level] {
						matchedIndex[idx] = true
					}
				}
			}

			if len(components) > 0 {
				totalFilters++
				componentFilter := make(map[int]bool)
				for _, component := range components {
					for _, idx := range seg.Index.ByComponent[component] {
						if len(matchedIndex) == 0 || matchedIndex[idx] {
							componentFilter[idx] = true
						}
					}
				}
				matchedIndex = componentFilter
			}

			if len(hosts) > 0 {
				totalFilters++
				hostFilter := make(map[int]bool)
				for _, host := range hosts {
					for _, idx := range seg.Index.ByHost[host] {
						if len(matchedIndex) == 0 || matchedIndex[idx] {
							hostFilter[idx] = true
						}
					}
				}
				matchedIndex = hostFilter
			}

			if len(reqIDs) > 0 {
				totalFilters++
				requestFilter := make(map[int]bool)
				for _, reqID := range reqIDs {
					for _, idx := range seg.Index.ByReqID[reqID] {
						if len(matchedIndex) == 0 || matchedIndex[idx] {
							requestFilter[idx] = true
						}
					}
				}
				matchedIndex = requestFilter
			}

			var localResults []models.LogEntry

			// Case: no filters and no time constraints
			if totalFilters == 0 && startTime.IsZero() && endTime.IsZero() {
				localResults = append(localResults, seg.LogEntries...)
			} else if totalFilters == 0 {
				for _, entry := range seg.LogEntries {
					if !startTime.IsZero() && entry.Time.Before(startTime) {
						continue
					}
					if !endTime.IsZero() && entry.Time.After(endTime) {
						continue
					}
					localResults = append(localResults, entry)
				}
			}

			for idx := range matchedIndex {
				entry := seg.LogEntries[idx]
				if !startTime.IsZero() && entry.Time.Before(startTime) {
					continue
				}
				if !endTime.IsZero() && entry.Time.After(endTime) {
					continue
				}
				localResults = append(localResults, entry)
			}

			// Lock before appending to shared result slice
			mu.Lock()
			result = append(result, localResults...)
			mu.Unlock()
		}(seg)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Filtering took:", elapsed)
	return result
}
