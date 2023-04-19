# 说明文档

---

## swag

- 常见报错链接
  `https://www.jianshu.com/p/b920a3f74cf2`

  - Fetch error Internal Server Error doc.json
    - 错误1: Fetch errorInternal Server Error doc.json 的时候。是因为没有导入依赖包doc的问题 [见router.go文件]
    - `import _ "project_name/docs"`
   - 错误2: swag不是命令错误
      - gopath环境变量
      - 用`go get -u github.com/swaggo/swag/cmd/swag`, 而不是用 `go install github.com/swaggo/swag/cmd/swag@latest`, Starting in Go 1.17, installing executables with go get is deprecated. go install may be used instead:

- github官方文件+使用说明
  `https://github.com/swaggo/gin-swagger`