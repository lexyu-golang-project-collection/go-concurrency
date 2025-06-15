# Overview

1. 基本模式 (Basic Patterns)
    - Generator Pattern (1-boring, 3-generator) ✨
    - Fan-in Pattern (4-fanin) ✨
    - Select Pattern (6-select-timeout) ✨
    - Quit Signal Pattern (7-quit-signal)

2. 通道管理 (Channel Management)
    - Daisy Chain Pattern (8-daisy-chan)
    - Ring Buffer Channel Pattern (17-ring-buffer-channel) -> `container/ring`

3. 並行控制 (Concurrency Control)
    - Bounded Parallelism Pattern (15-bounded-parallelism) ✨
    - Worker Pool Pattern (18-worker-pool) ✨

4. 同步模式 (Synchronization Patterns)
    - Ping-Pong Pattern (13-adv-pingpong)

5. 錯誤處理和恢復 (Error Handling and Recovery)
    - Restore Sequence Pattern (5-restore-sequence)

6. 上下文和取消 (Context and Cancellation)
    - Context Pattern (16-context) ✨

7. 訂閱模式 (Subscription Pattern)
    - Advanced Subscription Pattern (14-adv-subscription)

8. Google 搜索模式 (Google Search Patterns)
    - Google Search 1.0 (9-google1.0)
    - Google Search 2.0 (10-google2.0)
    - Google Search 2.1 (11-google2.1)
    - Google Search 3.0 (12-google3.0)


## Summary

|             | Usages       |
| ------------------- | -------------- |
| Context Pattern     | 管理取消與逾時 |
| Worker Pool Pattern | 背景任務、並發 |
| Select Pattern      | 超時、多路選擇 |
| Generator + Fan-in  | 資料流處理     |
| Bounded Parallelism | 控制併發數     |
