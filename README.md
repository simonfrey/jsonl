# JSON Lines (JSONL) golang library

This library provides you with a Reader and a Writer for the JSONL format. It automatically flushes underlying 
`http.ResponseWriter` if applicable. 

No external dependencies are required for this to run.

> JSON Lines is a convenient format for storing structured data that may be processed one record at a time. It works well with unix-style text processing tools and shell pipelines. It's a great format for log files. It's also a flexible format for passing messages between cooperating processes.

Source: [jsonlines.org](https://jsonlines.org/)

## Examples

*Dropping errors to keep code example shorter. You should definitely check them!*

### Write JSONL (JSON Lines)

The library automatically does the json marshaling for you on `Write(in interface{})`

```go
package main

import(
    "github.com/simonfrey/jsonl"
    "bytes"
)

func main(){
    buff := bytes.Buffer{}
    
    w := jsonl.NewWriter(&buff)
    w.Write("Hello")
    w.Write("World")
    w.Write(42)
   
	fmt.Println(buff.String())
	// Output: 
	//  "Hello"\n"World"\n42\n
}
```

### Read JSONL (JSON Lines)

More interesting than writing into JSONL is reading it and working with it. This library provides you with two functions
to do so.

#### Read Single Line and unmarshal into struct

If you exactly know what you expect you can directly read a single line from the data into a struct. The library does
the json unmarshaling for you.

```go
package main

import(
    "github.com/simonfrey/jsonl"
    "strings"
)

func main(){
    input := "\"Hello\"\n\"World\"\n42\n"
    
    r := jsonl.NewReader(strings.NewReader(input))
   
	outString := ""
	r.ReadSingleLine(&outString)
	
	fmt.Println(outString)
	// Output: 
	//  Hello
}
```

#### Read multiple entries with type detection

This is my favorite use case of JSONL. Reading different JSON types from one stream. This allows to build a form of streaming
data interface from the server (which sends JSONL with above writer function)

As you see for stricter type safety I use a dedicated `Type` field in the structs, so we can do a string/byte comparison
to figure out what type we are using. If you types do have unique fields you could also check for those. 

In the following example you learn how you can read JSONL typesafe in golang.

```go
package main

import(
    "github.com/simonfrey/jsonl"
    "strings"
)


type T1 struct {
    Type  string `json:"type"`
    Count int    `json:"am"`
}
type T2 struct {
    Type    string `json:"type"`
    Comment string `json:"am"`
}


func main(){
	input := "{\"type\":\"T1\",\"am\":2}\n{\"type\":\"T2\",\"am\":\"I am T2\"}\n{\"type\":\"T1\",\"am\":9999}\n"
    
    r := jsonl.NewReader(strings.NewReader(input))
    r.ReadLines(func(data []byte) error {
        switch {
        case bytes.Contains(data, []byte(`"T1"`)):
            // T1 struct type
            t := T1{}
            err := json.Unmarshal(data, &t)
            if err != nil {
            return fmt.Errorf("could not unmarshal into T1: %w", err)
            }
			fmt.Printf("%T: %d\n", t, t.Count)
        case bytes.Contains(data, []byte(`"T2"`)):
            // T2 struct type
            t := T2{}
            err := json.Unmarshal(data, &t)
            if err != nil {
            return fmt.Errorf("could not unmarshal into T2: %w", err)
            }
            fmt.Printf("%T: %s\n", t, t.Comment)
            output += fmt.Sprintf("%T:%s|", t, t.Comment)
        }
        return nil
    })
	
	// Output: 
	//  main.T1: 2
	//  main.T2: I am T2
	//  main.T1: 9999
}
```


## License

MIT
