package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
	"zhaoxin2025/model"
	"zhaoxin2025/logger"
)

// RefreshAccessToken 是由 cron 调度的函数，用于刷新 AccessToken。
func RefreshAccessToken() {
	println("定时任务: 开始尝试刷新 AccessToken... ")

	// 获取 AccessToken
	err := GetAccessToken()
	if err != nil {
		// 如果刷新失败
		logger.DatabaseLogger.Errorf("[Cron] 刷新 AccessToken 失败: %v", err)
		println("定时任务: 刷新 AccessToken 失败。将在 5 分钟后重试。", err)

		// 启动一个 goroutine，在 5 分钟后重试
		go func() {
			time.Sleep(5 * time.Minute)
			fmt.Println("重试任务: 5 分钟后尝试刷新 AccessToken... ")
			err := GetAccessToken()
			if err != nil {
				fmt.Printf("重试任务: 刷新 AccessToken 仍然失败: %v\n", err)
			} else {
				fmt.Println("重试任务: AccessToken 刷新成功。")
			}
		}()
		return
	}

	// 如果刷新成功
	println("定时任务: AccessToken 刷新成功。")
}

type Queue struct {
	Queue chan model.Stu
	wg    sync.WaitGroup
	once  sync.Once
}

var RegisterQueue = Queue{Queue: make(chan model.Stu, 100)}
var ResultQueue = Queue{Queue: make(chan model.Stu, 100)}
var TimeQueue = Queue{Queue: make(chan model.Stu, 100)}

func (q *Queue) AddMessage(stu model.Stu) error {
	select {
	case q.Queue <- stu: // 尝试将消息发送到 channel
		fmt.Printf("Producer: Published message OpenID: %s\n", stu.OpenID)
		return nil
	case <-time.After(time.Second): // 如果1秒内无法发送（队列已满），则超时
		return fmt.Errorf("producer: queue is full, failed to publish message NetID: %s", stu.OpenID)
	}
}

func (q *Queue) ConsumeMessage(handler func(model.Stu) error) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()          // 协程退出时减少等待组计数
		for msg := range q.Queue { // 循环从 channel 中接收消息，直到 channel 关闭
			if err := handler(msg); err != nil {
				logger.DatabaseLogger.Errorf("[Cron] 消费者发送信息失败: OpenID: %s, error: %v", msg.OpenID, err)
				fmt.Printf("Consumer: Failed to process message OpenID: %s, error: %v\n", msg.OpenID, err)
			}
			fmt.Printf("Consumer: Processed message OpenID: %s\n", msg.OpenID)
		}
		fmt.Printf("Consumer: Finished processing messages\n")
	}()
}

// Send 是由 cron 调度的函数，用于获取和发送订阅消息。
func Send() {
	TimeQueue.ConsumeMessage(SendTime)
	ResultQueue.ConsumeMessage(SendResult)
	RegisterQueue.ConsumeMessage(SendRegister)
	go func() {
		c := cron.New(cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		))
		// 每10分钟执行一次，获取需要发送的订阅消息
		if _, err := c.AddFunc("@every 10m", func() {
			var record []model.Stu
			if err := model.DB.Model(&model.Stu{}).Where("message > 0").Preload("Interv").Find(&record).Error; err != nil {
				fmt.Printf("定时任务: 查询学生信息失败: %v\n", err)
				return
			}
			for _, stu := range record {
				message := stu.Message
				if (message<<1)&1 == 1 && time.Now().Add(20*time.Minute).Before(stu.Interv.Time) && time.Now().Add(30*time.Minute).After(stu.Interv.Time) {
					err := TimeQueue.AddMessage(stu)
					if err != nil {
						fmt.Printf("添加面试时间订阅消息失败: %v\n", err)
					}
				}
			}
		}); err != nil {
			fmt.Printf("定时任务: 添加定时任务失败: %v\n", err)
		}
		c.Start()
	}()
	go func() {
		c := cron.New(cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		))
		// 每10分钟执行一次，获取并发送面试时间消息
		if _, err := c.AddFunc("@every 10m", func() {
			fd, err := AliyunSendItvTimeMsg()
			if err != nil {
				fmt.Printf("定时任务: 获取面试时间订阅消息失败: %v\n", err)
				return
			}
			for _, f := range fd {
				if f.ErrCode != 0 {
					fmt.Printf("定时任务: 发送面试时间订阅消息失败, NetID: %s, ErrCode: %d\n", f.NetID, f.ErrCode)
				} else {
					stu := model.Stu{NetID: f.NetID}
					if err := ResultQueue.AddMessage(stu); err != nil {
						fmt.Printf("定时任务: 发送面试时间订阅消息失败: %v\n", err)
					}
				}
			}
		}); err != nil {
			fmt.Printf("定时任务: 添加发送面试时间任务失败: %v\n", err)
		}
		c.Start()
	}()
}
