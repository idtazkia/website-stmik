package model

import (
	"github.com/idtazkia/stmik-admission-api/pkg/crypto"
)

// encryptCandidateFields encrypts sensitive fields before storing
func encryptCandidateFields(c *Candidate) error {
	enc := crypto.Get()
	if enc == nil {
		return crypto.ErrNotInitialized
	}

	var err error

	// Deterministic encryption for searchable fields
	c.Email, err = enc.EncryptNullableD(c.Email)
	if err != nil {
		return err
	}
	c.Phone, err = enc.EncryptNullableD(c.Phone)
	if err != nil {
		return err
	}

	// Probabilistic encryption for non-searchable fields
	c.Name, err = enc.EncryptNullableP(c.Name)
	if err != nil {
		return err
	}
	c.Address, err = enc.EncryptNullableP(c.Address)
	if err != nil {
		return err
	}
	c.City, err = enc.EncryptNullableP(c.City)
	if err != nil {
		return err
	}
	c.Province, err = enc.EncryptNullableP(c.Province)
	if err != nil {
		return err
	}
	c.HighSchool, err = enc.EncryptNullableP(c.HighSchool)
	if err != nil {
		return err
	}

	return nil
}

// decryptCandidateFields decrypts sensitive fields after reading
func decryptCandidateFields(c *Candidate) error {
	enc := crypto.Get()
	if enc == nil {
		return crypto.ErrNotInitialized
	}

	var err error

	// Deterministic decryption
	c.Email, err = enc.DecryptNullableD(c.Email)
	if err != nil {
		return err
	}
	c.Phone, err = enc.DecryptNullableD(c.Phone)
	if err != nil {
		return err
	}

	// Probabilistic decryption
	c.Name, err = enc.DecryptNullableP(c.Name)
	if err != nil {
		return err
	}
	c.Address, err = enc.DecryptNullableP(c.Address)
	if err != nil {
		return err
	}
	c.City, err = enc.DecryptNullableP(c.City)
	if err != nil {
		return err
	}
	c.Province, err = enc.DecryptNullableP(c.Province)
	if err != nil {
		return err
	}
	c.HighSchool, err = enc.DecryptNullableP(c.HighSchool)
	if err != nil {
		return err
	}

	return nil
}

// encryptEmail encrypts email for search
func encryptEmail(email string) (string, error) {
	enc := crypto.Get()
	if enc == nil {
		return "", crypto.ErrNotInitialized
	}
	return enc.EncryptDeterministic(email)
}

// encryptPhone encrypts phone for search
func encryptPhone(phone string) (string, error) {
	enc := crypto.Get()
	if enc == nil {
		return "", crypto.ErrNotInitialized
	}
	return enc.EncryptDeterministic(phone)
}

// encryptName encrypts name (probabilistic)
func encryptName(name string) (string, error) {
	enc := crypto.Get()
	if enc == nil {
		return "", crypto.ErrNotInitialized
	}
	return enc.EncryptProbabilistic(name)
}

// decryptName decrypts name
func decryptName(name string) (string, error) {
	enc := crypto.Get()
	if enc == nil {
		return "", crypto.ErrNotInitialized
	}
	return enc.DecryptProbabilistic(name)
}

// encryptNullableP encrypts a nullable string probabilistically
func encryptNullableP(s *string) (*string, error) {
	enc := crypto.Get()
	if enc == nil {
		return nil, crypto.ErrNotInitialized
	}
	return enc.EncryptNullableP(s)
}

// decryptNullableP decrypts a nullable string that was encrypted probabilistically
func decryptNullableP(s *string) (*string, error) {
	enc := crypto.Get()
	if enc == nil {
		return nil, crypto.ErrNotInitialized
	}
	return enc.DecryptNullableP(s)
}

// decryptNullableD decrypts a nullable string that was encrypted deterministically
func decryptNullableD(s *string) (*string, error) {
	enc := crypto.Get()
	if enc == nil {
		return nil, crypto.ErrNotInitialized
	}
	return enc.DecryptNullableD(s)
}

// encryptUserFields encrypts sensitive user fields before storing
func encryptUserFields(email, name, googleID string) (emailEnc, nameEnc string, googleIDEnc *string, err error) {
	enc := crypto.Get()
	if enc == nil {
		return "", "", nil, crypto.ErrNotInitialized
	}

	emailEnc, err = enc.EncryptDeterministic(email)
	if err != nil {
		return "", "", nil, err
	}

	nameEnc, err = enc.EncryptProbabilistic(name)
	if err != nil {
		return "", "", nil, err
	}

	if googleID != "" {
		g, err := enc.EncryptDeterministic(googleID)
		if err != nil {
			return "", "", nil, err
		}
		googleIDEnc = &g
	}

	return emailEnc, nameEnc, googleIDEnc, nil
}

// decryptUserFields decrypts sensitive user fields after reading
func decryptUserFields(email, name string, googleID *string) (emailDec, nameDec string, googleIDDec *string, err error) {
	enc := crypto.Get()
	if enc == nil {
		return "", "", nil, crypto.ErrNotInitialized
	}

	emailDec, err = enc.DecryptDeterministic(email)
	if err != nil {
		return "", "", nil, err
	}

	nameDec, err = enc.DecryptProbabilistic(name)
	if err != nil {
		return "", "", nil, err
	}

	if googleID != nil && *googleID != "" {
		g, err := enc.DecryptDeterministic(*googleID)
		if err != nil {
			return "", "", nil, err
		}
		googleIDDec = &g
	}

	return emailDec, nameDec, googleIDDec, nil
}

// encryptReferrerFields encrypts sensitive referrer fields
func encryptReferrerFields(name string, email, phone, bankName, bankAccount, accountHolder, institution *string) (
	nameEnc string, emailEnc, phoneEnc, bankNameEnc, bankAccountEnc, accountHolderEnc, institutionEnc *string, err error) {
	enc := crypto.Get()
	if enc == nil {
		return "", nil, nil, nil, nil, nil, nil, crypto.ErrNotInitialized
	}

	nameEnc, err = enc.EncryptProbabilistic(name)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	emailEnc, err = enc.EncryptNullableP(email)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	phoneEnc, err = enc.EncryptNullableP(phone)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	bankNameEnc, err = enc.EncryptNullableP(bankName)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	bankAccountEnc, err = enc.EncryptNullableP(bankAccount)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	accountHolderEnc, err = enc.EncryptNullableP(accountHolder)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	institutionEnc, err = enc.EncryptNullableP(institution)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	return nameEnc, emailEnc, phoneEnc, bankNameEnc, bankAccountEnc, accountHolderEnc, institutionEnc, nil
}

// decryptReferrerFields decrypts sensitive referrer fields
func decryptReferrerFields(name string, email, phone, bankName, bankAccount, accountHolder, institution *string) (
	nameDec string, emailDec, phoneDec, bankNameDec, bankAccountDec, accountHolderDec, institutionDec *string, err error) {
	enc := crypto.Get()
	if enc == nil {
		return "", nil, nil, nil, nil, nil, nil, crypto.ErrNotInitialized
	}

	nameDec, err = enc.DecryptProbabilistic(name)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	emailDec, err = enc.DecryptNullableP(email)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	phoneDec, err = enc.DecryptNullableP(phone)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	bankNameDec, err = enc.DecryptNullableP(bankName)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	bankAccountDec, err = enc.DecryptNullableP(bankAccount)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	accountHolderDec, err = enc.DecryptNullableP(accountHolder)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	institutionDec, err = enc.DecryptNullableP(institution)
	if err != nil {
		return "", nil, nil, nil, nil, nil, nil, err
	}

	return nameDec, emailDec, phoneDec, bankNameDec, bankAccountDec, accountHolderDec, institutionDec, nil
}

// encryptToken encrypts verification token
func encryptToken(token string) (string, error) {
	enc := crypto.Get()
	if enc == nil {
		return "", crypto.ErrNotInitialized
	}
	return enc.EncryptDeterministic(token)
}

// decryptToken decrypts verification token
func decryptToken(token string) (string, error) {
	enc := crypto.Get()
	if enc == nil {
		return "", crypto.ErrNotInitialized
	}
	return enc.DecryptDeterministic(token)
}
