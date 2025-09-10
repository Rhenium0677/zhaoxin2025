package service

import (
	"github.com/robfig/cron/v3"
	// "sync"
	"time"
	"zhaoxin2025/logger"
	"zhaoxin2025/model"
)

// RefreshAccessToken 是由 cron 调度的函数，用于刷新 AccessToken。
func RefreshAccessToken() {
	logger.DatabaseLogger.Infof("定时任务: 开始尝试刷新 AccessToken... ")

	// 获取 AccessToken
	err := GetAccessToken()
	if err != nil {
		// 如果刷新失败
		logger.DatabaseLogger.Errorf("[Cron] 刷新 AccessToken 失败: %v", err)

		// 启动一个 goroutine，在 5 分钟后重试
		go func() {
			time.Sleep(5 * time.Minute)
			err := GetAccessToken()
			if err != nil {
				logger.DatabaseLogger.Errorf("[Cron] Retry: 刷新 AccessToken 仍然失败: %v\n", err)
			}
		}()
		return
	}

	// 如果刷新成功
	logger.DatabaseLogger.Infof("[Cron] AccessToken 刷新成功。")
}

// type Queue struct {
// 	Queue chan model.Stu
// 	wg    sync.WaitGroup
// 	once  sync.Once
// }

// var RegisterQueue = Queue{Queue: make(chan model.Stu, 100)}
// var ResultQueue = Queue{Queue: make(chan model.Stu, 100)}
// var TimeQueue = Queue{Queue: make(chan model.Stu, 100)}

// func (q *Queue) AddMessage(stu model.Stu) error {
// 	select {
// 	case q.Queue <- stu: // 尝试将消息发送到 channel
// 		logger.DatabaseLogger.Infof("[Cron] Producer published message for %s\n", stu.OpenID)
// 		return nil
// 	case <-time.After(time.Second): // 如果1秒内无法发送（队列已满），则超时
// 		logger.DatabaseLogger.Errorf("[Cron] Producer failed to publish message for %s\n", stu.OpenID)
// 		return fmt.Errorf("producer: queue is full, failed to publish message NetID: %s", stu.OpenID)
// 	}
// }

// func (q *Queue) ConsumeMessage(handler func(model.Stu) error) {
// 	q.wg.Add(1)
// 	go func() {
// 		defer q.wg.Done()          // 协程退出时减少等待组计数
// 		for msg := range q.Queue { // 循环从 channel 中接收消息，直到 channel 关闭
// 			if err := handler(msg); err != nil {
// 				logger.DatabaseLogger.Errorf("[Cron] Consumer failed to process message for %s, error: %v\n", msg.OpenID, err)
// 			}
// 			logger.DatabaseLogger.Infof("[Cron] Consumer processed message for %s", msg.OpenID)
// 		}
// 		logger.DatabaseLogger.Errorf("[Cron] Consumer: Finished processing messages\n")
// 	}()
// }

// Send 是由 cron 调度的函数，用于获取和发送订阅消息。
func Send() {
	// TimeQueue.ConsumeMessage(SendTime)
	// ResultQueue.ConsumeMessage(SendResult)
	// RegisterQueue.ConsumeMessage(SendRegister)
	go func() {
		c := cron.New(cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		))
		// 每10分钟执行一次，获取需要发送的订阅消息
		if _, err := c.AddFunc("@every 10m", func() {
			var record []model.Stu
			if err := model.DB.Model(&model.Stu{}).Where("message > 0").Preload("Interv").Find(&record).Error; err != nil {
				logger.DatabaseLogger.Errorf("[Cron] 查询学生信息失败: %v\n", err)
				return
			}
			for _, stu := range record {
				message := stu.Message
				if (message<<1)&1 == 1 && stu.Interv != nil && time.Now().Add(20*time.Minute).Before(stu.Interv.Time) && time.Now().Add(30*time.Minute).After(stu.Interv.Time) {
					err := SendTime(stu)
					if err != nil {
						logger.DatabaseLogger.Errorf("[Cron] 添加面试时间订阅消息失败: %v\n", err)
					}
				}
			}
		}); err != nil {
			logger.DatabaseLogger.Errorf("[Cron] 添加定时任务失败: %v\n", err)
		}
		c.Start()
	}()
	//go func() {
	//	c := cron.New(cron.WithChain(
	//		cron.SkipIfStillRunning(cron.DefaultLogger),
	//		cron.Recover(cron.DefaultLogger),
	//	))
	//	// 每10分钟执行一次，获取并发送面试时间消息
	//	if _, err := c.AddFunc("@every 10m", func() {
	//		fd, err := AliyunSendItvTimeMsg()
	//		if err != nil {
	//			logger.DatabaseLogger.Errorf("[Cron] 获取面试时间短信消息失败: %v\n", err)
	//			return
	//		}
	//		for _, f := range fd {
	//			if f.ErrCode != 0 {
	//				logger.DatabaseLogger.Errorf("[Cron] 发送面试时间短信消息失败, NetID: %s, ErrCode: %d\n", f.NetID, f.ErrCode)
	//			}
	//		}
	//	}); err != nil {
	//		logger.DatabaseLogger.Errorf("[Cron] 发送面试时间短信任务失败: %v\n", err)
	//	}
	//	c.Start()
	//}()
}
