package dao

import (
	"github.com/jinzhu/gorm"
	"math"
	"micro-service-tmpl/internal/AI/domain/vo"
)

const (
	NotLearned  = 0 // 未学习
	Learning    = 1 // 学习中
	Learned     = 2 // 学习成功
	LearnFailed = 3 // 学习失败

)

//Ai误报信息
type AiDistort struct {
	gorm.Model
	Uin     uint64 `gorm:"index;not null;default:'0';comment:'Uin'" json:"uin,string"`                                     // Uin
	Appid   uint64 `gorm:"index;not null;default:'0';comment:'用户ID'" json:"appid,string"`                                  // 用户ID
	Domain  string `gorm:"type:varchar(64);not null;default:'';comment:'域名'" json:"domain"`                                // 域名
	Payload string `gorm:"type:varchar(64);not null;default:'';comment:'载荷'" json:"Payload"`                               // 载荷
	From    string `gorm:"type:varchar(64);not null;default:'';comment:'来源'" json:"from"`                                  // 来源
	Remark  string `gorm:"type:varchar(64);not null;default:'';comment:'备注'" json:"remark"`                                // 备注
	Status  uint8  `gorm:"type:tinyint(4);not null;default:'0';comment:'状态 0：【未学习】 1【学习中】 2【学习成功】 3【学习失败】'" json:"status"` // 状态 0：【未学习】 1【学习中】 2【学习成功】 3【学习失败】
}

func (AiDistort) TableName() string {
	return "ai_distort"
}

// 添加记录
func (r *AiDistort) Add(model *AiDistort) (err error) {
	err = MasterDB.Create(model).Error
	return
}

// 更新保存记录
func (r *AiDistort) Save(model *AiDistort) (err error) {
	err = MasterDB.Save(model).Error
	return
}

// 软删除：结构体需要继承Base model 有delete_at字段
func (r *AiDistort) Delete(query interface{}, args ...interface{}) (err error) {
	//return r.db.Unscoped().Where(query, args...).Delete(&AiDistort{}).Error //硬删除
	return MasterDB.Where(query, args...).Delete(&AiDistort{}).Error
}

// 根据条件获取单挑记录
func (r *AiDistort) First(query interface{}, args ...interface{}) (model AiDistort, err error) {
	err = SlaveDB.Where(query, args...).First(&model).Error
	return
}

// 获取列表 数据量大时Count数据需另外请求接口
func (r *AiDistort) Find(query interface{}, page *vo.Pagination, args ...interface{}) (models []AiDistort, err error) {
	if page == nil {
		err = SlaveDB.Find(&models).Error
		//err = SlaveDB.Find(&models).Error //从
	} else {
		err = SlaveDB.Model(AiDistort{}).Where(query, args...).
			Count(&page.Total).Offset((page.Page - 1) * page.PageSize).
			Limit(page.PageSize).Find(&models).Error
		// 总条数
		page.TotalPage = int64(math.Ceil(float64(page.Total / page.PageSize)))
	}

	return
}

// 获取总记录条数
func (r *AiDistort) Count(where interface{}, args ...interface{}) (count int64, err error) {
	err = SlaveDB.Model(&AiDistort{}).Where(where, args...).Count(&count).Error
	return
}
