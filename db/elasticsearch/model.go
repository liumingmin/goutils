package elasticsearch

type BaseModel struct {
	KeyName   string
	IndexName string
	TypeName  string
}

type QueryModel struct {
	BaseModel
	Query   interface{} //elastic.Query
	Sort    []string
	Cursor  int
	Size    int
	Results interface{}
	Total   *int64
}

type SourceModel struct {
	BaseModel
	Source  string
	Results interface{}
	Total   *int64
}

type AggregateModel struct {
	BaseModel
	Source  string
	AggKeys []string
}
