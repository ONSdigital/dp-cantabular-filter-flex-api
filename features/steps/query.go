package steps

import (
	"encoding/json"
	"fmt"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
)

func (cf *CantabularFeature) createQueryBody(q string, vars cantabular.QueryData) (string, error) {
	s := struct {
		Q string               `json:"query"`
		V cantabular.QueryData `json:"variables"`
	}{
		Q: q,
		V: vars,
	}

	b, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal query body: %w", err)
	}

	return string(b), nil
}
