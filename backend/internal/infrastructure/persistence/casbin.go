package persistence

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"gorm.io/gorm"
)

const casbinModel = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (g(r.sub, p.sub, r.dom) || r.sub == p.sub) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
`

// InitCasbin 初始化 RBAC 执行器并从数据库同步策略
func InitCasbin(db *gorm.DB) *casbin.Enforcer {
	m, err := model.NewModelFromString(casbinModel)
	if err != nil {
		log.Fatalf("Casbin 模型加载失败: %v", err)
	}

	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		log.Fatalf("Casbin 初始化失败: %v", err)
	}

	if err := SyncCasbinPolicies(db, enforcer); err != nil {
		log.Printf("⚠️ Casbin 策略同步跳过（可能尚未迁移）: %v", err)
	}

	log.Println("✅ Casbin RBAC 初始化成功")
	return enforcer
}
