package echo

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Config struct {
	LogWriter io.Writer
	// Address 服务监听端口
	Address string `yaml:"address" json:"address" xml:"address"`
	// KeepAlive
	KeepAlive bool `yaml:"keep_alive" json:"keep_alive" xml:"keep_alive"`
	// AccessLogFormat http access 日志格式
	AccessLogFormat string `yaml:"access_log_format" json:"access_log_format" xml:"access_log_format"`
	// CORSConfig 跨域配置
	CORSConfig *middleware.CORSConfig `yaml:"cors" json:"cors_config" xml:"cors_config"`
	// SecureConfig 安全配置
	SecureConfig *middleware.SecureConfig `yaml:"secure" json:"secure_config" xml:"secure_config"`
	// RewriteRules 规则重写URL路径配置
	RewriteRules map[string]string `yaml:"rewrite_rules" json:"rewrite_rules" xml:"rewrite_rules"`
	// CookieDomain cookie域
	CookieDomain string `yaml:"cookie_domain" json:"cookie_domain" xml:"cookie_domain"`
	// StaticPrefix 静态路径前缀
	StaticPrefix string `yaml:"static_prefix" json:"static_prefix" xml:"static_prefix"`
	// StaticDir 静态文件目录
	StaticDir string `yaml:"static_dir" json:"static_dir" xml:"static_dir"`
	// BodyLimit body文件大小限制. Limit can be specified as 4x or 4xB, where x is one of the multiple from K, M, G, T or P.
	BodyLimit string `yaml:"body_limit" json:"body_limit" xml:"body_limit"`
	// TemplatePattern 视图模板渲染文件路径
	TemplatePattern string `yaml:"template_pattern" json:"template_pattern" xml:"template_pattern"`
}

const defaultAccessLogFormat = `${method} ${status} ${latency_human} ${id} ${host} ${uri} "${user_agent}" ${remote_ip} ${bytes_in} ${bytes_out}`

func New(cfg Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// 请求日志配置
	if len(cfg.AccessLogFormat) == 0 {
		cfg.AccessLogFormat = defaultAccessLogFormat
	}
	if cfg.LogWriter == nil {
		cfg.LogWriter = log.Writer()
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: cfg.AccessLogFormat + "\n",
		Output: cfg.LogWriter,
	}))

	// 路由重写配置
	if nil != cfg.RewriteRules {
		e.Use(middleware.Rewrite(cfg.RewriteRules))
	}

	// 跨域配置
	if cfg.CORSConfig != nil {
		e.Use(middleware.CORSWithConfig(*cfg.CORSConfig))
	} else {
		e.Use(middleware.CORS())
	}

	// 安全配置
	if cfg.SecureConfig != nil {
		e.Use(middleware.SecureWithConfig(*cfg.SecureConfig))
	} else {
		e.Use(middleware.Secure())
	}

	// 数据大小限制
	if len(cfg.BodyLimit) > 0 {
		e.Use(middleware.BodyLimit(cfg.BodyLimit))
	}

	// 静态文件配置
	if len(cfg.StaticPrefix) > 0 && len(cfg.StaticDir) > 0 {
		e.Static(cfg.StaticPrefix, cfg.StaticDir)
	}

	// 错误处理配置
	e.HTTPErrorHandler = ErrorHandler

	if len(cfg.TemplatePattern) > 0 {
		e.Renderer = &Template{
			templates: template.Must(template.ParseGlob(cfg.TemplatePattern)),
		}
	}

	return e
}

// ErrorHandler http 错误处理器
func ErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	_ = ctx.JSON(code, map[string]interface{}{
		"info": err.Error(),
	})
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
