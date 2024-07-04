package uuid

import "github.com/google/uuid"

type Generator struct{}

func NewGenerator() Generator {
    return Generator{}
}

func (generator Generator) GenerateUuidString() (string, error) {
    u, err := uuid.NewRandom()
    if err != nil {
        return "", err
    }

    return u.String(), nil
}
