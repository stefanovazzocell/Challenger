package provider

import (
	"bytes"
	"crypto/sha512"
	"math/rand"

	"github.com/stefanovazzocell/Challenger/provider/utils"
	"golang.org/x/crypto/pbkdf2"
)

const (
	queryLength int = 64
	// Suggested min difficulty
	DifficultyMin int = 10
	// Suggested max difficulty
	DifficultyMax int = 1000000
	// Suggested min multiplier
	MultiplierMin int = 10
	// Suggested max multiplier
	MultiplierMax int = 1000000
)

var (
	salt []byte = []byte{29, 215, 51, 123, 17, 173, 109, 69, 105, 225, 104, 175, 142, 141, 150, 11, 47, 217, 158, 208, 209, 170, 85, 55, 34, 158, 139, 82, 119, 133, 224, 162, 73, 185, 90, 9, 176, 36, 45, 10, 123, 38, 125, 213, 88, 14, 1, 192, 86, 170, 176, 193, 44, 130, 127, 238, 157, 37, 210, 113, 133, 5, 49, 238}
)

func init() {
	err := utils.SeedRandom()
	if err != nil {
		panic(err)
	}
}

// The representation of a challenge to a user's browser
type Challenge struct {
	Id         string `json:"id"`
	Query      []byte `json:"q"`
	Challenge  []byte `json:"c"`
	Difficulty int    `json:"d"`
}

// A solution to a Challenge
type Solution struct {
	Id       string `json:"id"`
	Response []byte `json:"r"`
}

// A verifier to be stored by the server
type Verifier struct {
	Id       string
	Expected []byte
}

// Generate a new challenge and the matching verifier
//
// difficulty affects how much work both the server and the user have to do
// multiplier indicates the upper bound of how much more work the user has to do compared to the server
func NewChallenge(difficulty int, multiplier int) (Challenge, Verifier) {
	// Generate challenge and verifier
	challenge := Challenge{
		Id:         utils.RandomString(64),
		Query:      utils.RandomBytes(queryLength),
		Difficulty: difficulty,
	}
	var v uint32 = uint32(rand.Intn(multiplier))
	verifier := Verifier{
		Id: challenge.Id,
		Expected: []byte{
			byte(v),
			byte(v >> 8),
			byte(v >> 16),
			byte(v >> 24),
		},
	}
	// Compute challenge
	challenge.Challenge = pbkdf2.Key(append(verifier.Expected, challenge.Query...), salt, difficulty, 64, sha512.New)

	return challenge, verifier
}

// Solve challenge locally (for testing purposes)
func (challenge Challenge) solve() Solution {
	solution := Solution{
		Id:       challenge.Id,
		Response: []byte{0, 0, 0, 0},
	}
	for {
		// Verify guess
		if bytes.Equal(pbkdf2.Key(append(solution.Response, challenge.Query...), salt, challenge.Difficulty, 64, sha512.New), challenge.Challenge) {
			return solution
		}
		// Increase guess
		solution.Response[0] += 1
		if solution.Response[0] == 0 {
			solution.Response[1] += 1
			if solution.Response[1] == 0 {
				solution.Response[2] += 1
				if solution.Response[2] == 0 {
					solution.Response[3] += 1
				}
			}
		}
	}
}

// Returns true if the solution is valid
func (verifier Verifier) Verify(solution Solution) bool {
	return verifier.Id == solution.Id && bytes.Equal(verifier.Expected, solution.Response)
}
