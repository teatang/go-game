package tetris

import "math/rand"

// ============================================
// 常量定义 - 游戏参数配置
// ============================================

const (
	BoardWidth  = 10 // 游戏面板宽度（格子数）
	BoardHeight = 20 // 游戏面板高度（格子数）
)

// ============================================
// 形状定义 - 7种经典俄罗斯方块
// ============================================
// 每个形状用一个二维数组表示
// 1 表示方块存在，0 表示空白
var Shapes = [][][]int{
	{{1, 1, 1, 1}},             // I - 青色长条
	{{1, 1}, {1, 1}},           // O - 黄色正方形
	{{0, 1, 0}, {1, 1, 1}},     // T - 紫色T形
	{{0, 1, 1}, {1, 1, 0}},     // S - 绿色S形
	{{1, 1, 0}, {0, 1, 1}},     // Z - 红色Z形
	{{1, 0, 0}, {1, 1, 1}},     // J - 蓝色J形
	{{0, 0, 1}, {1, 1, 1}},     // L - 橙色L形
}

// Colors 每种形状对应的显示颜色
var Colors = []string{
	"cyan",    // I
	"yellow",  // O
	"fuchsia", // T
	"lime",    // S
	"red",     // Z
	"navy",    // J
	"olive",   // L
}

// ============================================
// RNG 随机数接口
// ============================================
// 用于支持测试时的依赖注入

type RNG interface {
	Intn(n int) int
}

// randRNG 标准库随机数的实现
type randRNG struct{}

// Intn 返回 [0, n) 范围内的随机整数
func (r randRNG) Intn(n int) int {
	return rand.Intn(n)
}

// ============================================
// Game 结构体 - 核心游戏状态
// ============================================

type Game struct {
	// 游戏面板：0表示空，非0表示该位置的方块颜色索引+1
	board [][]int

	// 当前方块信息
	currPiece int     // 当前方块的形状索引 (0-6)
	currShape [][]int // 当前方块的形状数据
	nextPiece int     // 下一个方块的形状索引（+1存储，0表示未设置）

	// 方块在面板上的位置
	pieceX int // 方块左上角在面板的X坐标
	pieceY int // 方块左上角在面板的Y坐标

	// 游戏状态
	score    int  // 当前得分
	lines    int  // 消除的总行数
	level    int  // 当前等级（影响下落速度）
	paused   bool // 游戏是否暂停
	gameOver bool // 游戏是否结束

	// 依赖组件
	rng RNG // 随机数生成器
}

// ============================================
// 工厂方法
// ============================================

// NewGame 创建并初始化新的俄罗斯方块游戏
func NewGame() *Game {
	// 初始化游戏面板（boardHeight 行 x boardWidth 列）
	board := make([][]int, BoardHeight)
	for i := range board {
		board[i] = make([]int, BoardWidth)
	}

	return &Game{
		board:   board,
		pieceX:  BoardWidth/2 - 1,
		pieceY:  0,
		level:   1,
		rng:     randRNG{},
	}
}

// ============================================
// 核心游戏逻辑 - 方块生成与控制
// ============================================

// spawnPiece 生成并放置新方块
//
// 逻辑说明：
// 1. 如果有预存的 nextPiece，使用它作为当前方块
// 2. 否则随机生成一个方块
// 3. 在面板中央上方放置方块
// 4. 预生成下一个方块
// 5. 检查方块是否还能放置（无法放置则游戏结束）
func (g *Game) spawnPiece() {
	// 选择当前方块
	if g.nextPiece == 0 {
		g.currPiece = g.rng.Intn(len(Shapes))
	} else {
		// nextPiece 存储的是颜色索引+1，所以需要减1
		g.currPiece = g.nextPiece - 1
	}
	g.currShape = Shapes[g.currPiece]

	// 设置方块位置（居中）
	g.pieceX = BoardWidth/2 - len(g.currShape[0])/2
	g.pieceY = 0

	// 预生成下一个方块（+1 是因为0表示"未设置"状态）
	g.nextPiece = g.rng.Intn(len(Shapes)) + 1

	// 检查碰撞：如果新方块无法放置，游戏结束
	if g.collides() {
		g.gameOver = true
	}
}

// collides 碰撞检测
// 检测当前方块是否与边界或其他已锁定方块发生碰撞
//
// 碰撞检测规则：
// 1. 检查方块的每个单元格是否超出面板边界
// 2. 检查方块的每个单元格是否与已锁定的方块重叠
// 注意：pieceY < 0 的情况不视为碰撞（方块刚从顶部出现）
func (g *Game) collides() bool {
	for y, row := range g.currShape {
		for x, cell := range row {
			if cell == 1 {
				boardX := g.pieceX + x
				boardY := g.pieceY + y

				// 边界检测
				if boardX < 0 || boardX >= BoardWidth || boardY >= BoardHeight {
					return true
				}
				// 已锁定方块检测
				if boardY >= 0 && g.board[boardY][boardX] != 0 {
					return true
				}
			}
		}
	}
	return false
}

// rotate 旋转当前方块（顺时针90度）
//
// 旋转算法：
// 1. 创建一个新的矩阵，行列互换
// 2. 通过 formula: rotated[x][rows-1-y] = cell 实现顺时针旋转
// 3. 如果旋转后发生碰撞，则回滚到原形状
func (g *Game) rotate() {
	rows := len(g.currShape)
	cols := len(g.currShape[0])

	// 创建旋转后的新矩阵
	rotated := make([][]int, cols)
	for i := range rotated {
		rotated[i] = make([]int, rows)
	}

	// 顺时针旋转
	for y, row := range g.currShape {
		for x, cell := range row {
			rotated[x][rows-1-y] = cell
		}
	}

	// 尝试应用旋转，碰撞则回滚
	oldShape := g.currShape
	g.currShape = rotated
	if g.collides() {
		g.currShape = oldShape
	}
}

// move 尝试移动方块
// dx, dy 分别表示X和Y方向的位移
// 返回值表示移动是否成功
func (g *Game) move(dx, dy int) bool {
	g.pieceX += dx
	g.pieceY += dy
	if g.collides() {
		g.pieceX -= dx
		g.pieceY -= dy
		return false
	}
	return true
}

// drop 让方块下落一格
// 返回值：如果方块落地返回 false，否则返回 true
func (g *Game) drop() bool {
	if !g.move(0, 1) {
		// 方块落地，执行锁定、消除、生成新方块
		g.lockPiece()
		g.clearLines()
		g.spawnPiece()
		return false
	}
	return true
}

// ============================================
// 游戏逻辑 - 方块锁定与消除
// ============================================

// lockPiece 将当前方块锁定在游戏板上
// 将方块的每个单元格标记到 board 数组中
func (g *Game) lockPiece() {
	for y, row := range g.currShape {
		for x, cell := range row {
			if cell == 1 {
				boardY := g.pieceY + y
				boardX := g.pieceX + x
				if boardY >= 0 && boardY < BoardHeight && boardX >= 0 && boardX < BoardWidth {
					// +1 是因为 board 中 0 表示空白
					g.board[boardY][boardX] = g.currPiece + 1
				}
			}
		}
	}
}

// clearLines 检测并消除已满的行
//
// 算法：
// 1. 从底部向上扫描每一行
// 2. 如果某行没有空白单元格，则该行已满
// 3. 删除已满行，其上方的所有行下移一行
// 4. 在顶部添加新空白行
//
// 得分计算：
// 消除1行: 100 * level
// 消除2行: 300 * level
// 消除3行: 500 * level
// 消除4行: 800 * level
func (g *Game) clearLines() {
	linesCleared := 0

	// 从底部向上扫描
	for y := BoardHeight - 1; y >= 0; y-- {
		// 检查当前行是否已满
		complete := true
		for x := 0; x < BoardWidth; x++ {
			if g.board[y][x] == 0 {
				complete = false
				break
			}
		}

		if complete {
			linesCleared++

			// 上方行下移
			for removeY := y; removeY > 0; removeY-- {
				for x := 0; x < BoardWidth; x++ {
					g.board[removeY][x] = g.board[removeY-1][x]
				}
			}
			// 清空顶部行
			for x := 0; x < BoardWidth; x++ {
				g.board[0][x] = 0
			}
			// 重新检查当前行（因为上方行下移了）
			y++
		}
	}

	// 更新分数和等级
	if linesCleared > 0 {
		g.lines += linesCleared

		// 得分表
		scoreTable := []int{0, 100, 300, 500, 800}
		g.score += scoreTable[linesCleared] * g.level

		// 每消除10行升一级
		g.level = g.lines/10 + 1
	}
}

// ============================================
// 幽灵方块功能（Ghost Piece）
// ============================================

// getGhostPosition 计算当前方块最终会落到的位置
// 返回值：幽灵方块的X和Y坐标
//
// 算法：
// 1. 从当前位置开始，模拟方块持续下落
// 2. 每次检查是否可以继续下落
// 3. 当无法下落时，返回最后的有效位置
func (g *Game) getGhostPosition() (int, int) {
	ghostY := g.pieceY

	for {
		canMove := true
		for y, row := range g.currShape {
			for x, cell := range row {
				if cell == 1 {
					boardX := g.pieceX + x
					boardY := ghostY + y + 1

					// 检查是否到底或撞到方块
					if boardY >= BoardHeight {
						canMove = false
						break
					}
					if boardY >= 0 && g.board[boardY][boardX] != 0 {
						canMove = false
						break
					}
				}
			}
			if !canMove {
				break
			}
		}
		if !canMove {
			break
		}
		ghostY++
	}

	return g.pieceX, ghostY
}

// ============================================
// 游戏控制
// ============================================

// reset 重置游戏到初始状态
func (g *Game) reset() {
	// 清空面板
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = 0
		}
	}

	// 重置状态
	g.score = 0
	g.lines = 0
	g.level = 1
	g.nextPiece = 0
	g.paused = false
	g.gameOver = false

	// 生成第一个方块
	g.spawnPiece()
}

// getDropInterval 获取当前等级对应的下落间隔
// 等级越高，下落越快
// 最小间隔限制为50毫秒
func (g *Game) getDropInterval() int {
	interval := 500 / g.level
	if interval < 50 {
		interval = 50
	}
	return interval
}
