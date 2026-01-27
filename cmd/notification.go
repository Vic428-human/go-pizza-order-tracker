package main

type Notification struct {
	clients map[strig]map[chan string]bool
}


// 1. 按主題分組推送
// 當有新訂單事件發生時，系統只需遍歷對應主題的 channel 清單，即可只推送給相關訂單的客戶端，避免浪費資源全域廣播。
// 這種巢狀 map[string]map[chan string]bool 結構實現多播通知（Pub/Sub Pattern）：
clients map[string]map[string]bool = {
	// 代表通知主題或事件類型 => 外層 key（如 "order-123"、"admin:new_orders"）
	"order-123": { // 內層 map，好處：訂單123更新時，只通知兩個相關客戶端，不會打擾其他用戶。
		// 儲存該主題下所有感興趣的客戶端 channel，bool 作為簡單存在標記（避免重複註冊）
		0xc0000a4000: true,
		0xc0000a4060: true,
	},
	
	"order-456": { // 只有追蹤訂單456的客戶端
		0xc0000a40c0: true,
	},
	"admin:new_orders": { // 管理員全域訂單通知
		0xc0000a4180: true,
	},
}

// 2.支援動態訂閱/取消
// 客戶端訂閱 => n.clients["order-123"][0xc0000a4000] = true
// 客戶端斷線時清理 => delete(n.clients["order-123"], 0xc0000a4000)
