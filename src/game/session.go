package game

import (
	"time"

	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/notes"
	"github.com/faiface/pixel/pixelgl"
)

type Session struct {
	Song       []notes.SongNote
	Played     []playedNote
	Score      float64
	Start      time.Time
	Duration   float64 // in seconds
	SongCursor int     // number of passsed notes in song
	updateMode func(dt float64, note notes.Pitch)
}

type playedNote struct {
	notes.SongNote
	Correct bool
}

func NewSession(song []notes.SongNote, mode string) *Session {
	s := &Session{
		Played: make([]playedNote, 0, 100),
		Song:   song,
		Start:  time.Now(),
	}
	if mode == "challenge" {
		s.updateMode = s.challengeUpdate
	} else {
		s.updateMode = s.trainingUpdate
	}
	return s
}

func (s *Session) currentNote() notes.Pitch {
	// skip played notes
	for s.SongCursor < len(s.Song) && s.Song[s.SongCursor].End() < s.Duration {
		s.SongCursor += 1
	}
	// no notes to play left
	if s.SongCursor >= len(s.Song) {
		return notes.Pause
	}
	// already should be playing some note
	if s.Song[s.SongCursor].Time < s.Duration {
		return s.Song[s.SongCursor].Pitch
	}
	// otherwise - not yet playing anything
	return notes.Pause
}

func (s *Session) Update(dt float64, note notes.Pitch) {
	s.updateMode(dt, note)
}

// Notes don't stop, and you need to hit correct ones in time
func (s *Session) challengeUpdate(dt float64, note notes.Pitch) {
	s.Duration = time.Since(s.Start).Seconds()
	playingCorrectly := note == s.currentNote()
	if note.Name != "p" {
		if playingCorrectly {
			s.Score += dt
		} else {
			s.Score -= dt
		}
		if len(s.Played) == 0 || s.Played[len(s.Played)-1].End() > 0 { // no note currently playing
			s.Played = append(s.Played, playedNote{
				SongNote: notes.SongNote{ // create new note
					Time:     s.Duration,
					Pitch:    note,
					Duration: -1.0,
				},
				Correct: playingCorrectly,
			})
		} else if s.Played[len(s.Played)-1].Pitch != note || s.Played[len(s.Played)-1].Correct != playingCorrectly { // note changed
			s.Played[len(s.Played)-1].Duration = s.Duration - s.Played[len(s.Played)-1].Time // end current one
			s.Played = append(s.Played, playedNote{
				SongNote: notes.SongNote{ // create new note
					Time:     s.Duration,
					Pitch:    note,
					Duration: -1.0,
				},
				Correct: playingCorrectly,
			})
		}
	} else { // no note
		if len(s.Played) > 0 && s.Played[len(s.Played)-1].End() < 0 { // there is a note still playing
			s.Played[len(s.Played)-1].Duration = s.Duration - s.Played[len(s.Played)-1].Time // end it
		}
	}

	if len(s.Played) > 2 && s.Played[0].End() < s.Duration-config.TimeLinePosition/config.NoteSPS { // note not visible
		s.Played = s.Played[1:] // remove
	}
}

// Notes move only while you play correct note, to progress - play all the notes
func (s *Session) trainingUpdate(dt float64, note notes.Pitch) {
	playingCorrectly := note == s.currentNote()
	if !playingCorrectly {
		return
	}
	s.Duration += dt
	if note.Name != "p" {
		if len(s.Played) == 0 || s.Played[len(s.Played)-1].End() > 0 { // no note currently playing
			s.Played = append(s.Played, playedNote{
				SongNote: notes.SongNote{ // create new note
					Time:     s.Duration,
					Pitch:    note,
					Duration: -1.0,
				},
				Correct: true,
			})
		} else if s.Played[len(s.Played)-1].Pitch != note { // note changed
			s.Played[len(s.Played)-1].Duration = s.Duration - s.Played[len(s.Played)-1].Time // end current one
			s.Played = append(s.Played, playedNote{
				SongNote: notes.SongNote{ // create new note
					Time:     s.Duration,
					Pitch:    note,
					Duration: -1.0,
				},
				Correct: true,
			})
		}
	} else { // no note
		if len(s.Played) > 0 && s.Played[len(s.Played)-1].End() < 0 { // there is a note still playing
			s.Played[len(s.Played)-1].Duration = s.Duration - s.Played[len(s.Played)-1].Time // end it
		}
	}

	if len(s.Played) > 2 && s.Played[0].End() < s.Duration-config.TimeLinePosition/config.NoteSPS { // note not visible
		s.Played = s.Played[1:] // remove
	}
}

func (s Session) Render(win *pixelgl.Window) {
	renderNoteLines(win)
	renderNotes(win, s.Song, s.Played, s.Duration)
	renderScore(win, int(s.Score*100))
}
