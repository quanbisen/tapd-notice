package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"log"
	"tapd-notice/common/response"
	"tapd-notice/internal/service"
)

type Api struct {
	Context *gin.Context
	Orm     *gorm.DB
	Errors  error
}

func (e *Api) AddError(err error) {
	if e.Errors == nil {
		e.Errors = err
	} else if err != nil {
		e.Errors = fmt.Errorf("%v; %w", e.Errors, err)
	}
}

// MakeContext 设置http上下文
func (e *Api) MakeContext(c *gin.Context) *Api {
	e.Context = c
	return e
}

func (e *Api) MakeOrm() *Api {
	value, exist := e.Context.Get("db")
	if !exist {
		log.Println("数据库连接获取失败, ")
		e.AddError(errors.New("数据库连接获取失败"))
	}
	db, ok := value.(*gorm.DB)
	if !ok {
		log.Println("数据库对象类型转换失败")
		e.AddError(errors.New("数据库对象类型转换失败"))
	}
	e.Orm = db
	return e
}

func (e *Api) Bind(d interface{}, bindings ...binding.Binding) *Api {
	for i := range bindings {
		var err error
		if bindings[i] == nil {
			err = e.Context.ShouldBindUri(d)
		} else {
			err = e.Context.ShouldBindWith(d, bindings[i])
		}
		if err != nil && err.Error() == "EOF" && bindings[i].Name() != binding.JSON.Name() {
			log.Println("request body is not present anymore. ")
			err = nil
			continue
		}
		if err != nil {
			e.AddError(err)
			break
		}
	}
	return e
}

func (e *Api) MakeService(c *service.Service) *Api {
	c.Orm = e.Orm
	return e
}

// Error 通常错误数据处理
func (e *Api) Error(code int, err error, msg string) {
	response.Error(e.Context, code, err, msg)
}

// OK 通常成功数据处理
func (e *Api) OK(data interface{}, msg string) {
	response.OK(e.Context, data, msg)
}

// PageOK 分页数据处理
func (e *Api) PageOK(result interface{}, count int, pageIndex int, pageSize int, msg string) {
	response.PageOK(e.Context, result, count, pageIndex, pageSize, msg)
}

// Custom 兼容函数
func (e *Api) Custom(data gin.H) {
	response.Custom(e.Context, data)
}
