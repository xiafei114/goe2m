
# [goe2m](https://github.com/xiafei114/goe2m)

--------

#### excel 转 struct 工具,可以将excel 自动生成golang sturct结构，带大驼峰命名规则。带json标签
#### 此项目为例生成 [gin-admin](https://github.com/xiafei114/gin-admin) 的entity model bll ctl schema 开发

--------

## 1. 通过当前目录 config.yml 文件配置默认配置项
```
in_file_path : ./example.xlsx  # 输入文件
out_dir : ./target  # 输出目录
project_name : new_project
```

## 2. 可以使用命令行工具更新配置项

```
./goe2m -i=./example.xlsx -o=./target
```

