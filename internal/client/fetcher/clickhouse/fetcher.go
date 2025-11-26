package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	variable "github.com/w-h-a/caus/api/variable/v1alpha1"
	"github.com/w-h-a/caus/internal/client/fetcher"
)

type clickhouseFetcher struct {
	options fetcher.Options
	conn    *sql.DB
}

func (f *clickhouseFetcher) Fetch(ctx context.Context, v variable.VariableDefinition, start time.Time, end time.Time, step time.Duration) (map[time.Time]float64, error) {
	interval := fmt.Sprintf("%d", int(step.Minutes()))
	if interval == "0" {
		interval = "1" // lowest interval is 1m
	}

	results, err := f.aggregateSpans(
		ctx,
		v.TraceQuery.ServiceName,
		v.TraceQuery.Dimension,
		v.TraceQuery.AggregationOption,
		start,
		end,
		interval,
		v.TraceQuery.SpanName,
		v.TraceQuery.SpanKind,
		v.TraceQuery.AttributeQueries...,
	)
	if err != nil {
		return nil, fmt.Errorf("clickhouse fetcher failed to aggregate spans: %w", err)
	}

	data := map[time.Time]float64{}
	for _, r := range results {
		t, err := time.Parse("2006-01-02 15:04:05", r.Time)
		if err != nil {
			t, err = time.Parse(time.RFC3339, r.Time)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time string '%s': %w", r.Time, err)
			}
		}
		data[t.UTC().Truncate(0)] = r.Value
	}

	return data, nil
}

func (f *clickhouseFetcher) aggregateSpans(
	ctx context.Context,
	serviceName string,
	dimension string,
	aggregationOption string,
	start time.Time,
	end time.Time,
	interval string,
	spanName string,
	spanKind string,
	attributeQueries ...variable.AttributeQuery,
) ([]fetcher.Result, error) {
	aggregateQuery := ""

	switch dimension {
	case "duration":
		switch aggregationOption {
		case "avg":
			aggregateQuery = "avg(Duration) as value"
		case "p50":
			aggregateQuery = "quantile(0.50)(Duration) as value"
		case "p95":
			aggregateQuery = "quantile(0.95)(Duration) as value"
		case "p99":
			aggregateQuery = "quantile(0.99)(Duration) as value"
		}
	case "calls":
		aggregateQuery = "count(*) as value"
	}

	query := fmt.Sprintf(`SELECT toStartOfInterval(Timestamp, INTERVAL %s minute) as time, %s FROM default.otel_traces WHERE Timestamp>=? AND Timestamp<=?`, interval, aggregateQuery)
	args := []any{start, end}
	var err error

	query, args, err = f.buildSpanQuery(
		query,
		args,
		serviceName,
		spanName,
		spanKind,
		attributeQueries,
	)
	if err != nil {
		return nil, err
	}

	query += " GROUP BY time ORDER By time"

	results, err := f.readSpans(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (f *clickhouseFetcher) buildSpanQuery(
	query string,
	args []any,
	serviceName string,
	spanName string,
	spanKind string,
	attributeQueries []variable.AttributeQuery,
) (string, []interface{}, error) {
	if len(serviceName) != 0 {
		query += " AND ServiceName=?"
		args = append(args, serviceName)
	}

	if len(spanName) != 0 {
		query += " AND SpanName=?"
		args = append(args, spanName)
	}

	if len(spanKind) != 0 {
		query += " AND SpanKind=?"
		args = append(args, spanKind)
	}

	for _, attributeQ := range attributeQueries {
		if attributeQ.Key == "error" && attributeQ.Value == "true" {
			query += " AND (SpanAttributes['error']='true' OR StatusCode='Error')"
			continue
		}

		switch attributeQ.Operator {
		case "equals":
			query += " AND SpanAttributes[?]=?"
			args = append(args, attributeQ.Key, attributeQ.Value)
		case "contains":
			query += " AND SpanAttributes[?] ILIKE ?"
			args = append(args, attributeQ.Key, fmt.Sprintf("%%%s%%", attributeQ.Value))
		case "isnotnull":
			query += " AND mapContains(SpanAttributes, ?)"
			args = append(args, attributeQ.Key)
		}
	}

	return query, args, nil
}

func (f *clickhouseFetcher) readSpans(ctx context.Context, q string, as ...any) ([]fetcher.Result, error) {
	rows, err := f.conn.QueryContext(ctx, q, as...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []fetcher.Result

	for rows.Next() {
		var r fetcher.Result
		if err := rows.Scan(&r.Time, &r.Value); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func NewFetcher(opts ...fetcher.Option) fetcher.Fetcher {
	options := fetcher.NewOptions(opts...)

	// TODO: validate options

	cf := &clickhouseFetcher{
		options: options,
	}

	conn, err := sql.Open("clickhouse", options.Location)
	if err != nil {
		detail := fmt.Sprintf("clickhouse fetcher failed to connect: %v", err)
		panic(detail)
	}

	if err := conn.Ping(); err != nil {
		detail := fmt.Sprintf("clickhouse fetcher failed to ping: %v", err)
		panic(detail)
	}

	cf.conn = conn

	return cf
}
