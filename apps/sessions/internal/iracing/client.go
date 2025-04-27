package iracing

import (
	"log"

	"github.com/markphelps/optional"
	irapi "github.com/riccardotornesello/irapi-go"
	"github.com/riccardotornesello/irapi-go/api/results"
	"riccardotornesello.it/sharedtelemetry/iracing/sessions/config"
)

type Client struct {
	*irapi.IRacingApiClient
}

func NewClient(cfg config.Config) (*Client, error) {
	irClient, err := irapi.NewIRacingApiClient(cfg.IRacingEmail, cfg.IRacingPassword)
	if err != nil {
		return nil, err
	}
	return &Client{irClient}, nil
}

func (c *Client) FetchResults(subsessionId int) (*results.ResultsGetResponse, error) {
	log.Println("Fetching results")
	includeLicenses := optional.NewBool(false)
	res, err := c.Results.GetResults(results.ResultsGetParams{
		SubsessionId:    subsessionId,
		IncludeLicenses: &includeLicenses,
	})

	return res, err
}

func (c *Client) FetchDriverResults(subsessionId int, simsessionNumber int, custId int) (*results.ResultsLapDataResponse, error) {
	log.Println("Fetching driver results for", custId)
	custIdOpt := optional.NewInt(custId)
	res, err := c.Results.GetResultsLapData(results.ResultsLapDataParams{
		SubsessionId:     subsessionId,
		SimsessionNumber: simsessionNumber,
		CustId:           &custIdOpt,
	})

	return res, err
}
