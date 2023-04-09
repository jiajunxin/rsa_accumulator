// Package proof zkPoKE
// Protocol zkPoKE for R_{PoKE}
// Paper: Batching Techniques for Accumulators with Applications to IOPs and Stateless Blockchains
// Link: https://eprint.iacr.org/2018/1188.pdf
package proof

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"math/big"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

const (
	zkPoKEChallengeStatement = "u^x = w, x is an integer"
)

// ZKPoKEProof represents the proof of knowledge of exponentiation
type ZKPoKEProof struct {
	witness  *big.Int // witness x in Z, u^x = w
	commit   *zkPoKECommitment
	response *zkPoKEResponse
}

type zkPoKEChallenge struct {
	statement  string            // the statement for the challenge
	g, h, n    *big.Int          // public parameters: G, H, N
	commitment *zkPoKECommitment // the commitment of by the prover
}

// newZKPoKEChallenge generates a new challenge for proof of exponentiation
func newZKPoKEChallenge(pp *PublicParameters, commitment *zkPoKECommitment) *zkPoKEChallenge {
	return &zkPoKEChallenge{
		statement:  zkPoKEChallengeStatement,
		g:          pp.G,
		h:          pp.H,
		n:          pp.N,
		commitment: commitment,
	}
}

// Serialize generates the serialized data for proof of exponentiation challenge in byte format
func (zc *zkPoKEChallenge) serialize() []byte {
	var buf bytes.Buffer
	defer buf.Reset()
	buf.WriteString(zc.statement)
	buf.WriteString(zc.g.String())
	buf.WriteString(zc.h.String())
	buf.WriteString(zc.n.String())
	buf.WriteString(zc.commitment.z.String())
	buf.WriteString(zc.commitment.aG.String())
	buf.WriteString(zc.commitment.aU.String())
	return buf.Bytes()
}

// sha256 generates the SHA256 hash of the proof of exponentiation challenge
func (zc *zkPoKEChallenge) sha256() []byte {
	hashF := crypto.SHA256.New()
	defer hashF.Reset()
	_, err := hashF.Write(zc.serialize())
	if err != nil {
		panic(err)
	}
	sha256Result := hashF.Sum(nil)
	return sha256Result
}

// bigInt serializes the proof of exponentiation challenge to bytes,
// generates the SHA256 hash of the byte data,
// and convert the hash to big integer
func (zc *zkPoKEChallenge) bigInt() *big.Int {
	hashVal := zc.sha256()
	return new(big.Int).SetBytes(hashVal)
}

func (zc *zkPoKEChallenge) bigIntPrime() *big.Int {
	hashVal := zc.sha256()
	return accumulator.HashToPrime(hashVal)
}

// ZKPoKEProver is the prover for proof of exponentiation
type ZKPoKEProver struct {
	pp   *PublicParameters
	b    *big.Int
	k    *big.Int
	u    *big.Int
	rhoX *big.Int
	rhoK *big.Int
	x    *big.Int
}

// NewZKPoKEProver creates a new prover of proof of exponentiation
func NewZKPoKEProver(pp *PublicParameters) *ZKPoKEProver {
	// let |G| be N/4, calculate B = (2^(2*lambda))*|G| = N*(2^(2*lambda-2)
	b := new(big.Int).Set(pp.N)
	lsh := 2*securityParam - 2
	b.Lsh(b, uint(lsh))
	return &ZKPoKEProver{
		pp: pp,
		b:  b,
	}
}

// Prove generates the proof for proof of exponentiation
func (zp *ZKPoKEProver) Prove(u, x *big.Int) (*ZKPoKEProof, error) {
	zp.u = new(big.Int).Set(u)
	zp.x = new(big.Int).Set(x)
	challenge, err := zp.challenge()
	if err != nil {
		return nil, err
	}
	commitment, err := zp.commit(u, x)
	if err != nil {
		return nil, err
	}
	response, err := zp.response(challenge)
	if err != nil {
		return nil, err
	}
	return &ZKPoKEProof{
		witness:  x,
		commit:   commitment,
		response: response,
	}, nil
}

// chooseRand chooses a random number in [-B, B]
func (zp *ZKPoKEProver) chooseRand() (*big.Int, error) {
	// random number should be generated in [0, 2B]
	randRange := iPool.Get().(*big.Int)
	defer iPool.Put(randRange)
	randRange.Lsh(zp.b, 1)
	randRange.Add(randRange, big1)
	num, err := rand.Int(rand.Reader, randRange)
	if err != nil {
		return nil, err
	}
	// num' in [0, 2B], then (num' - B) in [-B, B]
	num.Sub(num, zp.b)
	return num, nil
}

// commit generates a commitment provided by the prover
func (zp *ZKPoKEProver) commit(u, x *big.Int) (*zkPoKECommitment, error) {
	zp.u = new(big.Int).Set(u)
	zp.x = new(big.Int).Set(x)
	// choose random k, rho_x, rho_y from [-B, B]
	var err error
	zp.k, err = zp.chooseRand()
	if err != nil {
		return nil, err
	}
	zp.rhoX, err = zp.chooseRand()
	if err != nil {
		return nil, err
	}
	zp.rhoK, err = zp.chooseRand()
	if err != nil {
		return nil, err
	}
	// z = Com(x; rho_x) =(g^x)(h^rho_x)
	z := com(zp.pp, zp.x, zp.rhoX)
	// A_g = Com(k; rho_k) = (g^k)(h^rho_k)
	aG := com(zp.pp, zp.k, zp.rhoK)
	// A_u =u^k
	aU := new(big.Int).Exp(u, zp.k, zp.pp.N)
	return &zkPoKECommitment{
		z:  z,
		aG: aG,
		aU: aU,
	}, nil
}

// challenge generates the challenge for proof of exponentiation
func (zp *ZKPoKEProver) challenge() (*zkPoKEChallenge, error) {
	commit, err := zp.commit(zp.u, zp.x)
	if err != nil {
		return nil, err
	}
	return newZKPoKEChallenge(zp.pp, commit), nil
}

// response generates the response for proof of exponentiation
func (zp *ZKPoKEProver) response(challenge *zkPoKEChallenge) (*zkPoKEResponse, error) {
	c := challenge.bigInt()
	l := challenge.bigIntPrime()
	// s_x = k + c*x
	sX := iPool.Get().(*big.Int)
	defer iPool.Put(sX)
	sX.Mul(c, zp.x)
	sX.Add(sX, zp.k)
	// s_rho = rho_k + c*rho_x
	sRho := iPool.Get().(*big.Int)
	defer iPool.Put(sRho)
	sRho.Mul(zp.rhoX, c)
	sRho.Add(sRho, zp.rhoK)
	// q_x * l  + r_x = s_x
	// q_rho * l + r_rho = s_rho
	qX := iPool.Get().(*big.Int)
	defer iPool.Put(qX)
	qX.Div(sX, l)
	rX := new(big.Int).Mod(sX, l)
	qRho := iPool.Get().(*big.Int)
	defer iPool.Put(qRho)
	qRho.Div(sRho, l)
	rRho := new(big.Int).Mod(sRho, l)
	// Q_g = Com(q_x; q_rho) = (g^q_x)(h^q_rho)
	qG := com(zp.pp, qX, qRho)
	// Q_u =u^q_x
	qU := new(big.Int).Exp(zp.u, qX, zp.pp.N)
	return &zkPoKEResponse{
		qG:   qG,
		qU:   qU,
		rX:   rX,
		rRho: rRho,
	}, nil
}

// ZKPoKEVerifier is the verifier for proof of exponentiation
type ZKPoKEVerifier struct {
	pp *PublicParameters
}

// NewZKPoKEVerifier creates a new verifier of proof of exponentiation
func NewZKPoKEVerifier(pp *PublicParameters) *ZKPoKEVerifier {
	return &ZKPoKEVerifier{
		pp: pp,
	}
}

// Verify checks the proof of exponentiation
func (zv *ZKPoKEVerifier) Verify(proof *ZKPoKEProof, u, w *big.Int) (bool, error) {
	challenge := zv.challenge(proof.commit)
	c := challenge.bigInt()
	l := challenge.bigIntPrime()
	return zv.VerifyResponse(c, l, u, w, proof.response, proof.commit)
}

// challenge generates a challenge for the verifier
func (zv *ZKPoKEVerifier) challenge(commit *zkPoKECommitment) *zkPoKEChallenge {
	challenge := newZKPoKEChallenge(zv.pp, commit)
	return challenge
}

// VerifyResponse verifies the response of the verifier
func (zv *ZKPoKEVerifier) VerifyResponse(c, l, u, w *big.Int, response *zkPoKEResponse, commit *zkPoKECommitment) (bool, error) {
	//// check if r_x, r_rho in [l]
	//if response.rX.Cmp(l) >= 0 || response.rRho.Cmp(l) >= 0 {
	//	return false, nil
	//}
	// Q_g^l * Com(r_x; r_rho) = A_g * zv^c
	lhs := iPool.Get().(*big.Int)
	defer iPool.Put(lhs)
	lhs.Exp(response.qG, l, zv.pp.N)
	lhs.Mul(lhs, com(zv.pp, response.rX, response.rRho))
	lhs.Mod(lhs, zv.pp.N)
	rhs := iPool.Get().(*big.Int)
	defer iPool.Put(rhs)
	rhs.Set(commit.aG)
	rhs.Mul(rhs, new(big.Int).Exp(commit.z, c, zv.pp.N))
	rhs.Mod(lhs, zv.pp.N)
	if lhs.Cmp(rhs) != 0 {
		return false, nil
	}
	// Q_u^l * u^r_x = A_u * w^c
	lhs.Exp(response.qU, l, zv.pp.N)
	lhs.Mul(lhs, new(big.Int).Exp(u, response.rX, zv.pp.N))
	lhs.Mod(lhs, zv.pp.N)
	rhs.Set(commit.aU)
	rhs.Mul(rhs, new(big.Int).Exp(w, c, zv.pp.N))
	rhs.Mod(lhs, zv.pp.N)
	if lhs.Cmp(rhs) != 0 {
		return false, nil
	}
	return true, nil
}

// zkPoKECommitment is the commitment for proof of exponentiation sent by the prover
type zkPoKECommitment struct {
	z  *big.Int // z = Com(x;rhoX)
	aG *big.Int
	aU *big.Int
}

type zkPoKEResponse struct {
	qG   *big.Int
	qU   *big.Int
	rX   *big.Int
	rRho *big.Int
}

// com calculates (g^x)(h^r)
func com(pp *PublicParameters, x, r *big.Int) *big.Int {
	res := new(big.Int).Exp(pp.G, x, pp.N)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	res.Mul(res, opt.Exp(pp.H, r, pp.N))
	res.Mod(res, pp.N)
	return res
}
