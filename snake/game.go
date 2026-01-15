package snake

import "math/rand"

// ============================================
// 常量定义 - 游戏参数配置
// ============================================

const (
	BoardWidth  = 20 // 游戏面板宽度（格子数）
	BoardHeight = 15 // 游戏面板高度（格子数）
	SpeedNormal = 150
	SpeedFast   = 80
)

// ============================================
// 枚举类型 - 方向定义
// ============================================
// 用于表示蛇的移动方向

type Direction int

const (
	Up Direction = iota // 向上移动
	Down               // 向下移动
	Left               // 向左移动
	Right              // 向右移动
)

// ============================================
// 基础数据结构
// ============================================

// Point 表示游戏面板上的一个坐标点
// 用于表示蛇的身体 segments 和食物的位置
type Point struct {
	x, y int
}

// ============================================
// Game 结构体 - 核心游戏状态
// ============================================
// 包含游戏的所有状态信息

type Game struct {
	// 游戏面板：0表示空，1表示被蛇身体占用
	board [][]int

	// 蛇身体：从头部(snake[0])到尾部排列的坐标列表
	snake []Point

	// 食物在面板上的位置
	food Point

	// direction: 蛇当前的实际移动方向
	// nextDir: 用户输入的下一个方向（用于防止快速反向导致自杀）
	direction Direction
	nextDir   Direction

	// 游戏状态
	score    int  // 当前得分（每吃一个食物+10分）
	length   int  // 蛇的目标长度（随得分增加）
	paused   bool // 游戏是否暂停
	gameOver bool // 游戏是否结束

	// 依赖组件
	rng *rand.Rand // 随机数生成器（用于生成食物位置）
}

// ============================================
// 工厂方法
// ============================================

// NewGame 创建并初始化一个新的贪吃蛇游戏
// screen: 用于渲染的 tcell 屏幕对象
func NewGame(screen interface{}) *Game {
	// 初始化游戏面板
	board := make([][]int, BoardHeight)
	for i := range board {
		board[i] = make([]int, BoardWidth)
	}

	// 初始化蛇的起始位置（面板中间，向上移动）
	snake := []Point{
		{BoardWidth / 2, BoardHeight / 2},
		{BoardWidth / 2, BoardHeight / 2 + 1},
		{BoardWidth / 2, BoardHeight / 2 + 2},
	}

	return &Game{
		board:     board,
		snake:     snake,
		food:      Point{},
		direction: Up,
		nextDir:   Up,
		score:     0,
		length:    3,
		paused:    false,
		gameOver:  false,
		rng:       rand.New(rand.NewSource(rand.Int63())),
	}
}

// ============================================
// 核心游戏逻辑
// ============================================

// spawnFood 在空白位置生成一个新的食物
// 算法：收集所有空白位置，随机选择一个作为食物位置
func (g *Game) spawnFood() {
	var emptyPoints []Point

	// 遍历整个面板，找出所有空白位置
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if g.board[y][x] == 0 {
				emptyPoints = append(emptyPoints, Point{x, y})
			}
		}
	}

	// 如果有空白位置，随机选择一个作为食物
	if len(emptyPoints) > 0 {
		g.food = emptyPoints[g.rng.Intn(len(emptyPoints))]
	}
}

// collides 检测给定坐标是否会发生碰撞
// head: 要检测的坐标点
// 返回值：如果发生碰撞返回 true，否则返回 false
//
// 碰撞检测包括：
// 1. 撞墙检测：坐标超出面板边界
// 2. 撞自身检测：坐标与蛇身体（除尾部外）重合
func (g *Game) collides(head Point) bool {
	// 撞墙检测
	if head.x < 0 || head.x >= BoardWidth || head.y < 0 || head.y >= BoardHeight {
		return true
	}

	// 撞自身检测（跳过蛇尾，因为蛇会移动）
	for i := 0; i < len(g.snake)-1; i++ {
		if g.snake[i] == head {
			return true
		}
	}

	return false
}

// move 让蛇移动一格
// 返回值：移动是否成功（失败时游戏结束）
//
// 移动逻辑：
// 1. 更新实际方向为用户输入的方向
// 2. 计算新的头部位置
// 3. 检测碰撞（撞墙或撞自身则游戏结束）
// 4. 检测是否吃到食物（头部与食物重合）
// 5. 添加新头部，根据是否吃到食物决定是否移除尾部
func (g *Game) move() bool {
	// 更新实际移动方向
	g.direction = g.nextDir

	// 计算新头部位置
	head := g.snake[0]
	switch g.direction {
	case Up:
		head.y--
	case Down:
		head.y++
	case Left:
		head.x--
	case Right:
		head.x++
	}

	// 碰撞检测
	if g.collides(head) {
		g.gameOver = true
		return false
	}

	// 检测是否吃到食物
	ateFood := (head == g.food)

	// 添加新头部到蛇身
	g.snake = append([]Point{head}, g.snake...)

	// 处理食物逻辑
	if ateFood {
		g.score += 10  // 增加得分
		g.length++     // 增加目标长度
		g.spawnFood()  // 生成新食物
	} else {
		// 未吃到食物，移除尾部以保持长度
		if len(g.snake) > g.length {
			g.snake = g.snake[:len(g.snake)-1]
		}
	}

	return true
}

// getSpeed 根据当前得分计算移动速度
// 返回值：移动间隔（毫秒），分数越高速度越快
//
// 速度计算公式：基础速度 - 得分/5
// 最小速度限制为 SpeedFast
func (g *Game) getSpeed() int {
	speed := SpeedNormal - g.score/5
	if speed < SpeedFast {
		return SpeedFast
	}
	return speed
}

// reset 重置游戏到初始状态
// 用于游戏结束后重新开始
func (g *Game) reset() {
	// 清空游戏面板
	for y := range g.board {
		for x := range g.board[y] {
			g.board[y][x] = 0
		}
	}

	// 重置蛇的位置
	g.snake = []Point{
		{BoardWidth / 2, BoardHeight / 2},
		{BoardWidth / 2, BoardHeight / 2 + 1},
		{BoardWidth / 2, BoardHeight / 2 + 2},
	}

	// 重置游戏状态
	g.direction = Up
	g.nextDir = Up
	g.score = 0
	g.length = 3
	g.paused = false
	g.gameOver = false

	// 生成新的食物
	g.spawnFood()
}
