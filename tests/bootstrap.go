package tests

import (
	"github.com/goal-web/application"
	"github.com/goal-web/config"
	"github.com/goal-web/contracts"
	"github.com/goal-web/database"
	"github.com/goal-web/migration"
	"github.com/goal-web/redis"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
)

func initApp(path ...string) contracts.Application {
	app := application.Singleton()

	// 设置异常处理器
	app.Singleton("exceptions.handler", func() contracts.ExceptionHandler {
		return exceptions.DefaultExceptionHandler{}
	})

	app.RegisterServices(
		config.NewService(config.NewToml(config.File("test.toml")), map[string]contracts.ConfigProvider{
			"database": func(env contracts.Env) any {
				return database.Config{
					Default: env.StringOptional("db.connection", "mysql"),
					Connections: map[string]contracts.Fields{
						"mysql": {
							"driver":          "mysql",
							"host":            env.GetString("db.host"),
							"port":            env.GetString("db.port"),
							"database":        env.GetString("db.database"),
							"username":        env.GetString("db.username"),
							"password":        env.GetString("db.password"),
							"unix_socket":     env.GetString("db.unix_socket"),
							"charset":         env.StringOptional("db.charset", "utf8mb4"),
							"collation":       env.StringOptional("db.collation", "utf8mb4_unicode_ci"),
							"prefix":          env.GetString("db.prefix"),
							"strict":          env.GetBool("db.struct"),
							"max_connections": env.GetInt("db.max_connections"),
							"max_idles":       env.GetInt("db.max_idles"),
						},
					},
				}
			},
			"app": func(env contracts.Env) any {
				return application.Config{
					Name:     env.GetString("app.name"),
					Debug:    env.GetBool("app.debug"),
					Timezone: env.GetString("app.timezone"),
					Env:      env.GetString("app.env"),
					Locale:   env.GetString("app.locale"),
					Key:      env.GetString("app.key"),
				}
			},
			"redis": func(env contracts.Env) any {
				return redis.Config{
					Default: utils.StringOr(env.GetString("redis.default"), "default"),
					Stores: map[string]contracts.Fields{
						"default": {
							"network":  env.GetString("redis.network"),
							"host":     env.GetString("redis.host"),
							"port":     env.GetString("redis.port"),
							"username": env.GetString("redis.username"),
							"password": env.GetString("redis.password"),
							"db":       env.GetInt64("redis.db"),
							"retries":  env.GetInt64("redis.retries"),
						},
						"cache": {
							"network":  env.GetString("redis.cache.network"),
							"host":     env.GetString("redis.cache.host"),
							"port":     env.GetString("redis.cache.port"),
							"username": env.GetString("redis.cache.username"),
							"password": env.GetString("redis.cache.password"),
							"db":       env.GetInt64("redis.cache.db"),
							"retries":  env.GetInt64("redis.cache.retries"),
						},
					},
				}
			},
		}),
		database.NewService(),
		redis.NewService(),
		NewService(),
		migration.NewService(),
		//&http.serviceProvider{RouteCollectors: []any{
		//	// 路由收集器
		//	routes.V1Routes,
		//}},
	)

	go func() {
		if errors := app.Start(); len(errors) > 0 {
			logs.WithField("errors", errors).Fatal("goal 启动异常!")
		}
	}()
	return app
}
