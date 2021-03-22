package vo

//Ai误报信息
type AiDistort struct {
	Uin       uint64 `json:"uin,string"`   // Uin
	Appid     uint64 `json:"appid,string"` // 用户ID
	Domain    string `json:"domain"`       // 域名
	Payload   string `json:"Payload"`      // 载荷
	From      string `json:"from"`         // 来源
	Remark    string `json:"remark"`       // 备注
	Status    uint8  `json:"status"`       // 状态 0：【未学习】 1【学习中】 2【学习成功】 3【学习失败】
	CreatedAt int64  `json:"createdAt"`    // 创建时间戳
}

// Ai漏报信息
type AiFailure struct {
	Uin       uint64 `json:"userId,string"` // Uin
	Appid     uint64 `json:"appid,string"`  // 用户ID
	Domain    string `json:"domain"`        // 域名
	Payload   string `json:"Payload"`       // 载荷
	Sign      uint8  `json:"sign"`          // 标记 1：【其他】 2【XSS攻击】 3【SQL注入】
	From      string `json:"from"`          // 来源
	Remark    string `json:"remark"`        // 备注
	Status    uint8  `json:"status"`        // 状态 0：【未学习】 1【学习中】 2【学习成功】 3【学习失败】
	CreatedAt int64  `json:"createdAt"`     // 创建时间戳
}

// 页码结构体
type Pagination struct {
	Page      int64 `json:"page" example:"0"`      // 当前页
	PageSize  int64 `json:"pageSize" example:"20"` // 每页条数
	Total     int64 `json:"total"`                 // 总条数
	TotalPage int64 `json:"totalPage"`             // 总页数
}
