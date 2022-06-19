package elasticsearch

import "encoding/json"

const (
	ES6 = "es6"
	ES7 = "es7"
)

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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type MappingProperty struct {
	Type     string                      `json:"type"`               //
	Index    bool                        `json:"index,omitempty"`    //
	Analyzer string                      `json:"analyzer,omitempty"` //
	Fields   map[string]*MappingProperty `json:"fields,omitempty"`   //
	ExtProps map[string]interface{}      `json:"-"`
}

func (t *MappingProperty) MarshalJSON() ([]byte, error) {
	propJson, err := json.Marshal(*t)
	if err != nil {
		return nil, err
	}

	if len(t.ExtProps) == 0 {
		return propJson, nil
	}

	var props map[string]interface{}
	err = json.Unmarshal(propJson, &props)
	if err != nil {
		return nil, err
	}

	for k, v := range t.ExtProps {
		if _, ok := props[k]; !ok {
			props[k] = v
		}
	}
	return json.Marshal(props)
}

type MappingSettings struct {
	SettingsIndex `json:"index"`
}

type SettingsIndex struct {
	IgnoreMalformed  bool `json:"mapping.ignore_malformed,omitempty"` // true
	NumberOfReplicas int  `json:"number_of_replicas"`                 // 1
	NumberOfShards   int  `json:"number_of_shards"`                   // 3
}
