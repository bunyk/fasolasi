package game

import (
	"time"

	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/notes"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/constraints"
)

type Session struct {
	Song       []notes.SongNote
	Played     []playedNote
	Score      float64
	Start      time.Time
	Duration   float64 // session duration, progress of song in seconds
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

func (s Session) Finished() bool {
	return s.SongCursor >= len(s.Song)
}

func (s *Session) moveSongCursor() {
	// skip played notes
	for !s.Finished() && s.Song[s.SongCursor].End() < s.Duration {
		s.SongCursor += 1
	}
}

func (s *Session) nextNote() notes.SongNote {
	s.moveSongCursor()
	if s.Finished() {
		return notes.SongNote{}
	}
	return s.Song[s.SongCursor]
}

func (s *Session) currentNote() notes.Pitch {
	s.moveSongCursor()

	// no notes to play left
	if s.Finished() {
		return notes.Pause
	}
	nn := s.Song[s.SongCursor]

	// should be playing some note right now
	if nn.Time < s.Duration {
		return nn.Pitch
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

// Notes move only while you play correct note, to progress - play all the notes in correct orders.
// Obeying durations is optional.
func (s *Session) trainingUpdate(dt float64, note notes.Pitch) {
	if s.Finished() {
		s.Duration += dt
		s.Played = nil
		return
	}
	nn := s.nextNote()
	if note.Name != "p" {
		if len(s.Played) > 0 { // we were already playing some note
			if s.Played[0].Pitch == note { // still playing it
				s.Duration += dt
				if s.Duration > s.Played[0].End() { // Should have stopped already
					s.Score -= s.Duration - s.Played[0].End() // Decrease score
					s.Duration = s.Played[0].End()
					s.Played[0].Correct = false
				} else {
					s.Score += dt
				}
				return // and that's it for continuing playing note
			}
		}
		if note.Name == nn.Pitch.Name { // start playing current note
			s.Played = []playedNote{{
				SongNote: nn,
				Correct:  true,
			}}
			s.Score += 1.0
			s.Duration = nn.Time
			s.SongCursor += 1 // Prepare for next note
		}
	} else { // no note, probably stopped playing
		s.Played = nil
		s.Duration = nn.Time // Move timeline to next note
		if nn.Pitch.Name == "p" {
			s.Score += 1.0
			s.SongCursor += 1
		}
	}
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func (s Session) Render(win *pixelgl.Window) {
	renderNoteLines(win)
	renderNotes(win, s.Song, s.Played, s.Duration)
	renderScore(win, int(s.Score*100))
}
