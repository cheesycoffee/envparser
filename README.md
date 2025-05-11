# envparser

A lightweight Go package for parsing environment variables into struct fields using tags, with support for nested/embedded structs and custom decoding.

## ✨ Features

* Simple struct-tag-based configuration
* Support for various primitive types
* Nested and embedded structs
* JSON, XML, Form Data, Base64 decoding via struct tags
* Custom error aggregation
* Works with unexported structs in the same package

---

## 🍞 Installation

```bash
go get github.com/cheesycoffee/envparser
```

---

## 🧠 Supported Data Types

| Go Type                                             | Supported           |
| --------------------------------------------------- | ------------------- |
| `string`                                            | ✅                   |
| `int`, `int32`, `int64`                             | ✅                   |
| `uint`, `uint32`, `uint64`                          | ✅                   |
| `float32`, `float64`                                | ✅                   |
| `bool`                                              | ✅                   |
| `time.Duration`                                     | ✅                   |
| `time.Time` (RFC3339 format)                        | ✅                   |
| `[]string`                                          | ✅ (comma-separated) |
| `[]int`, `[]uint` , `[]uint32`, `[]uint64`.         | ✅ (comma-separated) |
| `[]float32`, `[]float64`.                           | ✅ (comma-separated) |
| Structs (anonymous/embedded)                        | ✅                   |
| Structs with `json`/`xml`/`form`/`base64` tags via `encoding:"xml"`/`encoding:"json"`/`encoding:"form"`/`encoding:"base64"` | ✅                   |

---

## 🔧 Usage

### 1. Basic Usage

```go
package main

import (
	"log"
	"time"
	"github.com/joho/godotenv"
	"github.com/cheesycoffee/envparser"
)

type Nested struct {
    NestedValueString string `env:"NESTED_VALUE_STRING"`
    NestedValueInt    int    `env:"NESTED_VALUE_INT"`
}

type Embeded struct {
    EmbededValueString string `env:"EMBEDED_VALUE_STRING"`
    EmbededValueInt    int    `env:"EMBEDED_VALUE_INT"`
}

type JSONData struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type XMLData struct {
    Name string `xml:"name"`
    Age  int    `xml:"age"`
}

type Config struct {
	AppName       string        `env:"APP_NAME"`
	Port          int           `env:"PORT"`
	Debug         bool          `env:"DEBUG"`
    PhiVal        float32       `env:"PHI_VAL"`
	Timeout       time.Duration `env:"TIMEOUT"`
	LaunchAt      time.Time     `env:"LAUNCH_AT"` // RFC3339 format
	AllowedHosts  []string      `env:"ALLOWED_HOSTS"`
	IDs           []uint64      `env:"UINT_IDS"`
    NestedValue // nested
    EmbededValue  EmbededValue // embeded
    JSONData      JSONData      `env:"JSON_VALUE" encoding:"json"`
    XMLDATA       XMLDATA       `env:"XML_VALUE" encoding:"xml"`
    FormValue     url.Values    `env:"FORM_VALUE" encoding:"form"`
    FileData      []byte        `env:"FILE_VALUE" encoding:"base64"`
}

func main() {
	_ = godotenv.Load()

	var cfg Config
	if err := envparser.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", cfg)
}
```

### .env Example

```
APP_NAME=EnvApp
PORT=8080
DEBUG=true
TIMEOUT=30s
LAUNCH_AT="2023-10-01T15:04:05Z"
ALLOWED_HOSTS="example.com,api.example.com"
PHI_VAL="3.14"
UINT_IDS="1,2,4,5"
NESTED_VALUE_STRING="nested value"
NESTED_VALUE_INT=5
EMBEDED_VALUE_STRING="embeded value string"
EMBEDED_VALUE_INT=5
JSON_VALUE="{\"name\":\"Alice\",\"age\":30}"
XML_VALUE="<XMLData><name>Alice</name><age>30</age></XMLData>"
FORM_VALUE="name=alice&age=30"
FILE_VALUE="SGVsbG8gR28gd29ybGQh"
```

---

## ⚠️ Error Handling

If multiple fields fail to parse, `Parse` aggregates and returns them all:

```go
err := envparser.Parse(&cfg)
if err != nil {
	log.Fatalf("Failed to load config:\n%v", err)
}
```

---

## 👀 Notes

* Ensure the target is passed as a **pointer to a struct**: `Parse(&cfg)`
* Environment variable keys must be explicitly defined with `env:"KEY"`
* If a field has no `env` tag or is marked `env:"-"`, it will be ignored
* Embedded/anonymous structs are parsed recursively
