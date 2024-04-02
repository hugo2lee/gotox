<!--
 * @Author: hugo  
 * @Date: 2024-04-02 10:25
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-02 14:23
 * @FilePath: \gotox\logx\README.md
 * @Description: 
 * 
 * Copyright (c) 2024 by hugo, All Rights Reserved. 
-->
# goto logx
对 zap 的使用整合

## 配置读取逻辑
使用logger结构体，内含configx, 读取log文件存放路径及日志记录级别，强依赖于config mode

### 内部逻辑
内部使用 zap 进行日志记录，lumberjack 进行日志切割

### 扩充logger方法
type logCustom struct {
	*logCli
}

func (l *logCustom) Fatal(msg string, args ...any) {
	l.logger.Sugar().Fatalf(msg, args...)
}

### 替换logger实现
logCli 实现了 Logger interface，你也可以自己实现接口无缝替换