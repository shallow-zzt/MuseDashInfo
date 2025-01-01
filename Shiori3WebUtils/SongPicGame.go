package MDWebPageUtils

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"

	"image/png"
	"math"
	"os"
)

type MiniGameContent struct {
	Rotateb64      string
	Answerb64      string
	Cuttedb64      string
	SongAnswerName string
}

type MiniGameResult struct {
	Result         int      `json:"result"`
	ResultNum      int      `json:"result-num"`
	PossibleResult []string `json:"possible-result"`
}

func rotateImage(src image.Image, angle float64) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	centerX, centerY := width/2, height/2

	rotated := image.NewRGBA(bounds)
	draw.Draw(rotated, rotated.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)

	angleRad := angle * math.Pi / 180
	cosAngle := math.Cos(angleRad)
	sinAngle := math.Sin(angleRad)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := x - centerX
			dy := y - centerY

			srcX := int(float64(dx)*cosAngle - float64(dy)*sinAngle + float64(centerX))
			srcY := int(float64(dx)*sinAngle + float64(dy)*cosAngle + float64(centerY))

			if srcX >= 0 && srcX < width && srcY >= 0 && srcY < height {
				rotated.Set(x, y, src.At(srcX, srcY))
			} else {
				rotated.Set(x, y, color.Transparent)
			}
		}
	}

	return rotated
}

func highlightImage(src image.Image, x0 int, y0 int, width int, height int, borderWidth int) image.Image {
	bounds := src.Bounds()
	highlighted := image.NewRGBA(bounds)
	draw.Draw(highlighted, highlighted.Bounds(), src, image.Point{}, draw.Src)

	red := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	for x := x0 - 1; x <= x0+width; x++ {
		for i := 0; i <= borderWidth; i++ {
			highlighted.Set(x, y0-1+i, red)
			highlighted.Set(x, y0+height-i, red)
		}

	}
	for y := y0 - 1; y <= y0+height; y++ {
		for i := 0; i <= borderWidth; i++ {
			highlighted.Set(x0-1+i, y, red)
			highlighted.Set(x0+width-i, y, red)
		}
	}

	return highlighted
}

func productingSongPic(songPicName string, isRotate bool) (image.Image, image.Image, image.Image) {

	var angle float64
	inputFile := songPicName + ".png"
	if isRotate {
		angle = rand.Float64() * 360.0
	} else {
		angle = 0.0
	}

	file, err := os.Open(inputFile)
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("解码 PNG 失败:", err)
		return nil, nil, nil
	}
	rotatedImg := rotateImage(img, angle)

	width := 80
	height := 80
	x0 := rand.Intn(312-width) + 64
	y0 := rand.Intn(312-height) + 64

	highlightImg := highlightImage(rotatedImg, x0, y0, width, height, 2)

	subImg := rotatedImg.(interface {
		SubImage(rect image.Rectangle) image.Image
	}).SubImage(image.Rect(x0, y0, x0+width, y0+height))

	return rotatedImg, highlightImg, subImg

}

func GenerateMiniGameContent(db *sql.DB, isRotate bool) MiniGameContent {

	songPicName, songAnswerName := GetRandomSongPicName(db)
	songPicPath := "SongAssets/Songpic/" + songPicName

	rotatePic, answerPic, cuttedPic := productingSongPic(songPicPath, isRotate)

	rotateBuf := new(bytes.Buffer)
	answerBuf := new(bytes.Buffer)
	cuttedBuf := new(bytes.Buffer)

	err := png.Encode(rotateBuf, rotatePic)
	err = png.Encode(answerBuf, answerPic)
	err = png.Encode(cuttedBuf, cuttedPic)
	if err != nil {
		fmt.Println(err)
	}

	rotateb64 := base64.StdEncoding.EncodeToString(rotateBuf.Bytes())
	answerb64 := base64.StdEncoding.EncodeToString(answerBuf.Bytes())
	cuttedb64 := base64.StdEncoding.EncodeToString(cuttedBuf.Bytes())

	return MiniGameContent{Rotateb64: rotateb64, Answerb64: answerb64, Cuttedb64: cuttedb64, SongAnswerName: songAnswerName}
}

func GetRandomSongPicName(db *sql.DB) (string, string) {
	var songPicName string
	var songName string
	getRandomSongPic := `SELECT music_pic_name,music_name
	FROM mdsong
	ORDER BY RANDOM()
	LIMIT 1`
	result := db.QueryRow(getRandomSongPic)
	result.Scan(&songPicName, &songName)
	return songPicName, songName
}

func GetPossibleAnswer(db *sql.DB, answer string) []string {
	var possibleAnswers []string
	getPossibleAnswer := `SELECT DISTINCT music_name
FROM(SELECT ms.music_album,ms.music_album_number,music_name,music_alias
FROM mdsong ms
LEFT OUTER JOIN mdalias ma
ON ms.music_album = ma.music_album AND ms.music_album_number = ma.music_album_number)
WHERE music_name LIKE "%` + answer + `%" OR music_alias LIKE "` + answer + `"`
	result, err := db.Query(getPossibleAnswer)

	for result.Next() {
		var name string
		err = result.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		possibleAnswers = append(possibleAnswers, name)
	}

	return possibleAnswers
}
