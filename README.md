# gotager
gotager is a tiny devtool for Golang coding。
You can add/update tags for fields in structs more easily. 

## Install
```
go get -u github.com/DennisMao/gotager
```

## Usage


### Simple
test.go
```
package test

type Hello struct {
	Name     string
	Id       int64 `bson:"hi" json:"Id"`
	Email    string
	N1N2     string
	_SET     string
	S_2_1_NN string
}
```

run gotager tool for test.go
```
$./gotaget  test.go 

```

test.go 
```
package test

type Hello struct {
	Name		string	`json:"name"`
	Id		    int64	`bson:"hi" json:"Id"`
	Email		string	`json:"email"`
	N1N2		string	`json:"n_1_n2"`
	_SET		string	`json:"_set"`
	S_2_1_NN	string	`json:"s_2_1_nn"`
}
```

## Todo
+ Support remove tags
+ Support inner tag param insertion
+ Support more conversion style
