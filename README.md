如果依赖其他包，需要在根目录执行go mod tidy; pkg目录下的包内不要包含go.mod
go-pkg包会有一些配置依赖，如需调整，请看gopkg/config.go