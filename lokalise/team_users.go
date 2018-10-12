package lokalise

import (
	"context"
	"fmt"

	"github.com/go-resty/resty"
	"github.com/lokalise/lokalise-go-sdk/model"
)

type TeamUsersService struct {
	httpClient *resty.Client
}

const (
	pathTeams = "teams"
)

func pathTeamUsers(teamID int64) string {
	return fmt.Sprintf("%s/%d/users", pathTeams, teamID)
}

func (c *TeamUsersService) Retrieve(ctx context.Context, teamID, userID int64) (model.TeamUserResponse, error) {
	var res model.TeamUserResponse
	req := c.httpClient.R().
		SetResult(&res).
		SetContext(ctx)
	resp, err := req.Get(fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID))
	if err != nil {
		return model.TeamUserResponse{}, err
	}
	return res, apiError(resp)
}
