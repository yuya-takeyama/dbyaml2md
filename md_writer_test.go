package main

import (
	"bytes"
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
`)
	table := new(model.Table)
	err := yaml.Unmarshal(tableYaml, table)
	if err != nil {
		t.Fatalf("Failed to unmarshal stub YAML: %s", err)
	}

	buf := new(bytes.Buffer)
	err = mdWriter.writeMarkdownFromTable(buf, table)
	if err != nil {
		t.Fatalf("Failed to write generated markdown into buffer: %s", err)
	}

	expected := []byte(`---
table: users
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
`)

	if bytes.Compare(expected, buf.Bytes()) != 0 {
		t.Fatalf("generated markdown is not as expected")
	}
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
	}

	buf := new(bytes.Buffer)
	mdWriter := &MdWriter{frontMatter}
	err := mdWriter.writeMarkdownFromTable(buf, table)
	if err != nil {
		t.Fatalf("Failed to write generated markdown into buffer: %s", err)
	}

	expectedPrefix := []byte(`---
array:
- 1
- true
- false
- null
string: foo
table: users
---
# users

Users table
`)

	if !bytes.HasPrefix(buf.Bytes(), expectedPrefix) {
		t.Fatalf("generated markdown's front matter is not as expected")
	}
}
