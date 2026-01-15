package tetris

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
// 3. 绘制已锁定的方块
// 4. 绘制幽灵方块（预览最终位置）
// 5. 绘制当前下落的方块
// 6. 绘制右侧信息面板
// 7. 绘制状态提示（暂停/游戏结束）
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

	// ---------- 3. 绘制已锁定的方块 ----------
	// 这些是之前落下方块并已锁定的
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if r.game.board[y][x] != 0 {
				color := Colors[r.game.board[y][x]-1]
				cellStyle := tcell.StyleDefault.Foreground(getColor(color))
				r.screen.SetContent(4+x*2, y+2, '■', nil, cellStyle)
				r.screen.SetContent(5+x*2, y+2, ' ', nil, cellStyle)
			}
		}
	}

	// ---------- 4. 绘制幽灵方块 ----------
	// 预览当前方块最终会落到的位置（灰色半透明效果）
	if !r.game.gameOver {
		ghostX, ghostY := r.game.getGhostPosition()
		ghostStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
		for y, row := range r.game.currShape {
			for x, cell := range row {
				if cell == 1 {
					drawX := 4 + (ghostX+x)*2
					drawY := ghostY + y + 2
					if drawY >= 2 && drawY < BoardHeight+2 {
						r.screen.SetContent(drawX, drawY, '░', nil, ghostStyle)
						r.screen.SetContent(drawX+1, drawY, ' ', nil, ghostStyle)
					}
				}
			}
		}
	}

	// ---------- 5. 绘制当前下落的方块 ----------
	color := Colors[r.game.currPiece]
	cellStyle := tcell.StyleDefault.Foreground(getColor(color))
	for y, row := range r.game.currShape {
		for x, cell := range row {
			if cell == 1 {
				drawX := 4 + (r.game.pieceX+x)*2
				drawY := r.game.pieceY + y + 2
				if drawY >= 2 && drawY < BoardHeight+2 {
					r.screen.SetContent(drawX, drawY, '■', nil, cellStyle)
					r.screen.SetContent(drawX+1, drawY, ' ', nil, cellStyle)
				}
			}
		}
	}

	// ---------- 6. 绘制右侧信息面板 ----------
	infoStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	nextX := BoardWidth*2 + 8

	// "NEXT" 标签
	r.screen.SetContent(nextX, 2, 'N', nil, infoStyle)
	r.screen.SetContent(nextX+1, 2, 'E', nil, infoStyle)
	r.screen.SetContent(nextX+2, 2, 'X', nil, infoStyle)
	r.screen.SetContent(nextX+3, 2, 'T', nil, infoStyle)

	// 下一个方块预览
	if r.game.nextPiece > 0 {
		nextPieceIdx := r.game.nextPiece - 1
		for y, row := range Shapes[nextPieceIdx] {
			for x, cell := range row {
				if cell == 1 {
					color := Colors[nextPieceIdx]
					cellStyle := tcell.StyleDefault.Foreground(getColor(color))
					r.screen.SetContent(nextX+x*2, y+4, '■', nil, cellStyle)
					r.screen.SetContent(nextX+x*2+1, y+4, ' ', nil, cellStyle)
				}
			}
		}
	}

	// 分数信息
	scoreText := fmt.Sprintf("SCORE: %d", r.game.score)
	for i, ch := range scoreText {
		r.screen.SetContent(nextX+i, 8, ch, nil, infoStyle)
	}
	linesText := fmt.Sprintf("LINES: %d", r.game.lines)
	for i, ch := range linesText {
		r.screen.SetContent(nextX+i, 10, ch, nil, infoStyle)
	}
	levelText := fmt.Sprintf("LEVEL: %d", r.game.level)
	for i, ch := range levelText {
		r.screen.SetContent(nextX+i, 12, ch, nil, infoStyle)
	}

	// 操作说明
	controls := []string{
		"CONTROLS:",
		"←→ : Move",
		"↑   : Rotate",
		"↓   : Soft Drop",
		"Space: Hard Drop",
		"P   : Pause",
		"Esc : Menu",
	}
	for i, ctrl := range controls {
		for j, ch := range ctrl {
			r.screen.SetContent(nextX+j, 16+i, ch, nil, infoStyle)
		}
	}

	// ---------- 7. 绘制状态提示 ----------
	if r.game.paused {
		for i, ch := range "PAUSED" {
			r.screen.SetContent(BoardWidth+2+i, BoardHeight/2+2, ch, nil, infoStyle)
		}
	}
	if r.game.gameOver {
		for i, ch := range "GAME OVER" {
			r.screen.SetContent(BoardWidth+1+i, BoardHeight/2+2, ch, nil, infoStyle)
		}
		for i, ch := range "Press R to restart" {
			r.screen.SetContent(BoardWidth-1+i, BoardHeight/2+4, ch, nil, infoStyle)
		}
	}

	// 刷新屏幕显示
	r.screen.Show()
}

// getColor 辅助函数：根据颜色名称返回 tcell.Color
func getColor(name string) tcell.Color {
	switch name {
	case "cyan":
		return tcell.ColorAqua
	case "yellow":
		return tcell.ColorYellow
	case "fuchsia":
		return tcell.ColorFuchsia
	case "lime":
		return tcell.ColorLime
	case "red":
		return tcell.ColorRed
	case "navy":
		return tcell.ColorNavy
	case "olive":
		return tcell.ColorOlive
	default:
		return tcell.ColorWhite
	}
}
