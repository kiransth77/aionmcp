# Self-Learning Engine Usage Guide

## Overview

The AionMCP Self-Learning Engine is an autonomous system that monitors tool executions, analyzes patterns, and generates actionable insights for optimization. This guide covers how to use and configure the learning system.

## Configuration

### Environment Variables

```bash
# Enable/disable learning
export AIONMCP_LEARNING_ENABLED=true

# Data collection settings
export AIONMCP_LEARNING_SAMPLE_RATE=1.0        # 0.0-1.0, collect all executions
export AIONMCP_LEARNING_RETENTION_DAYS=30      # How long to keep data
export AIONMCP_LEARNING_ASYNC_PROCESSING=true  # Non-blocking collection
export AIONMCP_LEARNING_INCLUDE_SUCCESSFUL=true # Learn from successes

# Storage settings
export AIONMCP_STORAGE_PATH="./data/aionmcp.db"
```

### Configuration File (config.yaml)

```yaml
learning:
  enabled: true
  sample_rate: 1.0
  retention_days: 30
  async_processing: true
  include_successful: true
  include_input_output: true
  pii_filter_enabled: true
  max_input_size: 1024
  max_output_size: 4096

storage:
  type: boltdb
  path: "./data/aionmcp.db"
```

## API Endpoints

### Learning Statistics

Get comprehensive system-wide learning statistics:

```bash
GET /api/v1/learning/stats
```

**Response:**
```json
{
  "total_executions": 150,
  "success_rate": 0.94,
  "average_latency": 250000000,
  "error_breakdown": {
    "network": 5,
    "validation": 3,
    "configuration": 1
  },
  "top_tools": [
    {
      "name": "openapi.petstore.listPets",
      "execution_count": 75,
      "success_rate": 0.96,
      "average_latency": 180000000,
      "last_used": "2025-10-26T10:30:00Z"
    }
  ],
  "recent_patterns": [...],
  "active_insights": [...],
  "last_updated": "2025-10-26T10:30:00Z"
}
```

### Insights Management

#### Get All Insights

```bash
GET /api/v1/learning/insights
```

#### Filter by Type or Priority

```bash
GET /api/v1/learning/insights?type=performance
GET /api/v1/learning/insights?priority=high
```

**Response:**
```json
{
  "insights": [
    {
      "id": "insight_abc123",
      "type": "performance",
      "priority": "high",
      "title": "Performance Issues in API Tool",
      "description": "Tool shows consistently slow performance",
      "suggestion": "Consider implementing caching or optimizing API calls",
      "evidence": [
        "Average latency: 2.5s",
        "Execution count: 100",
        "Success rate: 98%"
      ],
      "created_at": "2025-10-26T10:00:00Z",
      "metadata": {
        "tool_name": "api_tool",
        "pattern_id": "pattern_xyz789"
      }
    }
  ]
}
```

### Pattern Analysis

#### Get All Patterns

```bash
GET /api/v1/learning/patterns
```

#### Filter by Type

```bash
GET /api/v1/learning/patterns?type=error
GET /api/v1/learning/patterns?type=performance
GET /api/v1/learning/patterns?type=usage
```

#### Tool-Specific Patterns

```bash
GET /api/v1/learning/patterns?tool=openapi.petstore.listPets
```

### Tool-Specific Insights

Get insights for a specific tool:

```bash
GET /api/v1/learning/tools/openapi.petstore.listPets/insights
```

### Manual Analysis

Trigger manual pattern analysis and insight generation:

```bash
POST /api/v1/learning/analyze
```

**Response:**
```json
{
  "patterns_found": 3,
  "insights_generated": 2
}
```

### Configuration Management

Get current learning configuration:

```bash
GET /api/v1/learning/config
```

## Understanding Insights

### Insight Types

- **`optimization`**: General optimization opportunities
- **`configuration`**: Configuration-related improvements
- **`reliability`**: Reliability and error-related insights
- **`performance`**: Performance optimization suggestions
- **`usage`**: Usage pattern insights

### Priority Levels

- **`critical`**: Immediate attention required (e.g., >50% error rate)
- **`high`**: Should be addressed soon (e.g., performance bottlenecks)
- **`medium`**: Good to optimize when possible
- **`low`**: Minor improvements or informational

## Pattern Types

### Error Patterns

Detect recurring errors by:
- Tool name and error type combination
- Frequency threshold (≥3 occurrences)
- Error classification (network, validation, configuration, etc.)

### Performance Patterns

Identify performance issues:
- Tools with latency >2x system average
- Degrading performance trends
- Timeout-related issues

### Usage Patterns

Understand usage distribution:
- Tools consuming >50% of executions
- Underutilized tools
- Peak usage patterns

## Best Practices

### 1. Regular Monitoring

Check learning statistics regularly:

```bash
# Daily stats check
curl http://localhost:8080/api/v1/learning/stats | jq '.success_rate, .total_executions'

# Weekly insight review
curl http://localhost:8080/api/v1/learning/insights?priority=high
```

### 2. Acting on Insights

**High Priority Insights:**
- Review immediately
- Implement suggested optimizations
- Monitor impact after changes

**Medium Priority Insights:**
- Schedule for next maintenance window
- Plan optimization work
- Consider during refactoring

### 3. Configuration Tuning

**High Volume Systems:**
- Reduce sample rate: `sample_rate: 0.1` (10% sampling)
- Shorter retention: `retention_days: 7`
- Ensure async processing: `async_processing: true`

**Development/Testing:**
- Full sampling: `sample_rate: 1.0`
- Include all data: `include_input_output: true`
- Longer retention for analysis: `retention_days: 90`

### 4. Privacy Considerations

- Enable PII filtering: `pii_filter_enabled: true`
- Limit data size: `max_input_size: 1024`
- Regular cleanup with appropriate retention periods

## Troubleshooting

### Learning Not Working

1. **Check Configuration:**
   ```bash
   curl http://localhost:8080/api/v1/learning/config
   ```

2. **Verify Tool Executions:**
   ```bash
   curl http://localhost:8080/api/v1/learning/stats
   ```

3. **Check Logs:**
   ```bash
   # Look for learning-related errors
   grep -i learning /var/log/aionmcp.log
   ```

### No Patterns Detected

- Ensure sufficient executions (>10 for pattern detection)
- Check pattern thresholds in analyzer configuration
- Verify data retention hasn't cleaned up recent data

### Performance Impact

- Enable async processing: `async_processing: true`
- Reduce sample rate for high-volume scenarios
- Monitor storage growth and adjust retention

## Integration Examples

### Monitoring Dashboard

```bash
#!/bin/bash
# Simple monitoring script

echo "=== AionMCP Learning Stats ==="
curl -s http://localhost:8080/api/v1/learning/stats | jq '
{
  total_executions,
  success_rate,
  average_latency_ms: (.average_latency / 1000000),
  error_count: (.error_breakdown | add)
}'

echo -e "\n=== High Priority Insights ==="
curl -s http://localhost:8080/api/v1/learning/insights?priority=high | jq '.insights[] | {
  title,
  description,
  suggestion
}'
```

### Automated Analysis

```bash
#!/bin/bash
# Trigger analysis and alert on critical insights

curl -s -X POST http://localhost:8080/api/v1/learning/analyze > /dev/null

CRITICAL_COUNT=$(curl -s http://localhost:8080/api/v1/learning/insights?priority=critical | jq '.insights | length')

if [ "$CRITICAL_COUNT" -gt 0 ]; then
    echo "ALERT: $CRITICAL_COUNT critical insights found!"
    curl -s http://localhost:8080/api/v1/learning/insights?priority=critical | jq '.insights[] | .title'
fi
```

## Advanced Usage

### Custom Analysis Queries

The learning system provides rich data that can be queried programmatically:

```python
import requests
import json

# Get tool performance comparison
stats = requests.get('http://localhost:8080/api/v1/learning/stats').json()

for tool in stats['top_tools']:
    latency_ms = tool['average_latency'] / 1_000_000
    print(f"{tool['name']}: {latency_ms:.1f}ms avg, {tool['success_rate']:.1%} success")
```

### Integration with Monitoring Systems

Export metrics to Prometheus, Grafana, or other monitoring systems:

```yaml
# Example Grafana dashboard query
# Success rate over time
rate(aionmcp_executions_total[5m])

# Tool performance
histogram_quantile(0.95, aionmcp_execution_duration_seconds)
```

## Security Considerations

1. **API Access Control:** Secure learning endpoints with authentication
2. **Data Privacy:** Enable PII filtering for sensitive environments  
3. **Storage Security:** Protect BoltDB files with appropriate file permissions
4. **Network Security:** Use HTTPS for production deployments
5. **Data Retention:** Comply with data retention policies and regulations