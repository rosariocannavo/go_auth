package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func CheckSig(from, sigHex string, msg []byte) bool {
	sig := hexutil.MustDecode(sigHex)

	msg = accounts.TextHash(msg)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	fmt.Printf("ECDSA Signature: %x\n", sig)
	fmt.Printf("  R: %x\n", sig[0:32])  // 32 bytes
	fmt.Printf("  S: %x\n", sig[32:64]) // 32 bytes
	fmt.Printf("  V: %x\n", sig[64:])

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	return strings.EqualFold(from, recoveredAddr.Hex())
}

func GenerateRandomNonce() (string, error) {
	// Define the length of the nonce
	nonceLength := 16

	// Create a byte slice to store the random nonce
	nonce := make([]byte, nonceLength)

	// Read random bytes into the nonce slice
	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to a hexadecimal string
	nonceString := "0x" + hex.EncodeToString(nonce)

	return nonceString, nil
}
