package repo

import (
	"context"
	"fmt"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
)

func (r *Repo) GetRelation(ctx context.Context, relationName string, destinationID int) ([]*models.Relation, error) {
	result := make([]*models.Relation, 0)

	query := fmt.Sprintf(`select r.id as id, r.name as name from %s r
		left join movies_%s_relation mr on mr.destination_id = r.id
		where mr.source_id = $1`, relationName, relationName)

	if err := r.db.SelectContext(ctx, &result, query, destinationID); err != nil {
		return nil, err
	}

	return result, nil
}
