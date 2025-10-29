package selflearn

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/boltdb/bolt"
	"go.uber.org/zap"
)

// BoltStorage implements Storage interface using BoltDB
type BoltStorage struct {
	db     *bolt.DB
	logger *zap.Logger
}

// Bucket names for different data types
const (
	ExecutionsBucket = "executions"
	PatternsBucket   = "patterns"
	InsightsBucket   = "insights"
	StatsBucket      = "stats"
)

// NewBoltStorage creates a new BoltDB storage instance
func NewBoltStorage(dbPath string, logger *zap.Logger) (*BoltStorage, error) {
	// Ensure directory exists
	if err := ensureDir(filepath.Dir(dbPath)); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	storage := &BoltStorage{
		db:     db,
		logger: logger,
	}

	// Initialize buckets
	if err := storage.initBuckets(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize buckets: %w", err)
	}

	return storage, nil
}

// initBuckets creates the required buckets if they don't exist
func (s *BoltStorage) initBuckets() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		buckets := []string{ExecutionsBucket, PatternsBucket, InsightsBucket, StatsBucket}
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}
		return nil
	})
}

// StoreExecution stores an execution record
func (s *BoltStorage) StoreExecution(ctx context.Context, record ExecutionRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal execution record: %w", err)
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ExecutionsBucket))
		if bucket == nil {
			return fmt.Errorf("executions bucket not found")
		}

		// Use timestamp + ID as key for time-based ordering
		key := fmt.Sprintf("%d_%s", record.Timestamp.Unix(), record.ID)
		return bucket.Put([]byte(key), data)
	})
}

// GetExecution retrieves an execution record by ID
func (s *BoltStorage) GetExecution(ctx context.Context, id string) (ExecutionRecord, error) {
	var record ExecutionRecord

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ExecutionsBucket))
		if bucket == nil {
			return fmt.Errorf("executions bucket not found")
		}

		// Search for the record by ID (since key includes timestamp)
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var exec ExecutionRecord
			if err := json.Unmarshal(v, &exec); err != nil {
				continue // Skip invalid records
			}
			if exec.ID == id {
				record = exec
				return nil
			}
		}
		return fmt.Errorf("execution record not found: %s", id)
	})

	return record, err
}

// GetExecutionsByTool retrieves execution records for a specific tool
func (s *BoltStorage) GetExecutionsByTool(ctx context.Context, toolName string, limit int) ([]ExecutionRecord, error) {
	var records []ExecutionRecord

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ExecutionsBucket))
		if bucket == nil {
			return fmt.Errorf("executions bucket not found")
		}

		cursor := bucket.Cursor()
		count := 0

		// Iterate in reverse order (newest first)
		for k, v := cursor.Last(); k != nil && count < limit; k, v = cursor.Prev() {
			var record ExecutionRecord
			if err := json.Unmarshal(v, &record); err != nil {
				s.logger.Warn("Failed to unmarshal execution record", zap.Error(err))
				continue
			}

			if record.ToolName == toolName {
				records = append(records, record)
				count++
			}
		}

		return nil
	})

	return records, err
}

// GetExecutionsByTimeRange retrieves execution records within a time range
func (s *BoltStorage) GetExecutionsByTimeRange(ctx context.Context, start, end time.Time, limit int) ([]ExecutionRecord, error) {
	var records []ExecutionRecord

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ExecutionsBucket))
		if bucket == nil {
			return fmt.Errorf("executions bucket not found")
		}

		cursor := bucket.Cursor()
		count := 0

		// Create key range for time-based search
		startKey := []byte(fmt.Sprintf("%d", start.Unix()))
		endKey := []byte(fmt.Sprintf("%d", end.Unix()))

		for k, v := cursor.Seek(startKey); k != nil && count < limit; k, v = cursor.Next() {
			// Check if we've exceeded the end time
			if string(k) > string(endKey) {
				break
			}

			var record ExecutionRecord
			if err := json.Unmarshal(v, &record); err != nil {
				s.logger.Warn("Failed to unmarshal execution record", zap.Error(err))
				continue
			}

			// Double-check time range (inclusive)
			if !record.Timestamp.Before(start) && !record.Timestamp.After(end) {
				records = append(records, record)
				count++
			}
		}

		return nil
	})

	return records, err
}

// GetExecutionStats calculates and returns learning statistics
func (s *BoltStorage) GetExecutionStats(ctx context.Context) (LearningStats, error) {
	stats := LearningStats{
		ErrorBreakdown: make(map[ErrorType]int),
		TopTools:       []ToolStat{},
		LastUpdated:    time.Now().UTC(),
	}

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ExecutionsBucket))
		if bucket == nil {
			return fmt.Errorf("executions bucket not found")
		}

		toolStats := make(map[string]*ToolStat)
		var totalDuration time.Duration
		var successCount int64

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var record ExecutionRecord
			if err := json.Unmarshal(v, &record); err != nil {
				continue
			}

			stats.TotalExecutions++
			totalDuration += record.Duration

			if record.Success {
				successCount++
			} else {
				stats.ErrorBreakdown[record.ErrorType]++
			}

			// Update tool statistics
			if toolStat, exists := toolStats[record.ToolName]; exists {
				toolStat.ExecutionCount++
				if record.Success {
					toolStat.SuccessRate = (toolStat.SuccessRate*(float64(toolStat.ExecutionCount)-1) + 1) / float64(toolStat.ExecutionCount)
				} else {
					toolStat.SuccessRate = toolStat.SuccessRate * (float64(toolStat.ExecutionCount) - 1) / float64(toolStat.ExecutionCount)
				}
				toolStat.AverageLatency = time.Duration((int64(toolStat.AverageLatency)*(toolStat.ExecutionCount-1) + int64(record.Duration)) / toolStat.ExecutionCount)
				if record.Timestamp.After(toolStat.LastUsed) {
					toolStat.LastUsed = record.Timestamp
				}
			} else {
				successRate := 0.0
				if record.Success {
					successRate = 1.0
				}
				toolStats[record.ToolName] = &ToolStat{
					Name:           record.ToolName,
					ExecutionCount: 1,
					SuccessRate:    successRate,
					AverageLatency: record.Duration,
					LastUsed:       record.Timestamp,
				}
			}
		}

		// Calculate overall statistics
		if stats.TotalExecutions > 0 {
			stats.SuccessRate = float64(successCount) / float64(stats.TotalExecutions)
			stats.AverageLatency = totalDuration / time.Duration(stats.TotalExecutions)
		}

		// Convert tool stats to slice and sort by execution count
		for _, stat := range toolStats {
			stats.TopTools = append(stats.TopTools, *stat)
		}
		sort.Slice(stats.TopTools, func(i, j int) bool {
			return stats.TopTools[i].ExecutionCount > stats.TopTools[j].ExecutionCount
		})

		// Limit to top 10 tools
		if len(stats.TopTools) > 10 {
			stats.TopTools = stats.TopTools[:10]
		}

		return nil
	})

	return stats, err
}

// StorePattern stores a pattern
func (s *BoltStorage) StorePattern(ctx context.Context, pattern Pattern) error {
	data, err := json.Marshal(pattern)
	if err != nil {
		return fmt.Errorf("failed to marshal pattern: %w", err)
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PatternsBucket))
		return bucket.Put([]byte(pattern.ID), data)
	})
}

// GetPattern retrieves a pattern by ID
func (s *BoltStorage) GetPattern(ctx context.Context, id string) (Pattern, error) {
	var pattern Pattern

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PatternsBucket))
		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("pattern not found: %s", id)
		}
		return json.Unmarshal(data, &pattern)
	})

	return pattern, err
}

// GetPatterns retrieves patterns by type
func (s *BoltStorage) GetPatterns(ctx context.Context, patternType PatternType, limit int) ([]Pattern, error) {
	var patterns []Pattern

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PatternsBucket))
		cursor := bucket.Cursor()
		count := 0

		for k, v := cursor.First(); k != nil && count < limit; k, v = cursor.Next() {
			var pattern Pattern
			if err := json.Unmarshal(v, &pattern); err != nil {
				continue
			}

			if patternType == "" || pattern.Type == patternType {
				patterns = append(patterns, pattern)
				count++
			}
		}

		return nil
	})

	return patterns, err
}

// UpdatePattern updates an existing pattern
func (s *BoltStorage) UpdatePattern(ctx context.Context, pattern Pattern) error {
	return s.StorePattern(ctx, pattern)
}

// DeletePattern deletes a pattern
func (s *BoltStorage) DeletePattern(ctx context.Context, id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PatternsBucket))
		return bucket.Delete([]byte(id))
	})
}

// StoreInsight stores an insight
func (s *BoltStorage) StoreInsight(ctx context.Context, insight Insight) error {
	data, err := json.Marshal(insight)
	if err != nil {
		return fmt.Errorf("failed to marshal insight: %w", err)
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(InsightsBucket))
		return bucket.Put([]byte(insight.ID), data)
	})
}

// GetInsight retrieves an insight by ID
func (s *BoltStorage) GetInsight(ctx context.Context, id string) (Insight, error) {
	var insight Insight

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(InsightsBucket))
		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("insight not found: %s", id)
		}
		return json.Unmarshal(data, &insight)
	})

	return insight, err
}

// GetInsights retrieves insights by type
func (s *BoltStorage) GetInsights(ctx context.Context, insightType InsightType, limit int) ([]Insight, error) {
	var insights []Insight

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(InsightsBucket))
		cursor := bucket.Cursor()
		count := 0

		for k, v := cursor.First(); k != nil && count < limit; k, v = cursor.Next() {
			var insight Insight
			if err := json.Unmarshal(v, &insight); err != nil {
				continue
			}

			if insightType == "" || insight.Type == insightType {
				insights = append(insights, insight)
				count++
			}
		}

		return nil
	})

	return insights, err
}

// GetInsightsByPriority retrieves insights by priority
func (s *BoltStorage) GetInsightsByPriority(ctx context.Context, priority Priority, limit int) ([]Insight, error) {
	var insights []Insight

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(InsightsBucket))
		cursor := bucket.Cursor()
		count := 0

		for k, v := cursor.First(); k != nil && count < limit; k, v = cursor.Next() {
			var insight Insight
			if err := json.Unmarshal(v, &insight); err != nil {
				continue
			}

			if priority == "" || insight.Priority == priority {
				insights = append(insights, insight)
				count++
			}
		}

		return nil
	})

	return insights, err
}

// UpdateInsight updates an existing insight
func (s *BoltStorage) UpdateInsight(ctx context.Context, insight Insight) error {
	return s.StoreInsight(ctx, insight)
}

// DeleteInsight deletes an insight
func (s *BoltStorage) DeleteInsight(ctx context.Context, id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(InsightsBucket))
		return bucket.Delete([]byte(id))
	})
}

// Cleanup removes old records based on retention period
func (s *BoltStorage) Cleanup(ctx context.Context, retentionPeriod time.Duration) error {
	cutoff := time.Now().Add(-retentionPeriod)
	
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ExecutionsBucket))
		cursor := bucket.Cursor()
		
		var keysToDelete [][]byte
		
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var record ExecutionRecord
			if err := json.Unmarshal(v, &record); err != nil {
				// Delete invalid records
				keysToDelete = append(keysToDelete, k)
				continue
			}
			
			if record.Timestamp.Before(cutoff) {
				keysToDelete = append(keysToDelete, k)
			}
		}
		
		// Delete old records
		for _, key := range keysToDelete {
			if err := bucket.Delete(key); err != nil {
				s.logger.Warn("Failed to delete old record", zap.Error(err))
			}
		}
		
		s.logger.Info("Cleanup completed", zap.Int("deleted_records", len(keysToDelete)))
		return nil
	})
}

// Close closes the BoltDB connection
func (s *BoltStorage) Close() error {
	return s.db.Close()
}

// ensureDir creates directory if it doesn't exist
func ensureDir(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, 0755)
}