package libkb

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	keybase1 "github.com/keybase/client/protocol/go"
	jsonw "github.com/keybase/go-jsonw"
)

func UidFromHex(s string) (keybase1.UID, error) {
	u, err := keybase1.UIDFromString(s)
	if err != nil {
		var nilUID keybase1.UID
		return nilUID, err
	}
	return u, nil
}

func GetUID(w *jsonw.Wrapper) (keybase1.UID, error) {
	s, err := w.GetString()
	var nilUID keybase1.UID
	if err != nil {
		return nilUID, err
	}
	return UidFromHex(s)
}

func GetUidVoid(w *jsonw.Wrapper, u *keybase1.UID, e *error) {
	uid, err := GetUID(w)
	if err == nil {
		*u = uid
	} else if e != nil && *e == nil {
		*e = err
	}
}

func UIDWrapper(uid keybase1.UID) *jsonw.Wrapper {
	return jsonw.NewString(uid.String())
}

func UIDArg(uid keybase1.UID) HttpValue {
	return S{Val: uid.String()}
}

const (
	UID_LEN      = keybase1.UID_LEN
	UID_SUFFIX   = keybase1.UID_SUFFIX
	UID_SUFFIX_2 = keybase1.UID_SUFFIX_2
)

// UsernameToUID works for users created after "Fri Feb  6 19:33:08 EST 2015"
func UsernameToUID(s string) keybase1.UID {
	h := sha256.Sum256([]byte(strings.ToLower(s)))
	var uid [keybase1.UID_LEN]byte
	copy(uid[:], h[0:UID_LEN-1])
	uid[UID_LEN-1] = UID_SUFFIX_2
	ret, _ := keybase1.UIDFromString(hex.EncodeToString(uid[:]))
	return ret
}

func CheckUIDAgainstUsername(uid keybase1.UID, username string) (err error) {
	u2 := UsernameToUID(username)
	if uid.NotEqual(u2) {
		err = UidMismatchError{fmt.Sprintf("%s != %s (via %s)", uid, u2, username)}
	}
	return
}
