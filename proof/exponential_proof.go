// Package proof exponentiation proof
// Protocol ZKPoKE for R_{PoKE}
// Paper: Batching Techniques for Accumulators with Applications to IOPs and Stateless Blockchains
// Link: https://eprint.iacr.org/2018/1188.pdf
package proof

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"math/big"

	"github.com/rsa_accumulator/accumulator"
)

const (
	epChallengeStatement = "u^x = w, x is an integer"
)

// ExponentiationProof represents the proof for proof of exponentiation
type ExponentiationProof struct {
	witness  *big.Int // witness x in Z, u^x = w
	commit   *epCommitment
	response *epResponse
}

type epChallenge struct {
	statement  string        // the statement for the challenge
	g, h, n    *big.Int      // public parameters: G, H, N
	commitment *epCommitment // the commitment of by the prover
}

// newEPChallenge generates a new challenge for proof of exponentiation
func newEPChallenge(pp *PublicParameters, commitment *epCommitment) *epChallenge {
	return &epChallenge{
		statement:  epChallengeStatement,
		g:          pp.G,
		h:          pp.H,
		n:          pp.N,
		commitment: commitment,
	}
}

// Serialize generates the serialized data for proof of exponentiation challenge in byte format
func (e *epChallenge) serialize() []byte {
	var buf bytes.Buffer
	defer buf.Reset()
	buf.WriteString(e.statement)
	buf.WriteString(e.g.String())
	buf.WriteString(e.h.String())
	buf.WriteString(e.n.String())
	buf.WriteString(e.commitment.z.String())
	buf.WriteString(e.commitment.aG.String())
	buf.WriteString(e.commitment.aU.String())
	return buf.Bytes()
}

// sha256 generates the SHA256 hash of the proof of exponentiation challenge
func (e *epChallenge) sha256() []byte {
	hashF := crypto.SHA256.New()
	defer hashF.Reset()
	hashF.Write(e.serialize())
	sha256Result := hashF.Sum(nil)
	return sha256Result
}

// bigInt serializes the proof of exponentiation challenge to bytes,
// generates the SHA256 hash of the byte data,
// and convert the hash to big integer
func (e *epChallenge) bigInt() *big.Int {
	hashVal := e.sha256()
	return new(big.Int).SetBytes(hashVal)
}

func (e *epChallenge) bigIntPrime() *big.Int {
	hashVal := e.sha256()
	return accumulator.HashToPrime(hashVal)
}

// ExpProver is the prover for proof of exponentiation
type ExpProver struct {
	pp   *PublicParameters
	b    *big.Int
	k    *big.Int
	u    *big.Int
	rhoX *big.Int
	rhoK *big.Int
	x    *big.Int
}

// NewExpProver creates a new prover of proof of exponentiation
func NewExpProver(pp *PublicParameters) *ExpProver {
	// let |G| be N/4, calculate B = (2^(2*lambda))*|G| = N*(2^(2*lambda-2)
	b := new(big.Int).Set(pp.N)
	lsh := 2*securityParam - 2
	b.Lsh(b, uint(lsh))
	return &ExpProver{
		pp: pp,
		b:  b,
	}
}

// Prove generates the proof for proof of exponentiation
func (e *ExpProver) Prove(u, x *big.Int) (*ExponentiationProof, error) {
	e.u = new(big.Int).Set(u)
	e.x = new(big.Int).Set(x)
	challenge, err := e.challenge()
	if err != nil {
		return nil, err
	}
	commitment, err := e.commit(u, x)
	if err != nil {
		return nil, err
	}
	response, err := e.response(challenge)
	if err != nil {
		return nil, err
	}
	return &ExponentiationProof{
		witness:  x,
		commit:   commitment,
		response: response,
	}, nil
}

// chooseRand chooses a random number in [-B, B]
func (e *ExpProver) chooseRand() (*big.Int, error) {
	// random number should be generated in [0, 2B]
	randRange := new(big.Int).Lsh(e.b, 1)
	randRange.Add(randRange, big1)
	num, err := rand.Int(rand.Reader, randRange)
	if err != nil {
		return nil, err
	}
	// num' in [0, 2B], then (num' - B) in [-B, B]
	num.Sub(num, e.b)
	return num, nil
}

// commit generates a commitment provided by the prover
func (e *ExpProver) commit(u, x *big.Int) (*epCommitment, error) {
	e.u = new(big.Int).Set(u)
	e.x = new(big.Int).Set(x)
	// choose random k, rho_x, rho_y from [-B, B]
	var err error
	e.k, err = e.chooseRand()
	if err != nil {
		return nil, err
	}
	e.rhoX, err = e.chooseRand()
	if err != nil {
		return nil, err
	}
	e.rhoK, err = e.chooseRand()
	if err != nil {
		return nil, err
	}
	// z = Com(x; rho_x) =(g^x)(h^rho_x)
	z := com(e.pp, e.x, e.rhoX)
	// A_g = Com(k; rho_k) = (g^k)(h^rho_k)
	aG := com(e.pp, e.k, e.rhoK)
	// A_u =u^k
	aU := new(big.Int).Exp(u, e.k, e.pp.N)
	return &epCommitment{
		z:  z,
		aG: aG,
		aU: aU,
	}, nil
}

// challenge generates the challenge for proof of exponentiation
func (e *ExpProver) challenge() (*epChallenge, error) {
	commit, err := e.commit(e.u, e.x)
	if err != nil {
		return nil, err
	}
	return newEPChallenge(e.pp, commit), nil
}

// response generates the response for proof of exponentiation
func (e *ExpProver) response(challenge *epChallenge) (*epResponse, error) {
	c := challenge.bigInt()
	l := challenge.bigIntPrime()
	// s_x = k + c*x
	sX := new(big.Int).Mul(c, e.x)
	sX.Add(sX, e.k)
	// s_rho = rho_k + c*rho_x
	sRho := new(big.Int).Mul(e.rhoX, c)
	sRho.Add(sRho, e.rhoK)
	// q_x * l  + r_x = s_x
	// q_rho * l + r_rho = s_rho
	qX := new(big.Int).Div(sX, l)
	rX := new(big.Int).Mod(sX, l)
	qRho := new(big.Int).Div(sRho, l)
	rRho := new(big.Int).Mod(sRho, l)
	// Q_g = Com(q_x; q_rho) = (g^q_x)(h^q_rho)
	qG := com(e.pp, qX, qRho)
	// Q_u =u^q_x
	qU := new(big.Int).Exp(e.u, qX, e.pp.N)
	return &epResponse{
		qG:   qG,
		qU:   qU,
		rX:   rX,
		rRho: rRho,
	}, nil
}

// ExpVerifier is the verifier for proof of exponentiation
type ExpVerifier struct {
	pp *PublicParameters
}

// NewExpVerifier creates a new verifier of proof of exponentiation
func NewExpVerifier(pp *PublicParameters) *ExpVerifier {
	return &ExpVerifier{
		pp: pp,
	}
}

// Verify checks the proof of exponentiation
func (e *ExpVerifier) Verify(proof *ExponentiationProof, u, w *big.Int) (bool, error) {
	challenge := e.challenge(proof.commit)
	c := challenge.bigInt()
	l := challenge.bigIntPrime()
	return e.VerifyResponse(c, l, u, w, proof.response, proof.commit)
}

// challenge generates a challenge for the verifier
func (e *ExpVerifier) challenge(commit *epCommitment) *epChallenge {
	challenge := newEPChallenge(e.pp, commit)
	return challenge
}

// VerifyResponse verifies the response of the verifier
func (e *ExpVerifier) VerifyResponse(c, l, u, w *big.Int, response *epResponse, commit *epCommitment) (bool, error) {
	//// check if r_x, r_rho in [left]
	//if response.rX.Cmp(l) >= 0 || response.rRho.Cmp(l) >= 0 {
	//	return false, nil
	//}
	// Q_g^left * Com(r_x; r_rho) = A_g * z^c
	left := new(big.Int).Exp(response.qG, l, e.pp.N)
	left.Mul(left, com(e.pp, response.rX, response.rRho))
	left.Mod(left, e.pp.N)
	right := new(big.Int).Set(commit.aG)
	right.Mul(right, new(big.Int).Exp(commit.z, c, e.pp.N))
	right.Mod(left, e.pp.N)
	if left.Cmp(right) != 0 {
		return false, nil
	}
	// Q_u^left * u^r_x = A_u * w^c
	left = new(big.Int).Exp(response.qU, l, e.pp.N)
	left.Mul(left, new(big.Int).Exp(u, response.rX, e.pp.N))
	left.Mod(left, e.pp.N)
	right = new(big.Int).Set(commit.aU)
	right.Mul(right, new(big.Int).Exp(w, c, e.pp.N))
	right.Mod(left, e.pp.N)
	if left.Cmp(right) != 0 {
		return false, nil
	}
	return true, nil
}

// epCommitment is the commitment for proof of exponentiation sent by the prover
type epCommitment struct {
	z  *big.Int // z = Com(x;rhoX)
	aG *big.Int
	aU *big.Int
}

type epResponse struct {
	qG   *big.Int
	qU   *big.Int
	rX   *big.Int
	rRho *big.Int
}

// com calculates (g^x)(h^r)
func com(pp *PublicParameters, x, r *big.Int) *big.Int {
	res := new(big.Int).Exp(pp.G, x, pp.N)
	res.Mul(res, new(big.Int).Exp(pp.H, r, pp.N))
	res.Mod(res, pp.N)
	return res
}
