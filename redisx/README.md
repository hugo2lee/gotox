<!--
 * @Author: hugo  
 * @Date: 2024-04-02 10:25
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-02 19:37
 * @FilePath: \gotox\redisx\README.md
 * @Description: 
 * 
 * Copyright (c) 2024 by hugo, All Rights Reserved. 
-->
# goto reidsx
对 go-redis 的使用整合

## redis链接
url="redis://:root@localhost:6379/0"

## 关闭redisCli
使用 func (c *redisCli) Close(ctx context.Context, wg *sync.WaitGroup) 进行优雅关闭