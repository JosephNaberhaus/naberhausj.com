package cache

import (
	"encoding/json"
	"fmt"
)

type Artifact interface{}

func UnmarshalArtifact(artifact Artifact, into interface{}) error {
	marshalled, err := json.Marshal(artifact)
	if err != nil {
		return fmt.Errorf("error marshalling cache artifact: %w", err)
	}

	err = json.Unmarshal(marshalled, into)
	if err != nil {
		return fmt.Errorf("error unmarshalling cache artifact: %w", err)
	}

	return nil
}
