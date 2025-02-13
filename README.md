
# excel转为go结构体和可以读json的工具



## 点击运行 exe 从 既可以读取和生成对应的文件



##生成文件夹
> 1.generateConfig
2.generateJsons
是包含excel转换的json文件和golang的文件

dirConf.json 是配置excel和json和go代码的输出地址

##生成后自己的项目如何调用？
可以使用generateConfig下生成的go代码里面的InitConfig() 就可以在自己项目使用读取配置。

## 配置表规则
1.第一行类型 2.第二行字段 3.第三行注释

