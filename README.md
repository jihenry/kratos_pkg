## 一、目录定位
1. 存放和工程无关的非业务类代码
2. 此目录和公共库有所不同，这个是基于kratos框架，属于kratos模板工程中的pkg
3. 这个工程库定位是基于kratos统一公司基础中间件，这也是为什么pkg里面直接是redis和mysql命名
## 二、目录说明
```
.
├── README.md
├── cache //缓存
├── db //数据库
├── log //日志封装
└── util //通用函数库
```