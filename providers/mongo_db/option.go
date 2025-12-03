package mongo_db

import (
	"context"
	"fmt"
	"time"

	"github.com/tqhuy-dev/xgen/utilities"
)

type Option struct {
	Host        string
	Port        int
	DB          string
	User        string
	Password    string
	Ctx         context.Context
	MinPoolSize uint64
	MaxPoolSize uint64
	MaxIdleTime time.Duration
	IsAdmin     bool
}

func (o *Option) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s%s", o.User, o.Password, o.Host, o.Port, o.DB, utilities.Ternary(o.IsAdmin, "?authSource=admin", ""))
}
