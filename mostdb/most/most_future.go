package most

import (
	"errors"
)

type FutureFunc interface {
	Execute() *MostResult
}

type FutureTask func() *MostResult

type Futures struct {
	Result         *MostResult      // 执行结果
	Error          error            // 执行错误
	done           bool             // 是否执行完成
	resultChan     chan bool        // 执行结果通道
	taskResultChan chan *MostResult // 任务执行结果通道
	tasks          []FutureTask     // 任务
}

func NewFuture() *Futures {
	return &Futures{
		resultChan:     make(chan bool, 1),         // 初始化执行结果管道
		taskResultChan: make(chan *MostResult, 64), // 初始化任务执行结果管道
	}
}

// 添加任务
func (f *Futures) AddTask(task FutureTask) *Futures {
	if !f.done {
		f.tasks = append(f.tasks, task)
	}
	return f
}

// 是否已执行完
func (f *Futures) IsDone(task FutureTask) bool {
	return f.done
}

// 异步执行全部任务，等待过半一致的结果
func (f *Futures) WaitForMostResult() {

	// 避免重复执行影响
	if f.done {
		return
	}

	// 异步执行
	size := len(f.tasks)
	for i := 0; i < size; i++ {
		fn := f.tasks[i]
		go func() {
			rs := fn()
			if !f.done {
				f.taskResultChan <- rs
			}
		}()
	}

	// 过半确认
	half := size / 2 // 向下取整
	mapCnt := make(map[string](int))
	cnt := 0
	for {
		rs := <-f.taskResultChan
		cnt++

		if !f.done && rs != nil && rs.Success {
			key := rs.Result.ToJson()
			n := mapCnt[key]
			n++
			if n > half {
				f.Result = rs        // 过半一致的结果
				f.done = true        // 完成
				f.resultChan <- true // 过半一致即可
				break
			}
			mapCnt[key] = n
		}

		if cnt >= size {
			if !f.done {
				f.Result = nil        // 结果nil
				f.done = true         // 完成
				f.resultChan <- false // 完成
				f.Error = errors.New("无过半一致的结果")
			}
			break
		}
	}

	defer close(f.taskResultChan)
	defer close(f.resultChan)

	<-f.resultChan
}
