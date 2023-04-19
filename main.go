package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"threads-service/global"
	"threads-service/internal/model"
	"threads-service/internal/routers"
	"threads-service/pkg/logger"
	"threads-service/pkg/setting"
	"threads-service/pkg/tracer"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	port    string
	runMode string
	config  string

	isVersion    bool
	buildTime    string
	buildVersion string
	gitCommitID  string
)

func init() {

	_ = setupFlag()

	// setupFlag 必须要比 setupSetting 先执行

	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setUpSetting err:%v", err)
		panic(err)
	}
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupEngine err:%v", err)
	}
	steupLogger()
	log.Println("init.setupLogger 日志初始化完成")
	err = setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

// @title			战舰世界论坛
// @version		1.0
// @description	战舰世界论坛, 畅言好玩战舰
// @termsOfService	发表意见平台技术条款
func main() {
	if isVersion {
		fmt.Printf("build_time: %s\n", buildTime)
		fmt.Printf("build_version: %s\n", buildVersion)
		fmt.Printf("git_commit_id: %s\n", gitCommitID)
	}

	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe() err: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("server exiting")
}

func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要是有的配置文件路径")

	flag.BoolVar(&isVersion, "version", false, "是否打印编译信息")
	flag.Parse()

	return nil
}

// 在初始化函数中调用该函数
func setupSetting() error {
	settings, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		log.Fatal("settings, err := setting.NewSetting()", err)
		return err
	}

	if port != "" {
		global.ServerSetting.HttpPort = port
	}

	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}

	err = settings.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		log.Fatal("err = settings.ReadSection(\"Server\", &global.ServerSetting)", err)
		return err
	}
	log.Println("读取的服务器配置:", global.ServerSetting)

	err = settings.ReadSection("App", &global.AppSetting)
	if err != nil {
		log.Fatal("err = settings.ReadSection(\"App\", &global.AppSetting)", err)
		return err
	}
	log.Println("读取的应用配置:", global.AppSetting)

	err = settings.ReadSection("DataBase", &global.DataBaseSetting)
	if err != nil {
		log.Fatal("err = settings.ReadSection(\"DataBase\", &global.DataBaseSetting)", err)
		return err
	}
	log.Println("读取的数据库配置:", global.DataBaseSetting)

	err = settings.ReadSection("JWT", &global.JwtSetting)
	if err != nil {
		log.Fatal("err=settings.ReadSection(\"JWT\", &global.JwtSetting):", err)
		return err
	}

	err = settings.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		log.Fatal("err=settings.ReadSection(\"Email\", &global.EmailSetting):", err)
		return err
	}

	global.AppSetting.DefaultContextTimeout *= time.Second
	global.JwtSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DataBaseSetting)
	return err
}

func steupLogger() {
	lj := &lumberjack.Logger{
		Filename:  global.AppSetting.LogPath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExtentionName,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}
	global.Logger = logger.NewLogger(lj, "", log.LstdFlags).WithCaller(2)
	log.Println("日志位置:", lj.Filename)
}

func setupTracer() error {
	jaegerTarcer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6832")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTarcer
	return nil
}
