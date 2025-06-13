package passwords

import "golang.org/x/crypto/bcrypt"

type PasswordService struct {}

func (passwordService PasswordService) HashPassword(plaintextPassword []byte) ([]byte, error) {
    return bcrypt.GenerateFromPassword(plaintextPassword, bcrypt.DefaultCost)
}
