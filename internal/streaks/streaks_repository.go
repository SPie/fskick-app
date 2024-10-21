package streaks

import (
	"github.com/spie/fskick/internal/db"
	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

type Streak struct {
    Number int
    Player players.Player
}

type PlayerWithAttendances struct {
    players.Player
    attendances []games.Attendance
}

type StreaksRepository struct {
    dbHandler db.Handler
}

func NewStreaksRepository(dbHandler db.Handler) StreaksRepository {
    return StreaksRepository{dbHandler: dbHandler}
}
