# gotager
gotager是一个go代码开发辅助工具。
它可以让你更轻松地添加和更新结构体字段的标签。
[![License](http://img.shields.io/:license-Apache%202-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.txt)

## 安装
```
go get -u github.com/DennisMao/gotager
```

## 使用方法

### 基本使用
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
+ 支持tags删除
+ 支持tags内部参数插入
+ 支持更多转换风格
