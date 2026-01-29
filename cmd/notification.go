package main

import "sync"

/*
RWMutex (讀寫鎖)：區分「讀」和「寫」兩種操作：
讀鎖 (RLock)：允許多個 goroutine 同時讀取資源，只要沒有 goroutine 在寫。
寫鎖 (Lock)：只允許一個 goroutine 寫入資源，並且會阻塞所有其他的讀和寫。
*/
type NotificationManager struct {
	clients map[string]map[chan string]bool
	mu      sync.RWMutex // 讀寫鎖 (Read-Write Mutex)
}

/* clients: make(map[string]map[chan string]bool)
當有新訂單事件發生時，系統只需遍歷對應主題的 channel 清單，即可只推送給相關訂單的客戶端，避免浪費資源全域廣播。
這種巢狀 map[string]map[chan string]bool 結構實現多播通知（Pub/Sub Pattern）：
clients map[string]map[chan string]bool = {
	"order-123": { // 頻道名稱 (TOPIC)
		// 該群組內所有客戶端的 channel，當消息有發布的時候，只有下列這些client有訂閱過該頻道(TOPIC)的才會收到訊息
		0xc0000a4000: true,
		0xc0000a4060: true,
	},
	"order-456": {
		0xc0000a40c0: true,
	},
	"admin:new_orders": {
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

// 1. 訂閱頻道
func (n *NotificationManager) Subscribe(topic string, client chan string) {
	n.mu.Lock()         // 🔒 上鎖
	defer n.mu.Unlock() // 🔓 自動解鎖

	if n.clients[topic] == nil { // ⚠️ Race Condition，使用 Lock 跟 Unlock 就不會有這問題
		n.clients[topic] = make(map[chan string]bool)
	}
	n.clients[topic][client] = true
}

// 2. 取消訂閱頻道
func (n *NotificationManager) Unsubscribe(topic string, client chan string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if clients, ok := n.clients[topic]; ok {
		delete(clients, client)
	}

	if len(n.clients[topic]) == 0 {
		delete(n.clients, topic)
	}

	close(client)
}

// 2. 對特定頻道發送通知
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
