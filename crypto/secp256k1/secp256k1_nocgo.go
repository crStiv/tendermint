//go:build !libsecp256k1
// +build !libsecp256k1

package secp256k1

import (
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

// used to reject malleable signatures
// see:
//   - https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/signature_nocgo.go#L90-L93
//   - https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/crypto.go#L39
var secp256k1halfN = new(big.Int).Rsh(secp256k1.S256().N, 1)

// Sign creates an ECDSA signature on curve Secp256k1, using SHA256 on the msg.
// The returned signature will be of the form R || S (in lower-S form).
func (privKey PrivKeySecp256k1) Sign(msg []byte) ([]byte, error) {
	// [peppermint] sign with ethcrypto
	// priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	// sig, err := priv.Sign(crypto.Sha256(msg))
	privateObject, err := ethCrypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, err
	}
	// sigBytes := serializeSig(sig)
	// return sigBytes, nil
	return ethCrypto.Sign(ethCrypto.Keccak256(msg), privateObject)
}

// VerifyBytes verifies a signature of the form R || S.
// It rejects signatures which are not in lower-S form.
func (pubKey PubKeySecp256k1) VerifyBytes(msg []byte, sigStr []byte) bool {
	// if len(sigStr) != 64 {
	// 	return false
	// }
	// pub, err := secp256k1.ParsePubKey(pubKey[:], secp256k1.S256())
	// if err != nil {
	// 	return false
	// }
	// // parse the signature:
	// signature := signatureFromBytes(sigStr)
	// // Reject malleable signatures. libsecp256k1 does this check but btcec doesn't.
	// // see: https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/signature_nocgo.go#L90-L93
	// if signature.S.Cmp(secp256k1halfN) > 0 {
	// 	return false
	// }
	hash := ethCrypto.Keccak256(msg)
	return ethCrypto.VerifySignature(pubKey[:], hash, sigStr[:64])
	// return signature.Verify(crypto.Sha256(msg), pub)
}
