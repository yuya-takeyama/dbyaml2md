package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/yuya-takeyama/db2yaml/model"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestBasicTable(t *testing.T) {
	mdWriter := &MdWriter{make(FrontMatter)}

	tableYaml := []byte(`name: users
columns:
- name: id
  type: int
  auto_increment: true
  comment: User ID
- name: name
  type: varchar
  length: 128
  comment: User name
- name: birth
  type: datetime
  nullable: true
  comment: Birthday
indexes:
- name: PRIMARY
  unique: true
  columns:
  - name: id
- name: username
  unique: true
  columns:
  - name: name
  - name: id
comment: Users table
ddl: CREATE TABLE users
`)
	table := new(model.Table)
	err := yaml.Unmarshal(tableYaml, table)
	if err != nil {
		t.Fatalf("Failed to unmarshal stub YAML: %s", err)
	}

	buf := new(bytes.Buffer)
	err = mdWriter.writeMarkdown(buf, table)
	if err != nil {
		t.Fatalf("Failed to write generated markdown into buffer: %s", err)
	}

	expected := `---
table:
  name: users
  columns:
  - name: id
    type: int
    auto_increment: true
    comment: User ID
  - name: name
    type: varchar
    length: 128
    comment: User name
  - name: birth
    type: datetime
    nullable: true
    comment: Birthday
  indexes:
  - name: PRIMARY
    unique: true
    columns:
    - name: id
  - name: username
    unique: true
    columns:
    - name: name
    - name: id
  comment: Users table
  ddl: CREATE TABLE users
---
# users

Users table

## Columns

Name|Description|Type|Length|Default|Nullable|AUTO_INCREMENT|
----|-----------|----|-----:|-------|-------:|-------------:|
id|User ID|int||||✓|
name|User name|varchar|128||||
birth|Birthday|datetime|||✓||

## Indexes

<table>
  <thead>
    <tr>
      <th>Name</th>
      <th>Unique</th>
      <th>Columns</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>PRIMARY</td>
      <td style="text-align: right">✓</td>
      <td>
        <ul>
          <li>id</li>
        </ul>
      </td>
    </tr>
    <tr>
      <td>username</td>
      <td style="text-align: right">✓</td>
      <td>
        <ul>
          <li>name</li>
          <li>id</li>
        </ul>
      </td>
    </tr>
  </tbody>
</table>
`

	assert.Equal(t, buf.String(), expected)
}

func TestFrontMatter(t *testing.T) {
	frontMatter := make(FrontMatter)

	frontMatter["array"] = []interface{}{1, true, false, nil}
	frontMatter["string"] = "foo"

	table := &model.Table{
		Name:    "users",
		Columns: make([]*model.Column, 0),
		Indexes: make([]*model.Index, 0),
		Comment: "Users table",
		DDL:     "CREATE TABLE users",
	}

	buf := new(bytes.Buffer)
	mdWriter := &MdWriter{frontMatter}
	err := mdWriter.writeMarkdown(buf, table)
	if err != nil {
		t.Fatalf("Failed to write generated markdown into buffer: %s", err)
	}

	expectedPrefix := `---
array:
- 1
- true
- false
- null
string: foo
table:
  name: users
  columns: []
  indexes: []
  comment: Users table
  ddl: CREATE TABLE users
---
# users

Users table
`

	assert.Contains(t, buf.String(), expectedPrefix)
}
