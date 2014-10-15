// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package elasticsearch

import (
	"encoding/json"
	"fmt"
)

// Query DSL - Queries

// Query represents a query in the elasticsearch DSL.
type Query interface {
	json.Marshaler
}

// Filter represents a filter in the elasticsearch DSL.
type Filter interface {
	json.Marshaler
}

// BoostField creates a string which represents a field name with a boost value.
func BoostField(field string, boost float64) string {
	return fmt.Sprintf("%s^%f", field, boost)
}

// MatchAllQuery provides a query that matches all
// documents in the index.
type MatchAllQuery struct {
}

func (m MatchAllQuery) MarshalJSON() ([]byte, error) {
	return marshalJSON("match_all", struct{}{})
}

// MatchQuery provides a query that matches against
// a complete field.
type MatchQuery struct {
	Field string
	Query string
	Type  string
}

func (m MatchQuery) MarshalJSON() ([]byte, error) {
	params := map[string]interface{}{"query": m.Query}
	if m.Type != "" {
		params["type"] = m.Type
	}

	return marshalJSON("match", map[string]interface{}{m.Field: params})
}

// MultiMatchQuery provides a query that matches on a number of fields.
type MultiMatchQuery struct {
	Query  string
	Fields []string
}

func (m MultiMatchQuery) MarshalJSON() ([]byte, error) {
	return marshalJSON("multi_match", map[string]interface{}{
		"query":  m.Query,
		"fields": m.Fields,
	})
}

// FilteredQuery provides a query that includes a filter.
type FilteredQuery struct {
	Query  Query
	Filter Filter
}

func (f FilteredQuery) MarshalJSON() ([]byte, error) {
	return marshalJSON("filtered", map[string]interface{}{
		"query":  f.Query,
		"filter": f.Filter,
	})
}

// FunctionScoreQuery provides a query that adjusts the scoring of a
// query by applying functions to it.
type FunctionScoreQuery struct {
	Query     Query
	Functions []Function
}

// FunctionScoreQuery provides a query that includes.
func (f FunctionScoreQuery) MarshalJSON() ([]byte, error) {
	return marshalJSON("function_score", map[string]interface{}{
		"query":     f.Query,
		"functions": f.Functions,
	})
}

// Function is a function definition for use with a FunctionScoreQuery.
type Function struct {
	Function string
	Field    string
	Scale    string
}

func (f Function) MarshalJSON() ([]byte, error) {
	return marshalJSON(f.Function, map[string]interface{}{
		f.Field: map[string]interface{}{
			"scale": f.Scale,
		},
	})
}

// AndFilter provides a filter that matches if all of the internal
// filters match.
type AndFilter []Filter

func (a AndFilter) MarshalJSON() ([]byte, error) {
	return marshalJSON("and", map[string]interface{}{
		"filters": []Filter(a),
	})
}

// OrFilter provides a filter that matches if any of the internal
// filters match.
type OrFilter []Filter

func (o OrFilter) MarshalJSON() ([]byte, error) {
	return marshalJSON("or", map[string]interface{}{
		"filters": []Filter(o),
	})
}

// NotFilter provides a filter that matches the opposite of the
// wrapped filter.
type NotFilter struct {
	Filter Filter
}

func (n NotFilter) MarshalJSON() ([]byte, error) {
	return marshalJSON("not", n.Filter)
}

// QueryFilter provides a filter that matches when a query matches
// on a result
type QueryFilter struct {
	Query Query
}

func (q QueryFilter) MarshalJSON() ([]byte, error) {
	return marshalJSON("query", q.Query)
}

// RegexpFilter provides a filter that matches a field against a
// regular expression.
type RegexpFilter struct {
	Field  string
	Regexp string
}

func (r RegexpFilter) MarshalJSON() ([]byte, error) {
	return marshalJSON("regexp", map[string]string{r.Field: r.Regexp})
}

// TermFilter provides a filter that requires a field to match.
type TermFilter struct {
	Field string
	Value string
}

func (t TermFilter) MarshalJSON() ([]byte, error) {
	return marshalJSON("term", map[string]string{t.Field: t.Value})
}

// QueryDSL provides a structure to put together a query using the
// elasticsearch DSL.
type QueryDSL struct {
	Size  int    `json:"size,omitempty"`
	Query Query  `json:"query,omitempty"`
	Sort  []Sort `json:"sort,omitempty"`
}

type Sort struct {
	Field string
	Order Order
}

type Order struct {
	Order string `json:"order"`
}

func (s Sort) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]Order{
		s.Field: {s.Order.Order},
	})
}

// marshalJSON provides a helper that creates json objects in a form
// often required by the elasticsearch query DSL. The objects created
// take the following form:
//	{
//		name: obj
//	}
func marshalJSON(name string, obj interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{name: obj})
}
