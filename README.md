### 檔案架構

```code
yourproject/
├── go.mod
├── cmd/
│ ├── handler.go // 集中管理業務邏輯，工廠函式，也易於外來測試，或是更換其他DB等使用
│ └── validators.go // 其他輔助檔案 : 定義驗證規則，不管前端傳甚麼給後端，都會多層過濾驗證，確保傳的內容符合規格定義的種類
│
├── internal/   // 私有庫（只能本專案用）
│ └── model/    // 這裡寫 package model
│ └── order.go  // 做驗證的時候，可以在這邊取得一些 var model 變數，因為驗證的規則是基於各個模組裡的結構體
│
│
│
│
└── pkg/ // 可重用的公共庫（optional）
```
