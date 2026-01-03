```code
yourproject/
├── go.mod
├── cmd/
│ ├── main.go // 這裡是 package main 一定有 func main()，只負責「啟動程式」
│ └── validators.go // 其他輔助檔案，也寫 package main
│
├── internal/ // 私有庫（只能本專案用）
│ └── model/ // 這裡寫 package model
│ └── model/ order.go
│
│
│
│
└── pkg/ // 可重用的公共庫（optional）
```
