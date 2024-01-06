package repo

import (
	"context"
	"fmt"
	"github.com/KushnerykPavel/mab-recomendation-api/internal/models"
	"github.com/lib/pq"
)

func (r *Repo) GetLastInteraction(ctx context.Context, relationName string, userID string, entityID int) (*models.Interaction, error) {
	var result models.Interaction

	query := fmt.Sprintf(`select user_id, entity_id, alpha, beta from users_%s_interactions
where user_id = $1 and entity_id = $2
order by created_at desc
limit 1;`, relationName)

	if err := r.db.GetContext(ctx, &result, query, userID, entityID); err != nil {
		return &models.Interaction{
			UserID:   userID,
			EntityID: entityID,
			Alpha:    1.0,
			Beta:     1.0,
		}, err
	}

	return &result, nil
}

// GetLastInteractions return list of las week user interactions with top 10 entities
func (r *Repo) GetLastInteractions(ctx context.Context, relationName models.EntityType, userID string) ([]*models.Interaction, error) {
	result := make([]*models.Interaction, 0)

	query := fmt.Sprintf(`select entity_id, alpha, beta from users_%s_interactions
where user_id = $1 and created_at > now() - interval '7 day'
order by alpha desc
limit 10;`, relationName)

	if err := r.db.SelectContext(ctx, &result, query, userID); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repo) SetInteraction(ctx context.Context, relationName string, interaction *models.Interaction) error {
	var err error

	query := fmt.Sprintf(`insert into users_%s_interactions (user_id, entity_id, alpha, beta, created_at) 
	values (:user_id, :entity_id, :alpha, :beta, :created_at)
ON CONFLICT ON CONSTRAINT users_%s_interactions_pk
    DO
        UPDATE SET (created_at, alpha, beta) = (:created_at, :alpha, :beta);`, relationName, relationName)

	_, err = r.db.NamedExecContext(ctx, query, interaction)

	return err
}

func (r *Repo) GetMovieCandidatesForEntities(ctx context.Context, relationName models.EntityType, userID string, entityIDList []int) ([]*models.Interaction, error) {
	result := make([]*models.Interaction, 0)

	query := fmt.Sprintf(`select distinct id as entity_id, coalesce(umi.alpha, 1) as alpha , coalesce(umi.beta, 1) as beta from movies
left join movies_%s_relation mgr on movies.id = mgr.source_id
left join users_movies_interactions umi on movies.id = umi.entity_id
where id not in (select entity_id from users_movies_interactions where user_id = $1 and created_at > now() - interval '1 day')
and  mgr.destination_id = any($2)
order by 2 desc
limit 10;`, relationName)

	if err := r.db.SelectContext(ctx, &result, query, userID, pq.Array(entityIDList)); err != nil {
		return nil, err
	}

	return result, nil
}
