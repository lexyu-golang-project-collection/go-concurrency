package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Step struct {
	id      int
	data    interface{}
	waitFor chan bool
}

type StepProcessor struct {
	lastSuccessful int
	results        []string
}

func generateSequence() (<-chan Step, <-chan error) {
	steps := make(chan Step)
	errors := make(chan error, 1)

	testData := []string{
		"初始化系統",
		"載入配置",
		"連接數據庫",
		"啟動服務",
		"完成設定",
	}

	go func() {
		defer close(steps)
		defer close(errors)

		for i, data := range testData {
			// 為每個步驟創建等待channel
			waitForIt := make(chan bool)

			step := Step{
				id:      i + 1,
				data:    data,
				waitFor: waitForIt,
			}

			steps <- step

			// 等待步驟確認或超時
			select {
			case <-waitForIt:
				fmt.Printf("步驟 %d 已確認完成\n", i+1)
			case <-time.After(5 * time.Second):
				errors <- fmt.Errorf("步驟 %d 等待超時", i+1)
				return
			}
		}
	}()

	return steps, errors
}

func (sp *StepProcessor) processStep(step Step) error {
	fmt.Printf("正在處理步驟 %d: %v\n", step.id, step.data)

	time.Sleep(300 * time.Millisecond)

	if rand.Float32() < 0.1 {
		return fmt.Errorf("處理步驟 %d 時發生錯誤", step.id)
	}

	sp.lastSuccessful = step.id
	sp.results = append(sp.results, fmt.Sprintf("步驟 %d (%v) 處理完成", step.id, step.data))

	return nil
}

func (sp *StepProcessor) processWithRetry(step Step, maxRetries int) bool {
	for retry := 0; retry <= maxRetries; retry++ {
		if retry > 0 {
			fmt.Printf("重試步驟 %d (第 %d 次)\n", step.id, retry)
			time.Sleep(500 * time.Millisecond) // 重試前稍微等待
		}

		err := sp.processStep(step)
		if err == nil {
			return true // 處理成功
		}

		if retry < maxRetries {
			fmt.Printf("錯誤: %v (已重試 %d/%d)\n", err, retry, maxRetries)
		} else {
			fmt.Printf("步驟 %d 重試次數超過上限，放棄處理\n", step.id)
		}
	}
	return false // 所有重試都失敗
}

func NewStepProcessor() *StepProcessor {
	return &StepProcessor{
		lastSuccessful: 0,
		results:        make([]string, 0),
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	processor := NewStepProcessor()
	steps, errors := generateSequence()
	maxRetries := 3

	fmt.Println("開始處理序列...")
	fmt.Println("====================")

	for {
		select {
		case step, ok := <-steps:
			if !ok {
				fmt.Println("====================")
				fmt.Println("序列處理完成！結果：")
				for _, result := range processor.results {
					fmt.Println(result)
				}
				return
			}

			// 使用重試機制處理步驟
			success := processor.processWithRetry(step, maxRetries)

			// 無論成功失敗，都要發送確認信號避免超時
			step.waitFor <- success

		case err := <-errors:
			fmt.Printf("序列錯誤: %v\n", err)
			return
		}
	}
}
