package streaks

import (
	"fmt"
	"slices"

	"github.com/spie/fskick/internal/games"
	"github.com/spie/fskick/internal/players"
)

type Manager struct {
    attendanceRepository games.AttendanceRepository
}

func NewManager(attendanceRepository games.AttendanceRepository) Manager {
    return Manager{attendanceRepository: attendanceRepository}
}

func (manager Manager) GetLongestWinningAndLosingStreaks() (
    winningStreak Streak,
    logingStreak Streak,
    err error,
) {
    allPlayersWithAttendances, err := manager.attendanceRepository.GetAttendancesForAllPlayers()
    if err != nil {
        return Streak{}, Streak{}, fmt.Errorf("get longest streak: %w", err)
    }

    longestWinningStreak := Streak{}
    longestLosingStreak := Streak{}
    for _, playerWithAttendance := range allPlayersWithAttendances {
        winningStreak := GetLongestStreakForPlayer(
            playerWithAttendance.Player,
            playerWithAttendance.Attendances,
            true,
        )

        if winningStreak.Number > longestWinningStreak.Number {
            longestWinningStreak = winningStreak
        }

        losingStreak := GetLongestStreakForPlayer(
            playerWithAttendance.Player,
            playerWithAttendance.Attendances,
            false,
        )

        if losingStreak.Number > longestLosingStreak.Number {
            longestLosingStreak = losingStreak
        }
    }

    return longestWinningStreak, longestLosingStreak, nil
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

func GetLongestStreakForPlayer(player players.Player, attendances []games.Attendance, win bool) Streak {
    streak := Streak{Player: player, Number: 0}

    number := 0
    for _, attendance := range attendances {
        if attendance.Win !=  win {
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

    return streak
}
