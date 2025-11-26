package fetcher

type Result struct {
	Time  string  `db:"time"`
	Value float64 `db:"value"`
}
