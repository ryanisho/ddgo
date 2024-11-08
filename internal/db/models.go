package db

import (
	"encoding/json"
	"time"
)

type Metric struct {
    ID        int64             `json:"id"`
    Name      string            `json:"name"`
    Value     float64           `json:"value"`
    Tags      map[string]string `json:"tags"`
    Timestamp time.Time         `json:"timestamp"`
}

func (db *DB) SaveMetric(metric *Metric) error {
	tags, err := json.Marshal(metric.Tags)

	if err != nil {
		return err
	}

	query := `
		INSERT INTO metrics (name, value, tags, timestamp)
		VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(query, metric.Name, metric.Value, string(tags), metric.Timestamp)

	if err != nil {
		return err
	}

	metric.ID, _ = result.LastInsertId()
	return nil
}

func (db *DB) GetMetrics(name string, start, end time.Time) ([]Metric, error) {
	query := `
		SELECT id, name, values, tags, timestamp
		FROM metrics
		WHERE name = ? AND timestamp BETWEEN ? AND ?
		ORDER BY timestamp DESC
	`

	rows, err := db.Query(query, name, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []metrics
	for rows.Next() {
		var m Metric 
		var tagsJSON string
		err := rows.Scan(&m.ID, &m.Name, &m.Value, &tags.JSON, &m.Timestamp)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]bytes(tagsJSON), &m.Tags)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, m)
	}
	return metrics, rows.Err()
}

func (db *DB) CleanupOldMetrics(retention, time.Duration) error {
	cutoff := time.Now().Add(-retention)
	query := `DELETE FROM metrics WHERE timestamp < ?`
	_, err := db.Exec(query, cutoff)
	return err
}

type Alert struct {
    ID          string            `json:"id"`
    MetricName  string            `json:"metric_name"`
    Threshold   float64           `json:"threshold"`
    Operator    string            `json:"operator"`
    Duration    time.Duration     `json:"duration"`
    Tags        map[string]string `json:"tags"`
    IsTriggered bool             `json:"is_triggered"`
    LastCheck   time.Time         `json:"last_check"`
}

func (db *DB) SaveAlert(alert *Alert) error {
    tags, err := json.Marshal(alert.Tags)
    if err != nil {
        return err
    }

    query := `
        INSERT OR REPLACE INTO alerts 
        (id, metric_name, threshold, operator, duration, tags, is_triggered, last_check)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err = db.Exec(
        query,
        alert.ID,
        alert.MetricName,
        alert.Threshold,
        alert.Operator,
        int64(alert.Duration.Seconds()),
        string(tags),
        alert.IsTriggered,
        alert.LastCheck,
    )
    return err
}

func (db *DB) GetAlerts() ([]Alert, error) {
    query := `SELECT id, metric_name, threshold, operator, duration, tags, is_triggered, last_check FROM alerts`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alerts []Alert
    for rows.Next() {
        var a Alert
        var tagsJSON string
        var durationSecs int64
        err := rows.Scan(
            &a.ID,
            &a.MetricName,
            &a.Threshold,
            &a.Operator,
            &durationSecs,
            &tagsJSON,
            &a.IsTriggered,
            &a.LastCheck,
        )
        if err != nil {
            return nil, err
        }

        err = json.Unmarshal([]byte(tagsJSON), &a.Tags)
        if err != nil {
            return nil, err
        }

        a.Duration = time.Duration(durationSecs) * time.Second
        alerts = append(alerts, a)
    }

    return alerts, rows.Err()
}


type AlertHistory struct {
    ID         int64     `json:"id"`
    AlertID    string    `json:"alert_id"`
    MetricName string    `json:"metric_name"`
    Value      float64   `json:"value"`
    Threshold  float64   `json:"threshold"`
    Timestamp  time.Time `json:"timestamp"`
}

func (db *DB) SaveAlertHistory(history *AlertHistory) error {
    query := `
        INSERT INTO alert_history (alert_id, metric_name, value, threshold, timestamp)
        VALUES (?, ?, ?, ?, ?)
    `
    result, err := db.Exec(
        query,
        history.AlertID,
        history.MetricName,
        history.Value,
        history.Threshold,
        history.Timestamp,
    )
    if err != nil {
        return err
    }

    history.ID, _ = result.LastInsertId()
    return nil
}

func (db *DB) GetAlertHistory() ([]AlertHistory, error) {
    query := `
        SELECT id, alert_id, metric_name, value, threshold, timestamp
        FROM alert_history
        ORDER BY timestamp DESC
        LIMIT 100
    `
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var history []AlertHistory
    for rows.Next() {
        var h AlertHistory
        err := rows.Scan(
            &h.ID,
            &h.AlertID,
            &h.MetricName,
            &h.Value,
            &h.Threshold,
            &h.Timestamp,
        )
        if err != nil {
            return nil, err
        }
        history = append(history, h)
    }

    return history, rows.Err()
}