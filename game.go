package main

import (
	"image/color"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/font"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	sampleRate = beep.SampleRate(32000)
)

var (
	backgroundColor = color.NRGBA{0x00, 0x00, 0x00, 0xff}
	notesYPos       = map[string]int{
		"U": 360,
		"L": 372,
		"R": 378,
		"D": 385,
		"A": 400,
	}
)

type Game struct {
	ScreenWidth  int
	ScreenHeight int
	SongNotes    string
	PlayingSong  bool
	PlayingNote  bool
	CheckedSong  bool

	TextBox         *ebiten.Image
	BackgroundImage *ebiten.Image
	LinkImage       *ebiten.Image
	StaffLinesImage *ebiten.Image
	NoteImages      map[string]*ebiten.Image
	TextFont        font.Face
}

type ValidSong struct {
	Filename string
	SongName string
}

func NewGame(textFont font.Face) *Game {
	game := &Game{}
	game.ScreenWidth = 800
	game.ScreenHeight = 450
	game.SongNotes = ""
	game.PlayingSong = false
	game.PlayingNote = false
	game.CheckedSong = false

	game.TextFont = textFont
	game.TextBox = ebiten.NewImage(275, 25)
	game.ResetTextBox()
	game.BackgroundImage = loadImage("images/bg.png")
	game.LinkImage = loadImage("images/link.png")
	game.StaffLinesImage = loadImage("images/stafflines.png")
	game.NoteImages = map[string]*ebiten.Image{
		"U": loadImage("images/C-Up.png"),
		"D": loadImage("images/C-Down.png"),
		"L": loadImage("images/C-Left.png"),
		"R": loadImage("images/C-Right.png"),
		"A": loadImage("images/Button-A.png"),
	}

	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	return game
}

func loadImage(imageFilename string) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(imageFilename)
	if err != nil {
		log.Fatal(err)
	}
	return image
}

func (game *Game) ResetTextBox() {
	game.TextBox.Fill(color.NRGBA{0x00, 0x00, 0x00, 0x00})
	text.Draw(game.TextBox, "Play a song", game.TextFont, 10, 20, color.White)
}

func (game *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		game.PlaySingleNote("U")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		game.PlaySingleNote("D")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		game.PlaySingleNote("L")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		game.PlaySingleNote("R")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		game.PlaySingleNote("A")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		if !game.PlayingSong {
			game.SongNotes = ""
		}
	}

	if !game.CheckedSong {
		validSong := game.MatchSong(game.SongNotes)
		if validSong != nil {
			game.PlayingSong = true
			go func() {
				game.PlayValidSong(validSong.Filename, validSong.SongName)
				game.SongNotes = ""
				game.PlayingSong = false
				game.ResetTextBox()
			}()
		} else if len(game.SongNotes) >= 8 {
			go func() {
				playAudio("audio/Song_Error.wav")
				game.SongNotes = ""
			}()
		}
		game.CheckedSong = true
	}

	return nil
}

func (game *Game) MatchSong(notes string) *ValidSong {
	switch notes {
	case "LURLUR":
		return &ValidSong{
			Filename: "audio/songs/Zeldas_Lullaby.mp3",
			SongName: "Zelda's Lullaby",
		}
	case "ULRULR":
		return &ValidSong{
			Filename: "audio/songs/Eponas_Song.mp3",
			SongName: "Epona's Song",
		}
	case "DRLDRL":
		return &ValidSong{
			Filename: "audio/songs/Sarias_Song.mp3",
			SongName: "Saria's Song",
		}
	case "RDURDU":
		return &ValidSong{
			Filename: "audio/songs/Suns_Song.mp3",
			SongName: "Sun's Song",
		}
	case "RADRAD":
		return &ValidSong{
			Filename: "audio/songs/Song_of_Time.mp3",
			SongName: "Song of Time",
		}
	case "ADUADU":
		return &ValidSong{
			Filename: "audio/songs/Song_of_Storms.mp3",
			SongName: "Song of Storms",
		}
	case "AULRLR":
		return &ValidSong{
			Filename: "audio/songs/Minuet_Of_Woods.mp3",
			SongName: "Minuet of Forest",
		}
	case "DADARDRD":
		return &ValidSong{
			Filename: "audio/songs/Bolero_of_Fire.mp3",
			SongName: "Bolero of Fire",
		}
	case "ADRRL":
		return &ValidSong{
			Filename: "audio/songs/Serenade_of_Water.mp3",
			SongName: "Serenade of Water",
		}
	case "ADARDA":
		return &ValidSong{
			Filename: "audio/songs/Requiem_of_Spirit.mp3",
			SongName: "Requiem of Spirit",
		}
	case "LRRALRD":
		return &ValidSong{
			Filename: "audio/songs/Nocturne_of_Shadow.mp3",
			SongName: "Nocturne of Shadow",
		}
	case "URURLU":
		return &ValidSong{
			Filename: "audio/songs/Prelude_of_Light.mp3",
			SongName: "Prelude of Light",
		}
	case "LRDLRD":
		return &ValidSong{
			Filename: "audio/songs/Song_of_Healing.mp3",
			SongName: "Song of Healing",
		}
	case "DARDAR":
		return &ValidSong{
			Filename: "audio/songs/Inverted_Song_of_Time.mp3",
			SongName: "Inverted Song of Time",
		}
	case "RRAADD":
		return &ValidSong{
			Filename: "audio/songs/Song_of_Double_Time.mp3",
			SongName: "Song of Double Time",
		}
	case "DLUDLU":
		return &ValidSong{
			Filename: "audio/songs/Song_of_Soaring.mp3",
			SongName: "Song of Soaring",
		}
	case "ULULARA":
		return &ValidSong{
			Filename: "audio/songs/Sonata_of_Awakening.mp3",
			SongName: "Sonata of Awakening",
		}
	case "ARLARLRA":
		return &ValidSong{
			Filename: "audio/songs/Goron_Lullaby.mp3",
			SongName: "Goron Lullaby",
		}
	case "LULRDLR":
		return &ValidSong{
			Filename: "audio/songs/New_Wave_Bossa_Nova.mp3",
			SongName: "New Wave Bossa Nova",
		}
	case "RLRDRUL":
		return &ValidSong{
			Filename: "audio/songs/Elegy_of_Emptiness.mp3",
			SongName: "Elegy of Emptiness",
		}
	case "RDADRU":
		return &ValidSong{
			Filename: "audio/songs/Oath_to_Order.mp3",
			SongName: "Oath to Order",
		}
	}
	return nil
}

func (game *Game) PlayValidSong(songFilename string, songName string) {
	game.TextBox.Fill(color.NRGBA{0x00, 0x00, 0x00, 0x00})

	songNameColor := color.NRGBA{124, 162, 195, 0xff} // TODO: adjust color per song name
	text.Draw(game.TextBox, "You played", game.TextFont, 10, 20, color.White)
	text.Draw(game.TextBox, songName, game.TextFont, 100, 20, songNameColor)

	playAudio("audio/Song_Correct.wav")
	playAudio(songFilename)
}

func (game *Game) PlaySingleNote(note string) {
	if game.PlayingSong {
		return
	}

	game.PlayingNote = true
	game.CheckedSong = false
	game.SongNotes += note
	go func() {
		switch note {
		case "U":
			playAudio("audio/ocarina/D2.wav")
		case "L":
			playAudio("audio/ocarina/B.wav")
		case "R":
			playAudio("audio/ocarina/A.wav")
		case "D":
			playAudio("audio/ocarina/F.wav")
		case "A":
			playAudio("audio/ocarina/D.wav")
		}
		game.PlayingNote = false
	}()
}

func playAudio(audioFilename string) {
	f, err := os.Open(audioFilename)
	if err != nil {
		log.Fatal(err)
	}

	var streamer beep.StreamSeekCloser
	var format beep.Format

	fileType := filepath.Ext(audioFilename)
	if fileType == ".mp3" {
		streamer, format, err = mp3.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()
	} else if fileType == ".wav" {
		streamer, format, err = wav.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()
	} else {
		log.Fatal("Audio file must be MP3 or WAV")
	}

	done := make(chan bool)
	resampled := beep.Resample(4, format.SampleRate, sampleRate, streamer)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func (game *Game) DrawNotes(screen *ebiten.Image, notes string) {
	for i := 0; i < len(notes); i++ {
		note := string(notes[i])

		optionsNoteImage := &ebiten.DrawImageOptions{}
		optionsNoteImage.GeoM.Translate(float64(275+(25*i)), float64(notesYPos[note]))
		screen.DrawImage(game.NoteImages[note], optionsNoteImage)
	}
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	screen.DrawImage(game.BackgroundImage, nil)

	optionsLinkImage := &ebiten.DrawImageOptions{}
	optionsLinkImage.GeoM.Translate(325, 200)
	screen.DrawImage(game.LinkImage, optionsLinkImage)

	optionsStaffLinesImage := &ebiten.DrawImageOptions{}
	optionsStaffLinesImage.GeoM.Scale(1.5, 1)
	optionsStaffLinesImage.GeoM.Translate(110, 325)
	screen.DrawImage(game.StaffLinesImage, optionsStaffLinesImage)

	optionsTextBox := &ebiten.DrawImageOptions{}
	optionsTextBox.GeoM.Translate(250, 325)
	screen.DrawImage(game.TextBox, optionsTextBox)

	game.DrawNotes(screen, game.SongNotes)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.ScreenWidth, game.ScreenHeight
}
