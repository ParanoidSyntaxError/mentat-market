package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/web"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

func initFowardersSubCmds(s *Shell) []cli.Command {
	return []cli.Command{
		{
			Name:   "list",
			Usage:  "List all stored forwarders addresses",
			Action: s.ListForwarders,
		},
		{
			Name:   "track",
			Usage:  "Track a new forwarder",
			Action: s.TrackForwarder,
			Flags: []cli.Flag{
				cli.Int64Flag{
					Name:  "evm-chain-id, evmChainID, c",
					Usage: "chain ID, if left empty, EVM.ChainID will be used",
				},
				cli.StringFlag{
					Name:  "address, a",
					Usage: "The forwarding address (in hex format)",
				},
			},
		},
		{
			Name:   "delete",
			Usage:  "Delete a forwarder address",
			Action: s.DeleteForwarder,
		},
	}
}

type EVMForwarderPresenter struct {
	JAID // This is needed to render the id for a JSONAPI Resource as normal JSON
	presenters.EVMForwarderResource
}

var evmFwdsHeaders = []string{"ID", "Address", "Chain ID", "Created At"}

// ToRow presents the EVMForwarderResource as a slice of strings.
func (p *EVMForwarderPresenter) ToRow() []string {
	row := []string{
		p.GetID(),
		p.Address.String(),
		p.EVMChainID.ToInt().String(),
		p.CreatedAt.Format(time.RFC3339),
	}
	return row
}

// RenderTable implements TableRenderer
func (p *EVMForwarderPresenter) RenderTable(rt RendererTable) error {
	var rows [][]string
	rows = append(rows, p.ToRow())
	renderList(evmFwdsHeaders, rows, rt.Writer)

	return nil
}

// EVMForwarderPresenters implements TableRenderer for a slice of EVMForwarderPresenter.
type EVMForwarderPresenters []EVMForwarderPresenter

// RenderTable implements TableRenderer
func (ps EVMForwarderPresenters) RenderTable(rt RendererTable) error {
	var rows [][]string

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	renderList(evmFwdsHeaders, rows, rt.Writer)

	return nil
}

// ListForwarders list all forwarder addresses tracked by node
func (s *Shell) ListForwarders(c *cli.Context) (err error) {
	return s.getPage("/v2/nodes/evm/forwarders", c.Int("page"), &EVMForwarderPresenters{})
}

// DeleteForwarder deletes forwarder address from node db by id.
func (s *Shell) DeleteForwarder(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("must pass the forwarder id to be archived"))
	}
	resp, err := s.HTTP.Delete(s.ctx(), "/v2/nodes/evm/forwarders/"+c.Args().First())
	if err != nil {
		return s.errorOut(err)
	}
	_, err = s.parseResponse(resp)
	if err != nil {
		return s.errorOut(err)
	}

	fmt.Printf("Forwarder %v Deleted\n", c.Args().First())
	return nil
}

// TrackForwarder tracks forwarder address in db.
func (s *Shell) TrackForwarder(c *cli.Context) (err error) {
	addressHex := c.String("address")
	chainIDStr := c.String("evm-chain-id")

	addressBytes, err := hexutil.Decode(addressHex)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "could not decode address"))
	}
	address := gethCommon.BytesToAddress(addressBytes)

	var chainID *big.Int
	if chainIDStr != "" {
		var ok bool
		chainID, ok = big.NewInt(0).SetString(chainIDStr, 10)
		if !ok {
			return s.errorOut(errors.Wrap(err, "invalid evm-chain-id"))
		}
	}

	request, err := json.Marshal(web.TrackEVMForwarderRequest{
		EVMChainID: (*ubig.Big)(chainID),
		Address:    address,
	})
	if err != nil {
		return s.errorOut(err)
	}

	resp, err := s.HTTP.Post(s.ctx(), "/v2/nodes/evm/forwarders/track", bytes.NewReader(request))
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode >= 400 {
		body, rerr := io.ReadAll(resp.Body)
		if err != nil {
			err = multierr.Append(err, rerr)
			return s.errorOut(err)
		}
		fmt.Printf("Response: '%v', Status: %d\n", string(body), resp.StatusCode)
		return s.errorOut(err)
	}

	err = s.renderAPIResponse(resp, &EVMForwarderPresenter{}, "Forwarder created")
	return err
}
