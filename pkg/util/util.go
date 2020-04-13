package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	cursorImage, _, _ = ebitenutil.NewImageFromFile("asset/menu/cursor.png", ebiten.FilterDefault)
)

// Contains - 指定した要素を含んでいるか
func Contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func PaddingID(ID int) string {
	switch {
	case ID < 10:
		return fmt.Sprintf("00%d", ID)
	case ID < 100:
		return fmt.Sprintf("0%d", ID)
	default:
		return fmt.Sprintf("%d", ID)
	}
}

// PaddingLeft - 文字列を指定した文字数になるまで空白で前方にパディングする
func PaddingLeft(str string, to uint, char string) string {
	from := len(str)
	for i := uint(from); i < to; i++ {
		str = char + str
	}
	return str
}

// PaddingRight - 文字列を指定した文字数になるまで空白で後方にパディングする
func PaddingRight(str string, to uint) string {
	from := len(str)
	for i := uint(from); i < to; i++ {
		str += " "
	}
	return str
}

// Chance - 確率を与えるとその確率にしたがってtrueかfalseを返す
func Chance(probability float64) bool {
	rand.Seed(time.Now().UnixNano())
	p := rand.Intn(100)
	return float64(p) <= probability
}

// IsSwitchCommand - コマンドが交代コマンドかどうか
func IsSwitchCommand(command uint) bool {
	return command >= 5 && command <= 9
}

// SetCursor - カーソルをtargetに描画する
func SetCursor(target *ebiten.Image, cursor, y int) *ebiten.Image {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(8), float64(16*cursor+y))
	target.DrawImage(cursorImage, op)
	return target
}
