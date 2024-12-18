package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minireipaz/pkg/config"
	"minireipaz/pkg/domain/models"
	"net/http"
	"net/url"
)

type DashboardRepository struct {
	client          HTTPClient
	databaseHTTPURL string
	token           string
}

func NewDashboardRepository(client HTTPClient, clickhouseConfig config.ClickhouseConfig) *DashboardRepository {
	return &DashboardRepository{
		client:          client,
		databaseHTTPURL: clickhouseConfig.GetClickhouseURI(),
		token:           clickhouseConfig.GetClickhouseToken(),
	}
}

func (d *DashboardRepository) GetLastWorkflowData(userID string, limitCount uint64) (models.InfoDashboard, error) {
	u, err := url.Parse(d.databaseHTTPURL + "/user_workflow_stats.json")
	if err != nil {
		return models.InfoDashboard{}, err
	}

	q := u.Query()
	q.Set("token", d.token)
	q.Set("user_id", userID)
	q.Set("limit_count", fmt.Sprintf("%d", limitCount))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return models.InfoDashboard{}, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return models.InfoDashboard{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return models.InfoDashboard{}, fmt.Errorf("ERROR | response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.InfoDashboard
	// if err := json.Unmarshal(bodyBytes, &result); err != nil {
	// 	log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
	// 	return models.InfoDashboard{}, fmt.Errorf("ERROR | cannot decode token: %v", err)
	// }
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR | cannot decode body: %s %v", string(bodyBytes), err)
		return models.InfoDashboard{}, fmt.Errorf("ERROR | cannot decode token: %v", err)
	}

	return result, nil
}
