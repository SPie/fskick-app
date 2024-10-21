package streaks

import (
	"fmt"
	"slices"

	"github.com/spie/fskick/internal/games"
)

type Manager struct {
    attendanceRepository games.AttendanceRepository
}

func NewManager(attendanceRepository games.AttendanceRepository) Manager {
    return Manager{attendanceRepository: attendanceRepository}
}

func (manager Manager) GetLongestWiningStreak() (Streak, error) {
    allPlayersWithAttendances, err := manager.attendanceRepository.GetAttendancesForAllPlayers()
    if err != nil {
        return Streak{}, fmt.Errorf("get longest winning streak: %w", err)
    }

    longestWinningStreak := Streak{}
    for _, playerWithAttendance := range allPlayersWithAttendances {
        streak := Streak{Player: playerWithAttendance.Player, Number: 0}

        number := 0
        for _, attendance := range playerWithAttendance.Attendances {
            if !attendance.Win {
                if number > streak.Number {
                    streak.Number = number
                }

                number = 0
                continue
            }

            number++
        }

        if number > streak.Number {
            streak.Number = number
        }

        if streak.Number > longestWinningStreak.Number {
            longestWinningStreak = streak
        }
    }

    return longestWinningStreak, nil
}

func (manager Manager) GetLongestLosingStreak() (Streak, error) {
    allPlayersWithAttendances, err := manager.attendanceRepository.GetAttendancesForAllPlayers()
    if err != nil {
        return Streak{}, fmt.Errorf("get longest losing streak: %w", err)
    }

    longestLosingStreak := Streak{}
    for _, playerWithAttendance := range allPlayersWithAttendances {
        streak := Streak{Player: playerWithAttendance.Player, Number: 0}

        number := 0
        for _, attendance := range playerWithAttendance.Attendances {
            if attendance.Win {
                if number > streak.Number {
                    streak.Number = number
                }

                number = 0
                continue
            }

            number++
        }

        if number > streak.Number {
            streak.Number = number
        }

        if streak.Number > longestLosingStreak.Number {
            longestLosingStreak = streak
        }
    }

    return longestLosingStreak, nil
}

func (manager Manager) GetCurrentStreaks(win bool) ([]Streak, error) {
    allPlayersWithAttendances, err := manager.attendanceRepository.GetAttendancesForAllPlayers()
    if err != nil {
        return nil, fmt.Errorf("get current streaks: %w", err)
    }

    currentStreaks := make([]Streak, len(allPlayersWithAttendances))
    for i, playerWithAttendance := range allPlayersWithAttendances {
        streak := Streak{Player: playerWithAttendance.Player, Number: 0}
        slices.Reverse(playerWithAttendance.Attendances)
        for _, attendance := range playerWithAttendance.Attendances {
            if attendance.Win != win {
                break
            }

            streak.Number++
        }

        currentStreaks[i] = streak
    }

    slices.SortFunc(currentStreaks, func(a, b Streak) int {
        if a.Number < b.Number {
            return 1
        }
        if a.Number > b.Number {
            return -1
        }

        if a.Player.ID < b.Player.ID {
            return 1
        }

        return -1
    })

    return currentStreaks, nil
}
