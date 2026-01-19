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
│ └── model/    // 在這定義結構體變數跟業務邏輯會用到的方法，另外能避免非法資料進入資料庫，確保系統只處理有效的業務邏輯 (validators.go )，可預防前端人員在未來傳入非當時後端開出的規格的髒資料。
│ └── order.go  // 創建訂單
│ └── user.go  // 創建登入用帳號資訊，並對其進行加密
└── pkg/ // 可重用的公共庫（optional）
```

### 啟用專案
```
輸入 sqlite3 確認是否出現下方訊息，有的話代表有安裝sqlite cli
<!-- https://sqlite.org/download.html -->
<!-- 本專案使用的是 sqlite-tools-osx-x64-351020zip 沒有  ARM的那個版本  -->
<!-- 並且要記得設定環境變數，解壓縮後，改黨名sqlite，然後複製檔案路徑，貼到PATH重啟vscode，輸入 sqlite3 有出現下方就代表安裝成功 -->
SQLite version 3.51.2 2026-01-09 17:27:48
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open FILENAME" to reopen on a persistent database.
```

### 常見符號
```
& 表示取址，支持就地修改，如設置日誌或連接池
```

### 專案常用指令

```
<!-- 確定module都有載到在 go.sum -->
go mod tidy

<!-- 運行 main.go 裡的程序，如果不先運行是無法對DB Table進行資料修改的 -->
go run ./cmd

<!-- 查看 orders 的 DB -->
sqlite3 -header -column data/orders.db "SELECT * FROM orders;"

<!-- 對 orders 的 DB Table 寫入，修改特定欄位的資料 -->
sqlite3 data/orders.db "UPDATE orders SET status = 'Preparing' WHERE id='aULpvdIDR';" 

<!-- 對 users 的 DB Table 插入新的資料 -->
sqlite3 data/orders.db "INSERT INTO users (username, password) VALUES('admin', '\$2a\$12\$ZyZgQMjHvs41bMEX0i82jeqeWfz08Q9Vusx./QQJTNkfh2QWGLRa6');"

<!-- 查看 users 的 DB Table的內容 -->
sqlite3 -header -column data/orders.db "SELECT * FROM users;"

```

### 錯誤整理列表

####  Error: in prepare, no such table XXX
```
<!-- 表示 models結構體尚未定義在 Automigrate，導致資料表不能更新，把 User{} 放在 AutoMigrate 裡即可。-->
	
err = db.AutoMigrate(&Order{}, &OrderItem{}, &User{})

type User struct { 
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}
err = db.AutoMigrate(&Order{}, &OrderItem{}, &User{})
```

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
