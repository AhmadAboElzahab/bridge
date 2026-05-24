package models

type Operator struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type FilterItem struct {
	ID             string      `json:"id"`
	Type           string      `json:"type"`
	FieldID        uint        `json:"fieldId"`
	ColumnType     string      `json:"columnType"`
	Operator       Operator    `json:"operator"`
	SecondOperator *Operator   `json:"secondOperator,omitempty"`
	Value          interface{} `json:"value"`
}

type FilterGroup struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Conjunction string        `json:"conjunction"`
	Children    []interface{} `json:"children"`
}
