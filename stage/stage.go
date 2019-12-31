package stage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"demo/object"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

/*
ゲームに必要なもの
各タイルの画像データ
各タイルのインデックス
各タイルのプロパティ一覧
*/

// Stage マップのデータ
type Stage struct {
	Width      int              // マップの横幅(タイル)
	Height     int              // マップの立幅(タイル)
	Image      *ebiten.Image    // マップ全体を画像データにしたもの
	TileIndex  []int            // len = Width*Height
	Properties map[int]Property // タイル番号 => プロパティ
	Actions    []*Action
	Objects    []*object.Object
	Warps      []*Warp
}

// Property タイルのプロパティ
type Property struct {
	Block  int // 通行可能か
	Action int // このタイルに対して何らかのアクションが可能か？
}

// Load マップを読み込む関数
func (stage *Stage) Load(stagename string) {
	filename := fmt.Sprintf("%s/stage.json", stagename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	raw := new(rawStage)
	if err := json.Unmarshal(file, raw); err != nil {
		panic(err)
	}

	stage.Properties = map[int]Property{}
	stage.Width = raw.Width
	stage.Height = raw.Height

	stage.Image, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("%s/stage.png", stagename), ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	stage.TileIndex = make([]int, stage.Height*stage.Width)
	copy(stage.TileIndex, raw.Layers[0].Data)

	// 各タイルセットについて
	for _, tileset := range raw.Tilesets {
		firstGID := tileset.FirstGID
		source := tileset.Source
		filename := fmt.Sprintf("%s/%s", stagename, source)
		stage.loadProperties(firstGID, filename)
	}

	stage.loadActions(fmt.Sprintf("%s/actions.json", stagename))
	stage.loadObjects(fmt.Sprintf("%s/objects.json", stagename))
	stage.loadWarps(fmt.Sprintf("%s/warp.json", stagename))
}

// GetProperty Get tile property
func (stage *Stage) GetProperty(x, y int) (target *Property) {
	target = &Property{Block: 1}

	if x >= 0 && x/16 < stage.Width && y >= 0 && y/16 < stage.Height {
		index := (y/16)*stage.Width + (x / 16)
		tileIndex := stage.TileIndex[index]
		property, ok := stage.Properties[tileIndex]
		if ok {
			target = &property
		} else {
			target = &Property{}
		}
		return target
	}

	if warp := stage.GetWarp(x, y); warp != nil {
		return &Property{}
	}

	return target
}

// GetObject Get Object
func (stage *Stage) GetObject(x, y int) (target *object.Object) {
	for _, object := range stage.Objects {
		switch object.Direction {
		case "up":
			if object.X/16 == (x+15)/16 && ((object.Y+16)/16-1) == y/16 {
				target = object
			}
		case "down":
			if object.X/16 == (x+15)/16 && (object.Y+15)/16 == y/16 {
				target = object
			}
		case "right":
			if (object.X+15)/16 == x/16 && object.Y/16 == (y+15)/16 {
				target = object
			}
		case "left":
			if ((object.X+16)/16-1) == x/16 && object.Y/16 == (y+15)/16 {
				target = object
			}
		}

		if target != nil {
			break
		}
	}
	return target
}

// GetAction Get Action
func (stage *Stage) GetAction(x, y int) (target *Action) {
	for _, action := range stage.Actions {
		if action.X == x/16 && action.Y == y/16 {
			target = action
			break
		}
	}
	return target
}

// GetWarp Get warp point
func (stage *Stage) GetWarp(x, y int) (target *Warp) {
	for _, warp := range stage.Warps {
		if warp.X*16 == x && warp.Y*16 == y {
			target = warp
			break
		}
	}
	return target
}

func (stage *Stage) loadProperties(firstGID int, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	tileset := new(TileSet)
	if err := json.Unmarshal(file, tileset); err != nil {
		panic(err)
	}

	// 各タイルのプロパティをセットしていく
	for _, tile := range tileset.List {
		tileID := tile.ID + firstGID

		newProperty := Property{}
		for _, property := range tile.Properties {
			switch property.Name {
			case "block":
				newProperty.Block = property.Value
			case "action":
				newProperty.Action = property.Value
			}
		}
		stage.Properties[tileID] = newProperty
	}
}

func (stage *Stage) loadActions(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	raw := new(Actions)
	if err := json.Unmarshal(file, raw); err != nil {
		panic(err)
	}
	stage.Actions = raw.List
}

func (stage *Stage) loadObjects(filename string) {
	stage.Objects = object.Load(filename)
}

func (stage *Stage) loadWarps(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	warps := new(Warps)
	if err := json.Unmarshal(file, warps); err != nil {
		panic(err)
	}
	stage.Warps = warps.List
}
