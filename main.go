package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	snakepkg "go-game/snake"
	tetrispkg "go-game/tetris"
)

// ============================================
// 游戏类型
// ============================================

type GameType int

const (
	GameTetris GameType = iota
	GameSnake
)

// ============================================
// Menu - 主菜单
// ============================================

type Menu struct {
	screen   tcell.Screen
	selected int
	options  []string
}

// NewMenu 创建新菜单
func NewMenu(screen tcell.Screen) *Menu {
	return &Menu{
		screen:   screen,
		selected: 0,
		options: []string{
			"► 俄罗斯方块",
			"○ 贪吃蛇",
			"  退出游戏",
		},
	}
}

// Render 绘制菜单
func (m *Menu) Render() {
	m.screen.Clear()
	m.screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack))

	// 标题
	titleStyle := tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true)
	title := "TERMINAL GAMES"
	for i, ch := range title {
		m.screen.SetContent(10+i, 3, ch, nil, titleStyle)
	}

	// 副标题
	subtitleStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)
	subtitle := "Select a game to play"
	for i, ch := range subtitle {
		m.screen.SetContent(7+i, 5, ch, nil, subtitleStyle)
	}

	// 菜单选项
	for i, option := range m.options {
		var style tcell.Style
		if i == m.selected {
			style = tcell.StyleDefault.Foreground(tcell.ColorLime).Bold(true)
		} else {
			style = tcell.StyleDefault.Foreground(tcell.ColorWhite)
		}
		for j, ch := range option {
			m.screen.SetContent(8+j, 10+i*2, ch, nil, style)
		}
	}

	// 操作提示
	hintStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
	hints := []string{
		"↑↓ : Select",
		"Enter : Confirm",
		"Q : Quit",
	}
	for i, hint := range hints {
		for j, ch := range hint {
			m.screen.SetContent(8+j, 20+i, ch, nil, hintStyle)
		}
	}

	m.screen.Show()
}

// Run 运行菜单，返回选择的游戏类型
func (m *Menu) Run() GameType {
	m.Render()

	for {
		if m.screen.HasPendingEvent() {
			event := m.screen.PollEvent()
			if event != nil {
				switch ev := event.(type) {
				case *tcell.EventKey:
					if ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' || ev.Rune() == 'Q' {
						os.Exit(0)
					}

					switch ev.Key() {
					case tcell.KeyUp:
						if m.selected > 0 {
							m.selected--
							m.Render()
						}
					case tcell.KeyDown:
						if m.selected < len(m.options)-1 {
							m.selected++
							m.Render()
						}
					case tcell.KeyEnter:
						switch m.selected {
						case 0:
							return GameTetris
						case 1:
							return GameSnake
						case 2:
							os.Exit(0)
						}
					}
				case *tcell.EventResize:
					m.Render()
				}
			}
		}
	}
}

// ============================================
// 主程序入口
// ============================================

func main() {
	// 初始化屏幕
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create screen: %v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	screen.EnablePaste()
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack))

	// 主循环
	for {
		menu := NewMenu(screen)
		gameType := menu.Run()

		switch gameType {
		case GameTetris:
			tetrispkg.Run(screen)
		case GameSnake:
			snakepkg.Run(screen)
		}
	}
}
