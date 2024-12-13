// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

const (
	DataSource ResourceType = iota
	Resource
	Function
)

type ResourceType int

type InputGraphQLQuery struct {
	QueryName               string
	ObjectName              string
	Required                string
	Fields                  []Field
	GenqlientFields         []GenqlientField
	genqlientFieldsModify   []GenqlientField
	genqlientFieldsReadOnly []GenqlientField
	ResourceType            ResourceType
}

type Field struct {
	Name              string
	HumanReadableName string
	Type              string
}

type GenqlientField struct {
	Field
	Query                  string
	QueryNoPrefixReplaceId string
	InputObjectNames       string
	PlainObject            string
}

type DataSourceTemplateData struct {
	QueryName       string
	ObjectName      string
	Required        string
	StructName      string
	Fields          []Field
	GenqlientFields []GenqlientField
}

type ResourceTemplateData struct {
	QueryName               string
	ObjectName              string
	Required                string
	StructName              string
	Fields                  []Field
	GenqlientFields         []GenqlientField
	GenqlientFieldsModify   []GenqlientField
	GenqlientFieldsReadOnly []GenqlientField
}
type ProviderSourceTemplateData struct {
	DataSources []string
	Resources   []string
	Functions   []string
}

type TerraformComponents struct {
	dataSources []string
	resources   []string
}
