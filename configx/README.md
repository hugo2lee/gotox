<!--
 * @Author: hugo  
 * @Date: 2024-04-02 10:25
 * @LastEditors: hugo  
 * @LastEditTime: 2024-04-02 11:35
 * @FilePath: \gotox\config\README.md
 * @Description: 
 * 
 * Copyright (c) 2024 by hugo, All Rights Reserved. 
-->
# goto config
对 viper 的使用整合

## 配置读取逻辑
NewConfig - WithMode - WithPath - RUNMODE - CONFIGPATH - configKey - configValue
，使用了 WithX 就强制使用指定的环境或路径，即使配置错误

### 运行环境（模式）
读取系统变量 RUNMODE，可取值有 dev test prod, 默认值为 dev, 有with就使用with指定的

### 配置路径
默认值为 ./conf, 有with就使用with指定的

### 配置文件
本文件不可手动指定，运行环境与配置文件名强绑定，配置路径下的环境（模式）名.toml，默认值为{{RUNMODE}}.toml


### configFetch
里面定义了常用的配置项

### 扩充配置方法
type ConfigCustom struct {
	*ConfigCli
}

func (c *ConfigCustom) Key() string {
	return c.viper.GetString("key.value")
}