package elasticsearch

type QueryModel struct {
	IndexName string
	TypeName  string
	Query     interface{} //elastic.Query
	Sort      []string
	Cursor    int
	Size      int
	Results   interface{}
	Total     *int64
}

type SourceModel struct {
	IndexName string
	TypeName  string
	Source    string
	Results   interface{}
	Total     *int64
}

type AggregateModel struct {
	IndexName string
	TypeName  string
	Source    string
	AggKeys   []string
}
