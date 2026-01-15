package snake

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// ============================================
// 游戏入口 - 主循环
// ============================================
// Run 启动并运行贪吃蛇游戏
//
// 主循环逻辑：
// 1. 非阻塞方式检测用户输入事件
// 2. 根据时间间隔自动移动蛇
// 3. 渲染游戏画面
//
// 输入处理：
// - 方向键：控制蛇的移动方向（防止快速反向）
// - P 键：暂停/继续游戏
// - R 键：游戏结束时重新开始
// - Esc 键：返回主菜单
func Run(screen tcell.Screen) {
	game := NewGame(screen)
	renderer := NewRenderer(screen, game)
	game.spawnFood()
	renderer.Render()

	lastMove := time.Now()

	for {
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

					// 方向控制（防止快速反向导致自杀）
					switch ev.Key() {
					case tcell.KeyUp:
						if game.direction != Down {
							game.nextDir = Up
						}
					case tcell.KeyDown:
						if game.direction != Up {
							game.nextDir = Down
						}
					case tcell.KeyLeft:
						if game.direction != Right {
							game.nextDir = Left
						}
					case tcell.KeyRight:
						if game.direction != Left {
							game.nextDir = Right
						}
					case tcell.KeyRune:
						if ev.Rune() == 'r' || ev.Rune() == 'R' {
							game.reset()
						}
					}
					renderer.Render()

				case *tcell.EventResize:
					renderer.Render()
				}
			}
		}

		// ---------- 自动移动 ----------
		if !game.gameOver && !game.paused && time.Since(lastMove) > time.Duration(game.getSpeed())*time.Millisecond {
			game.move()
			renderer.Render()
			lastMove = time.Now()
		} else if !game.gameOver && !game.paused {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
