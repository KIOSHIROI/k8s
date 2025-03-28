package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"layer-scheduler/layer"

	klog "k8s.io/klog/v2" // 导入 Kubernetes 的日志库，用于记录日志信息
	"k8s.io/kubernetes/cmd/kube-scheduler/app" // 导入 Kubernetes 调度器的应用程序包
)

func main() {
	// 定义等待时间，设置为 10 秒，用于控制监听器的轮询间隔
	waitTs := 10 * time.Second
	// 定义本地缓存文件的路径，用于存储监听器获取的数据
	localCacheFile := "cache.json"

	// 记录日志，表示启动监听器
	klog.Infof("启动监听器")

	// 调用 layer.NewRegistry 函数，创建一个新的注册中心实例
	// 参数包括监听器的地址、认证信息等（此处为空字符串，可能不需要认证）
	reg, err := layer.NewRegistry(
		"http://localhost:5000", // 监听器的地址
		"",                      // 用户名（可能为空）
		"",                      // 密码（可能为空）
	)
	if err != nil {
		// 如果创建注册中心失败，记录错误日志并退出程序
		klog.Fatalf("监听器启动失败, err: %s", err)
		os.Exit(2) // 使用非零退出码表示错误
	}

	// 创建一个可取消的上下文，用于控制监听器的生命周期
	ctx, cancel := context.WithCancel(context.Background())
	
	// 启动监听器的 Watcher 方法，开始监听指定地址的数据变化
	// 参数包括等待时间、缓存文件路径和上下文
	go reg.Watcher(waitTs, localCacheFile, ctx)

	// 记录日志，表示启动调度器
	klog.Infof("启动调度器")

	// 创建 Kubernetes 调度器命令
	// 使用 layer.Name 和 layer.New 作为插件名称和构造函数，将自定义的 layer 插件集成到调度器中
	command := app.NewSchedulerCommand(
		app.WithPlugin(layer.Name, layer.New),
	)

	// 注释掉的代码，可能是用于测试的延迟逻辑
	// time.Sleep(100 * time.Second)

	// 确保在程序退出时取消上下文，停止监听器的 Watcher 方法
	defer cancel()

	// 执行调度器命令
	// 如果执行失败，记录错误日志并退出程序
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err) // 将错误信息输出到标准错误流
		os.Exit(1)                      // 使用非零退出码表示错误
	}
}