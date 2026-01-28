package main

import "sync"

type NotificationManager struct {
	clients map[string]map[chan string]bool
	mu      sync.RWMutex
}

/*
// 1. 按主題分組推送
// 當有新訂單事件發生時，系統只需遍歷對應主題的 channel 清單，即可只推送給相關訂單的客戶端，避免浪費資源全域廣播。
// 這種巢狀 map[string]map[chan string]bool 結構實現多播通知（Pub/Sub Pattern）：

	clients map[string]map[chan string]bool = {
		// 代表通知主題或事件類型 => 外層 key（如 "order-123"、"admin:new_orders"）
		// 可能是用戶ID、session ID 或 room ID（群組識別）
		"order-123": { // 內層 map，好處：訂單123更新時，只通知兩個相關客戶端，不會打擾其他用戶。
			// 儲存該主題下所有感興趣的客戶端 channel，bool 作為簡單存在標記（避免重複註冊）
			0xc0000a4000: true,
			// 該群組內所有客戶端的 channel
			0xc0000a4060: true,
		},
		"order-456": { // 只有追蹤訂單456的客戶端
			0xc0000a40c0: true,
		},
		"admin:new_orders": { // 管理員全域訂單通知
			0xc0000a4180: true,
		},
	}
*/
func NewNotification() *NotificationManager {
	return &NotificationManager{clients: make(map[string]map[chan string]bool)}
}

// 2.支援動態訂閱/取消
// 客戶端訂閱 => n.clients["order-123"][0xc0000a4000] = true
// 客戶端斷線時清理 => delete(n.clients["order-123"], 0xc0000a4000)

// 1. 客戶端連線時註冊
func (n *NotificationManager) Subscribe(topic string, client chan string) {
	if n.clients[topic] == nil {
		n.clients[topic] = make(map[chan string]bool)
	}
	n.clients[topic][client] = true
}

// 2. 發送通知
func (n *NotificationManager) Publish(room_id string, message string) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if clients, ok := n.clients[room_id]; ok {
		for client := range clients {
			select {
			case client <- message: // 非阻塞發送
			default: // 客戶端緩衝滿，忽略
			}
		}
	}
}
