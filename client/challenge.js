"use strict";

const challengeSalt = new Uint8Array([29, 215, 51, 123, 17, 173, 109, 69, 105, 225, 104, 175, 142, 141, 150, 11, 47, 217, 158, 208, 209, 170, 85, 55, 34, 158, 139, 82, 119, 133, 224, 162, 73, 185, 90, 9, 176, 36, 45, 10, 123, 38, 125, 213, 88, 14, 1, 192, 86, 170, 176, 193, 44, 130, 127, 238, 157, 37, 210, 113, 133, 5, 49, 238])

// Decode a base64 Query to an Array
function queryToArray(query) {
    return Array.from(atob(query), (c) => c.charCodeAt(0))
}

async function SolveChallenge(challengeObj) {
    let guess = [0,0,0,0].concat(queryToArray(challengeObj["q"]))
    const challenge = queryToArray(challengeObj["c"])
    while (true) {
        // Try current guess
        let guess_0 = solveOneChallenge([...guess], challengeObj["d"], challenge)
        guess[0]++
        let guess_1 = solveOneChallenge([...guess], challengeObj["d"], challenge)
        guess[0]++
        let guess_2 = solveOneChallenge([...guess], challengeObj["d"], challenge)
        guess[0]++
        let guess_3 = solveOneChallenge([...guess], challengeObj["d"], challenge)
        guess[0]++
        // Look for a valid solution
        if (await guess_0) {
            guess[0]-=4
            return JSON.stringify({id:challengeObj["id"],r:guess.slice(0,4)})
        }
        if (await guess_1) {
            guess[0]-=3
            return JSON.stringify({id:challengeObj["id"],r:guess.slice(0,4)})
        }
        if (await guess_2) {
            guess[0]-=2
            return JSON.stringify({id:challengeObj["id"],r:guess.slice(0,4)})
        }
        if (await guess_3) {
            guess[0]-=1
            return JSON.stringify({id:challengeObj["id"],r:guess.slice(0,4)})
        }
        // Generate next guess
        if (guess[0]==256) {
            guess[0]=0
            guess[1]++
            if (guess[1]==256) {
                guess[1]=0
                guess[2]++
                if (guess[2]==256) {
                    guess[2]=0
                    guess[3]++
                }
            }
        }
    }
}

async function solveOneChallenge(key, iterations, challenge) {
    const value = new Uint8Array(await crypto.subtle.deriveBits({
        name:     "PBKDF2",
        hash:     "SHA-512",
        salt:     challengeSalt,
        iterations: iterations
    }, await crypto.subtle.importKey("raw", new Uint8Array(key), "PBKDF2", false, ["deriveBits"]), 512))
    for (let i = 0; i < challenge.length; i++) {
        if (value[i] != challenge[i]) return false
    }
    return true
}