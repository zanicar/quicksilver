package balancer

import (
	"encoding/json"
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type balancerPoolPretty struct {
	Address            sdk.AccAddress `json:"address" yaml:"address"`
	Id                 uint64         `json:"id" yaml:"id"`
	PoolParams         PoolParams     `json:"pool_params" yaml:"pool_params"`
	FuturePoolGovernor string         `json:"future_pool_governor" yaml:"future_pool_governor"`
	TotalWeight        sdk.Dec        `json:"total_weight" yaml:"total_weight"`
	TotalShares        sdk.Coin       `json:"total_shares" yaml:"total_shares"`
	PoolAssets         []PoolAsset    `json:"pool_assets" yaml:"pool_assets"`
}

func (p Pool) String() string {
	out, err := p.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(out)
}

// MarshalJSON returns the JSON representation of a Pool.
func (p Pool) MarshalJSON() ([]byte, error) {
	if len(strings.TrimSpace(p.Address)) == 0 {
		return nil, errors.New("empty address string is not allowed")
	}

	accAddr, err := sdk.GetFromBech32(p.Address, "osmo")
	// accAddr, err := sdk.AccAddressFromBech32(p.Address)
	if err != nil {
		return nil, err
	}

	decTotalWeight := sdk.NewDecFromInt(p.TotalWeight)

	return json.Marshal(balancerPoolPretty{
		Address:            accAddr,
		Id:                 p.Id,
		PoolParams:         p.PoolParams,
		FuturePoolGovernor: p.FuturePoolGovernor,
		TotalWeight:        decTotalWeight,
		TotalShares:        p.TotalShares,
		PoolAssets:         p.PoolAssets,
	})
}

// UnmarshalJSON unmarshals raw JSON bytes into a Pool.
func (p *Pool) UnmarshalJSON(bz []byte) error {
	var alias balancerPoolPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	p.Address = alias.Address.String()
	p.Id = alias.Id
	p.PoolParams = alias.PoolParams
	p.FuturePoolGovernor = alias.FuturePoolGovernor
	p.TotalWeight = alias.TotalWeight.RoundInt()
	p.TotalShares = alias.TotalShares
	p.PoolAssets = alias.PoolAssets

	return nil
}
