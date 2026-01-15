package snake

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// ============================================
// Renderer - 游戏画面渲染器
// ============================================
// 负责将游戏状态绘制到终端屏幕

type Renderer struct {
	screen tcell.Screen // tcell 屏幕对象
	game   *Game        // 要渲染的游戏实例
}

// NewRenderer 创建渲染器实例
func NewRenderer(screen tcell.Screen, game *Game) *Renderer {
	return &Renderer{
		screen: screen,
		game:   game,
	}
}

// Render 绘制整个游戏画面
//
// 绘制顺序（从后到前）：
// 1. 清屏并设置背景色
// 2. 绘制游戏区域边框
// 3. 绘制蛇
// 4. 绘制食物
// 5. 绘制右侧信息面板
// 6. 绘制状态提示（暂停/游戏结束）
func (r *Renderer) Render() {
	// ---------- 1. 清屏 ----------
	r.screen.Clear()
	r.screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack))

	// ---------- 2. 绘制边框 ----------
	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	// 绘制左右边框
	for y := 0; y < BoardHeight+2; y++ {
		r.screen.SetContent(2, y+1, '|', nil, borderStyle)
		r.screen.SetContent(BoardWidth*2+3, y+1, '|', nil, borderStyle)
	}
	// 绘制上下边框
	for x := 0; x < BoardWidth*2+2; x++ {
		r.screen.SetContent(3+x, 1, '-', nil, borderStyle)
		r.screen.SetContent(3+x, BoardHeight+2, '-', nil, borderStyle)
	}

	// ---------- 3. 绘制蛇 ----------
	// 蛇头使用亮绿色，其他部分使用普通绿色
	for i, p := range r.game.snake {
		var snakeStyle tcell.Style
		if i == 0 {
			// 蛇头
			snakeStyle = tcell.StyleDefault.Foreground(tcell.ColorLime)
		} else {
			// 蛇身
			snakeStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen)
		}

		drawX := 4 + p.x*2
		drawY := p.y + 2
		r.screen.SetContent(drawX, drawY, '●', nil, snakeStyle)
		r.screen.SetContent(drawX+1, drawY, ' ', nil, snakeStyle)
	}

	// ---------- 4. 绘制食物 ----------
	foodStyle := tcell.StyleDefault.Foreground(tcell.ColorRed)
	drawX := 4 + r.game.food.x*2
	drawY := r.game.food.y + 2
	r.screen.SetContent(drawX, drawY, '★', nil, foodStyle)
	r.screen.SetContent(drawX+1, drawY, ' ', nil, foodStyle)

	// ---------- 5. 绘制右侧信息面板 ----------
	infoStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	nextX := BoardWidth*2 + 8

	// 游戏标题
	title := "SNAKE"
	for i, ch := range title {
		r.screen.SetContent(nextX+i, 2, ch, nil, infoStyle)
	}

	// 分数
	scoreText := fmt.Sprintf("SCORE: %d", r.game.score)
	for i, ch := range scoreText {
		r.screen.SetContent(nextX+i, 5, ch, nil, infoStyle)
	}

	// 操作说明
	controls := []string{
		"CONTROLS:",
		"↑↓←→ : Move",
		"P   : Pause",
		"R   : Restart",
		"Esc : Back to Menu",
	}
	for i, ctrl := range controls {
		for j, ch := range ctrl {
			r.screen.SetContent(nextX+j, 10+i, ch, nil, infoStyle)
		}
	}

	// ---------- 6. 绘制状态提示 ----------
	if r.game.paused {
		pauseText := "PAUSED"
		for i, ch := range pauseText {
			r.screen.SetContent(BoardWidth/2*2+i, BoardHeight/2+2, ch, nil, infoStyle)
		}
	}

	if r.game.gameOver {
		gameOverText := "GAME OVER"
		for i, ch := range gameOverText {
			r.screen.SetContent(BoardWidth/2*2+i, BoardHeight/2+2, ch, nil, infoStyle)
		}
		restartText := "Press R to restart"
		for i, ch := range restartText {
			r.screen.SetContent(BoardWidth/2*2-2+i, BoardHeight/2+4, ch, nil, infoStyle)
		}
	}

	// 刷新屏幕显示
	r.screen.Show()
}
