康威生命游戏（并行版）

这个目录里是一份用 Go + Ebiten/Ebitengine v2 写的 康威生命游戏。
它的思路来源于网上那篇经典的生命游戏文章（单线程版本），但是我做了三件事：

换成了现在在用的 github.com/hajimehoshi/ebiten/v2 的写法；

还是按照原来的生命游戏规则（B3 / S23）来算；

把“算下一代”这一段做成了并行的，用多个 goroutine 一起算；

绘图保持单线程（因为 Ebiten 的图片不是并发安全的）。

1. 目录结构
life/
├── main.go      # 核心代码：游戏结构体、并行更新、绘图
├── go.mod       # 模块声明：SETU/Con_dev/game_of_life
├── go.sum       # go 自动生成的依赖锁定文件
└── README.md
这个实现是参考了网上那篇 “Implementing Conway’s Game of Life in Go” 的思路，但是：

我没有照抄它老的 ebiten API，而是全改成了 ebiten/v2；

我把它包成了一个 ebiten.Game 结构；

我把它的那段“两层 for 循环算下一代”改成了 并行版。

所以可以说：规则是同一套，写法是新版的，实现是并行的。
