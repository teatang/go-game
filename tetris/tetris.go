package tetris

import (
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

// ============================================
// 游戏入口 - 主循环
// ============================================
// Run 启动并运行俄罗斯方块游戏
//
// 主循环逻辑：
// 1. 非阻塞方式检测用户输入事件
// 2. 根据时间间隔自动下落方块
// 3. 渲染游戏画面
//
// 输入处理：
// - ← →: 左右移动
// - ↑: 旋转
// - ↓: 软降（加速下落）
// - 空格: 硬降（直接落到底）
// - P: 暂停/继续
// - R: 游戏结束时重新开始
// - Esc: 返回主菜单
func Run(screen tcell.Screen) {
	game := NewGame()
	renderer := NewRenderer(screen, game)
	game.spawnPiece()
	renderer.Render()

	lastDrop := time.Now()

	for {
		// 计算下落间隔（毫秒）
		interval := game.getDropInterval()
		dropInterval := time.Duration(interval) * time.Millisecond

		// ---------- 处理用户输入 ----------
		if screen.HasPendingEvent() {
			event := screen.PollEvent()
			if event != nil {
				switch ev := event.(type) {
				case *tcell.EventKey:
					// 返回主菜单
					if ev.Key() == tcell.KeyEscape {
						return
					}

					// 退出游戏
					if ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' || ev.Rune() == 'Q' {
						os.Exit(0)
					}

					// 游戏结束时的操作
					if game.gameOver {
						if ev.Rune() == 'r' || ev.Rune() == 'R' {
							game.reset()
							renderer.Render()
						}
						continue
					}

					// 暂停/继续
					if ev.Rune() == 'p' || ev.Rune() == 'P' {
						game.paused = !game.paused
						renderer.Render()
						continue
					}
					if game.paused {
						continue
					}

					// 游戏控制
					switch ev.Key() {
					case tcell.KeyLeft:
						game.move(-1, 0)
					case tcell.KeyRight:
						game.move(1, 0)
					case tcell.KeyDown:
						game.drop()
					case tcell.KeyUp:
						game.rotate()
					case tcell.KeyRune:
						// 空格键：硬降（方块直接落到底）
						if ev.Rune() == ' ' {
							for game.drop() {
							}
						}
					}
					renderer.Render()

				case *tcell.EventResize:
					renderer.Render()
				}
			}
		}

		// ---------- 自动下落 ----------
		if !game.gameOver && !game.paused && time.Since(lastDrop) > dropInterval {
			game.drop()
			renderer.Render()
			lastDrop = time.Now()
		} else if !game.gameOver && !game.paused {
			// 避免CPU占用过高
			time.Sleep(10 * time.Millisecond)
		}
	}
}
