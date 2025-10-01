package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type Keyset struct {
	Private string
	Public  string
}

func main() {
	keyset, err := getOrInitKeyset()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Keyset: %s\n", keyset)
}

func getOrInitKeyset() (*Keyset, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	keyfile := path.Join(home, ".nostrkeys")

	file, err := os.ReadFile(keyfile)

	if err == nil {
		var keyset Keyset

		err := json.Unmarshal(file, &keyset)

		if err != nil {
			return nil, err
		}

		return &keyset, nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	keyset, err := createKeyset()

	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(keyset)

	if err != nil {
		return nil, err
	}

	err = os.WriteFile(keyfile, b, 0600)

	if err != nil {
		return nil, err
	}

	return keyset, nil
}

func createKeyset() (*Keyset, error) {
	sk := nostr.GeneratePrivateKey()
	pk, err := nostr.GetPublicKey(sk)

	if err != nil {
		return nil, err
	}

	nsec, err := nip19.EncodePrivateKey(sk)

	if err != nil {
		return nil, err
	}

	npub, err := nip19.EncodePublicKey(pk)

	if err != nil {
		return nil, err
	}

	return &Keyset{
		Private: nsec,
		Public:  npub,
	}, nil
}
