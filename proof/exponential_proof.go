package proof

import (
	"crypto/rand"
	"math/big"
)

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

func (e *ExpProver) Prove(u, w, x *big.Int) (*EPCommitment, error) {
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
	return &EPCommitment{
		z:  z,
		aG: aG,
		aU: aU,
	}, nil
}

func (e *ExpProver) Response(challenge *EPChallenge) (*EPResponse, error) {
	// s_x = k + c*x
	sX := new(big.Int).Mul(challenge.c, e.x)
	sX.Add(sX, e.k)
	// s_rho = rho_k + c*rho_x
	sRho := new(big.Int).Mul(e.rhoX, challenge.c)
	sRho.Add(sRho, e.rhoK)
	// q_x * l  + r_x = s_x
	// q_rho * l + r_rho = s_rho
	qX := new(big.Int).Div(sX, challenge.l)
	rX := new(big.Int).Mod(sX, challenge.l)
	qRho := new(big.Int).Div(sRho, challenge.l)
	rRho := new(big.Int).Mod(sRho, challenge.l)
	// Q_g = Com(q_x; q_rho) = (g^q_x)(h^q_rho)
	qG := com(e.pp, qX, qRho)
	// Q_u =u^q_x
	qU := new(big.Int).Exp(e.u, qX, e.pp.N)
	return &EPResponse{
		qG:   qG,
		qU:   qU,
		rX:   rX,
		rRho: rRho,
	}, nil
}

// ExpVerifier is the verifier for proof of exponentiation
type ExpVerifier struct {
	pp         *PublicParameters
	l          *big.Int
	c          *big.Int
	commitment *EPCommitment
}

func NewExpVerifier(pp *PublicParameters) *ExpVerifier {
	return &ExpVerifier{
		pp: pp,
	}
}

func (e *ExpVerifier) SetCommitment(commitment *EPCommitment) {
	e.commitment = commitment
}

func (e *ExpVerifier) Challenge() (*EPChallenge, error) {
	// choose random c in [0, 2^lambda]
	r := new(big.Int).Lsh(big1, securityParam)
	r.Add(r, big1)
	c, err := rand.Int(rand.Reader, r)
	if err != nil {
		return nil, err
	}
	e.c = c
	// get a random primes less than 2^lambda
	p, err := rand.Prime(rand.Reader, securityParam)
	if err != nil {
		return nil, err
	}
	e.l = p
	return &EPChallenge{
		c: c,
		l: p,
	}, nil
}

func (e *ExpVerifier) VerifyResponse(u, w *big.Int, response *EPResponse) (bool, error) {
	// check if r_x, r_rho in [l]
	if response.rX.Cmp(e.l) >= 0 || response.rRho.Cmp(e.l) >= 0 {
		return false, nil
	}
	// Q_g^l * Com(r_x; r_rho) = A_g * z^c
	l := new(big.Int).Exp(response.qG, e.l, e.pp.N)
	l.Mul(l, com(e.pp, response.rX, response.rRho))
	l.Mod(l, e.pp.N)
	r := new(big.Int).Set(e.commitment.aG)
	r.Mul(r, new(big.Int).Exp(e.commitment.z, e.c, e.pp.N))
	r.Mod(l, e.pp.N)
	if l.Cmp(r) != 0 {
		return false, nil
	}
	// Q_u^l * u^r_x = A_u * w^c
	l = new(big.Int).Exp(response.qU, e.l, e.pp.N)
	l.Mul(l, new(big.Int).Exp(u, response.rX, e.pp.N))
	l.Mod(l, e.pp.N)
	r = new(big.Int).Set(e.commitment.aU)
	r.Mul(r, new(big.Int).Exp(w, e.c, e.pp.N))
	r.Mod(l, e.pp.N)
	if l.Cmp(r) != 0 {
		return false, nil
	}
	return true, nil
}

// EPCommitment is the commitment for proof of exponentiation sent by the prover
type EPCommitment struct {
	z  *big.Int // z = Com(x;rhoX)
	aG *big.Int
	aU *big.Int
}

// EPChallenge is the challenge for proof of exponentiation sent by the verifier
type EPChallenge struct {
	c *big.Int
	l *big.Int
}

type EPResponse struct {
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
