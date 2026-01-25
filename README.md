### 檔案架構:

```code
yourproject/
├── go.mod
├── cmd/ // 每個子資料夾通常對應一個可執行檔，Handler 與 Middleware 屬於「應用程式入口」的一部分，Handler 是路由對應的邏輯處理器。
│ ├── handler.go // 集中管理業務邏輯，工廠函式，也易於外來測試，或是更換其他DB等使用
│ ├── main.go // 主程序，處理初始化db、結構體驗證規則，透過go-playground
│ ├── customer.go // 顧客訂單查詢、顯示訂單、創建訂單等業務邏輯規劃
│ ├── routes.go  // 規劃後端endpoint路由，透過 handler 綁定相關邏輯
│ ├── utils.go // 載入環境變數，例如引用db路徑或是port等資訊調用，並透過 html/│   template 套件處理資料的引用
│ └── validators.go // 其他輔助檔案 : 定義驗證規則，不管前端傳甚麼給後端，都會多層過
│ └── routes.go // 應用程式層級的邏輯
│ └──main.go // 應用程式層級的邏輯
濾驗證，確保傳的內容符合規格定義的種類
│
├── internal/   // 私有庫（只能本專案用）
│ └── model/    // 在這定義結構體變數跟業務邏輯會用到的方法，另外能避免非法資料進入資料庫，確保系統只處理有效的業務邏輯 (validators.go )，可預防前端人員在未來傳入非當時後端開出的規格的髒資料。
│ └── order.go  // 創建訂單
│ └── user.go  // 創建登入用帳號資訊，並對其進行加密
│
└── pkg/ // 可重用的公共庫（optional）
```

### 筆記區 
- [此專案筆記位置](https://www.notion.so/go-1c6a54651e3e80808c81ce1843e7931e)
- [Gradient Generator](https://gradienty.codes/)

### 啟用專案:
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

### 登入整體流程圖: /login Router → Handler → Model → DB
> 避免使用者反覆輸入帳號、密碼而產生的機制
```
+-------------------+
|      Router       |
+-------------------+
| + POST("/login", HandleLoginPost) |  => Router 綁定到 Handler 的方法，這些Handler方法都放在cmd/xx.go裡定義的的方法
+-------------------+
+-------------------+
|      Handler      | => cmd/handlers.go => 是一個中介層，負責處理路由進來的請求。它本身持有不同的 Model（例如 UserModel、OrderModel），並透過這些 Model 來執行資料存取或商業邏輯。
+-------------------+
| - users: UserModel| => 呼叫對應的 Model (users欄位)
+-------------------+
| + HandleLoginPost(c: gin.Context) | => 1.對登入資訊加鹽 2.透過user.ID + user.Username存到 Session (為了避免使用者反覆輸入帳號、密碼而產生的機制)
+-------------------+
            |
            | users  => 每個 Model 對應到資料庫中的某個領域（例如使用者、訂單），並且封裝了該領域的操作方法（例如 AuthenticateUser 用來驗證使用者）。
            ▼
+-------------------+
|    UserModel      | => h.users.AuthenticateUser
+-------------------+
| - DB: *gorm.DB    | => Model 封裝了資料庫操作邏輯
+-------------------+
| + AuthenticateUser(username: string, password: string): User | => 驗證使用者的操作方法，先接收username跟password然後對資料庫操作。
+-------------------+
            |
            | returns => 回傳加密後的使用者資訊
            ▼
+-------------------+
|       User        |
+-------------------+
| - ID: int         |
| - Username: string|
| - Password: string|
+-------------------+

```
## 專案選用套件說明:

### gin-contrib
> 使用原因: Gin框架是支持中間件的，在不同的場景下，中間件有不同的含義，而在Gin框架中，中間件可以看作是請求攔截器，主要用在請求處理函數被調用前後執行一些業務邏輯，比如用戶權限驗證，數據編碼轉換，記錄操作日誌，接口運行時間統計等。加上實作上我們是透過 middleware來帶一層對session判斷是否已經登入，進行路由導轉，所以選用 gin-contrib 是相對合適的。

```
<!-- 關鍵字: middleware gin session recomand --> https://zhuanlan.zhihu.com/p/30184285330 -->
https://github.com/gin-contrib/sessions
<!-- 進度條是要實時更新的，客戶端不需要主動發送請求 -->
https://github.com/gin-contrib/sse
```


### Gin Middleware (gin.HandlerFunc 範例)
> 透過session是否存在，來決定Redirect的路徑，對於網頁來說，會有區分登入後才能觀看的，跟沒有登入時也能觀看的頁面，下方的架構主要就是在實踐這一塊。

- **全域 Middleware (LoggerMiddleware)**
  - 請求進來時先記錄 `Request Path`

- **路由群組 Middleware (AuthMiddleware)**
  - 如果是 `/admin` 路由，會先檢查使用者是否登入
  - **沒登入** → Redirect `/login` 並中斷
  - **登入成功** → 繼續執行下一個 Handler

- **Handler**
  - `ServeAdminDashboard`：顯示後台首頁
  - `HandleOrderPut`：更新訂單資訊

- **全域 Middleware (LoggerMiddleware After)**
  - Handler 執行完畢後，再記錄 `Response Status`

```
// cmd/middleware.go
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // [步驟 1] 進入 /admin 路由群組時，先執行 AuthMiddleware

        userID := GetSessionString(c, "userID")

        <!-- 沒登入 → Redirect /login 並中斷。 -->
        if userID == "" {
            // [步驟 2] 如果沒有 userID，導向 /login 並中斷後續流程
            c.Redirect(http.StatusSeeOther, "/login")
            c.Abort()
            return
        }

        _, err := h.users.GetUserByID(userID)
        if err != nil {
            // [步驟 3] 如果 userID 無效，清除 session，導向 /login 並中斷
            ClearSession(c)
            c.Redirect(http.StatusSeeOther, "/login")
            c.Abort()
            return
        }
        <!-- 登入成功 → 繼續。 -->
        // [步驟 4] 驗證成功，繼續執行下一個 Handler
        c.Next()
    }
}

// cmd/routes.go
func setupRoutes(router *gin.Engine, h *Handler) {
    <!-- 不用登入就可以方問 -->
    router.GET("/", h.ServeNewOrderForm)

    <!-- 需要登入才可以方問 -->
    admin := router.Group("/admin")
    admin.Use(h.AuthMiddleware()) // [步驟 A] 進入 /admin 路由前，會先跑 AuthMiddleware
    {
        // [步驟 B] 如果通過驗證(登入成功後)，才會執行以下 Handler
        admin.GET("", h.ServeAdminDashboard)       // 顯示後台首頁
        admin.POST("/order/:id/update", h.HandleOrderPut) // 更新訂單
    }
}


// cmd/main.go 
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // [步驟 0-前] 全域 Middleware：請求進來時先印出 Request path
        println("Request path:", c.Request.URL.Path)

        c.Next() // [步驟 0-中] 繼續執行後續 Middleware / Handler

        // [步驟 0-後] Handler 執行完畢後，再印出 Response status
        println("Response status:", c.Writer.Status())
    }
}


// cmd/main.go
func main() {
    r := gin.Default()

    // [第一層] 全域 Middleware：LoggerMiddleware
    r.Use(LoggerMiddleware()) 

    h := NewHandler(dbModel)
    // [第二層] 設定路由，/admin 路由會套用 AuthMiddleware
    setupRoutes(r,h) 

    // [最後] 啟動伺服器，開始監聽請求
    r.Run(":8080")
}
```

### 常見符號
```
& 表示取址，支持就地修改，如設置日誌或連接池
```

### 知識點

#### ShouldBindJSON
> 是現代前端/API 在送結構化 JSON 資料，開發者想明確指定只收 JSON

```
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func ApiLogin(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil { 
        c.JSON(422, gin.H{"error": "格式錯誤或欄位缺失", "details": err.Error()})
        return
    }
    // 登入邏輯...
}
```

#### ShouldBind
> 傳統網頁表單在送資料（form 格式），Gin 需要自動判斷格式
```
if err := c.ShouldBind(&req); err != nil {
    c.JSON(422, gin.H{"error": "格式錯誤或欄位缺失", "details": err.Error()})
    return
}
```

### 常用sql語法
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
sqlite3 data/orders.db "INSERT INTO users (username, password) VALUES('admin', '$2a$12$7Fy63im5z3jHEDn08hQbzevdLJIkDOgi52S79B58nplylten5QKtq');"

<!-- 查看 users 的 DB Table的內容 -->
sqlite3 -header -column data/orders.db "SELECT * FROM users;"

<!-- 覆寫特定欄位對應的Value (bash模式下才可以看到完整的密碼，由於太長的關係，powershell看不到完整的密碼) -->
sqlite3 data/orders.db 'INSERT OR REPLACE INTO users (username, password) VALUES("admin", "$2a$11$KuR6igHoxf/yUKrl4IW0GO0ID6uRE3bWxQ5XMBno0N/cWNf8KtVi6");'

```

### 錯誤訊息快速查找

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

#### tailwind相關
```
w-full => 填满其直接父元素的宽度
<!-- 指的是 min-h-screen 這裡的塊級元素 -->
<div class="min-h-screen flex items-center justify-center p-4">
    <div class="bg-white p-8 rounded-2xl shadow-xl w-full"></div>
</div>

max-w-md => max-width: 28rem; /* 448px */ 元素宽度不会超过 448px，但如果父容器比 448px 窄，元素会自动缩小以适应，通常需要配合 w-full 或其他宽度类使用，才能在小屏幕上伸缩。
```

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




