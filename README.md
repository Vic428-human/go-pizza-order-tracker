### 檔案架構

```code
yourproject/
├── go.mod
├── cmd/
│ ├── handler.go // 集中管理業務邏輯，工廠函式，也易於外來測試，或是更換其他DB等使用
│ ├── main.go // 主程序，處理初始化db、結構體驗證規則，透過go-playground
│ ├── customer.go // 顧客訂單查詢、顯示訂單、創建訂單等業務邏輯規劃
│ ├── routes.go  // 規劃後端endpoint路由，透過 handler 綁定相關邏輯
│ ├── utils.go // 載入環境變數，例如引用db路徑或是port等資訊調用，並透過 html/template 套件處理資料的引用
│ └── validators.go // 其他輔助檔案 : 定義驗證規則，不管前端傳甚麼給後端，都會多層過
│
濾驗證，確保傳的內容符合規格定義的種類
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

本專案常用指令

```
<!-- 確定module都有載到在 go.sum -->
go mod tidy
<!-- 運行 main.go裡的程序 -->
go run ./cmd
```

### 錯誤整理列表

#### Undefined validation function 'min' on field 'Phone'

```
<!-- 這錯誤跟 binding 的寫法錯誤有關 -->
<!-- 修正如下 -->
type OrderRequest struct {
	Phone        string   `json:"phone" binding:"required,min=10,max=20"`
}
```

#### (\*Handler).HandleNewOrderPost: Instructions: form.Instructions[i],

```
<!--
1. 表單沒有傳 Instructions ，你的 HTML form 可能沒有包含 Instructions 欄位，或是 name 不對，導致 form.Instructions 沒有值。
2.Gin 綁定問題 ， 如果你用 c.Bind(&form) 或 c.ShouldBind(&form)，但表單欄位名稱跟 struct tag 不一致，slice 就會是空的。
解法: 確認 tmpl 跟 結構體的 binding 都有正確綁定
-->
type OrderReuqest struct {
	Instructions []string `form:"instructions" binding:"max=200"`
}
<textarea maxlength="200" name="instructions"></textarea>
```

#### Error #01: template: customer.tmpl:42:35: executing "customer.tmpl" at <$index>: wrong type for value; expected int; got string

```
<!-- , $status 因為當時少寫到這個，所以 "add"引用的時候因為沒寫到 $status 這個字串 ，所以誤把 $index 當成字串了-->
{{range $index, $status := .Statuses}}
<div class="flex-1 flex justify-center mx-2">
    <div id="step{{add $index 1}}"
        class="size-14 bg-gray-300 rounded-full flex items-center justify-center text-white font-bold tansition-all duration-300 z-10 shadow-md">
        {{add $index 1}}
    </div>
</div>
{{end}}
```

### 其他實用技巧

#### 生成 svg 檔案

```
@"
<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
    <g>
        <circle cx="50" cy="50" r="45" fill="#FBC02D"/>
        <circle cx="50" cy="50" r="40" fill="#FFD54F"/>
        <path d="M50,50 L50,5 A45,45 0 0 1 95,50 Z" fill="#FBC02D"/>
        <path d="M50,50 L50,10 A40,40 0 0 1 90,50 Z" fill="#FFD54F"/>
        <circle cx="70" cy="30" r="4" fill="#D84315"/>
        <circle cx="30" cy="70" r="4" fill="#D84315"/>
        <circle cx="70" cy="70" r="4" fill="#D84315"/>
        <circle cx="30" cy="30" r="4" fill="#D84315"/>
        <circle cx="50" cy="50" r="4" fill="#D84315"/>
    </g>
</svg>
"@ | Out-File -FilePath pizza.svg -Encoding utf8

```
