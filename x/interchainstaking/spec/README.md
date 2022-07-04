# Interchain Staking

## Abstract

Module, `x/interchainstaking`, defines and implements the core Quicksilver
protocol. _(wip)_

## Contents

1. [Concepts](#Concepts)
1. [State](#State)
1. [Messages](#Messages)
1. [Transactions](#Transactions)
1. [Events](#Events)
1. [Hooks](#Hooks)
1. [Queries](#Queries)
1. [Keepers](#Keepers)
1. [Parameters](#Parameters)
1. [Begin Block](#Begin-Block)
1. [After Epoch End](#After-Epoch-End)

## Concepts

Key concepts, mechanisms and core logic for the module;

## State

### RegisteredZone

A `RegisteredZone` represents a Cosmos based blockchain (or zone) that
integrates with the Quicksilver protocol via Interchain Accounts (ICS) and
Interblockchain Communication (IBC).

```
type RegisteredZone struct {
    ConnectionId                 string
    ChainId                      string
    DepositAddress               *ICAAccount
    WithdrawalAddress            *ICAAccount
    PerformanceAddress           *ICAAccount
    DelegationAddresses          []*ICAAccount
    AccountPrefix                string
    LocalDenom                   string
    BaseDenom                    string
    RedemptionRate               sdk.Dec
    LastRedemptionRate           sdk.Dec
    Validators                   []*Validator
    AggregateIntent              map[string]*ValidatorIntent
    MultiSend                    bool
    LiquidityModule              bool
    WithdrawalWaitgroup          uint32
    IbcNextValidatorsHash        []byte
    ValidatorSelectionAllocation sdk.Coins
    HoldingsAllocation           sdk.Coins
}
```

* **ConnectionId** - remote zone connection identifier;
* **ChainId** - remote zone identifier;
* **DepositAddress** - remote zone deposit address;
* **WithdrawalAddress** - remote zone withdrawal address;
* **PerformanceAddress** - remote zone performance address (each validator gets an exact equal delegation from this account to measure performance);
* **DelegationAddresses** - remote zone delegation addresses to represent granular voting power;
* **AccountPrefix** - remote zone account address prefix;
* **LocalDenom** - protocol denomination (qAsset), e.g. uqatom;
* **BaseDenom** - remote zone denomination (uStake), e.g. uatom;
* **RedemptionRate** - redemption rate between protocol qAsset and native remote asset;
* **LastRedemptionRate** - redemption rate as at previous epoch boundary (used to prevent epoch boundary gaming);
* **Validators** - list of validators on the remote zone;
* **AggregateIntent** - the aggregated delegation intent of the protocol for this remote zone;
* **MultiSend** - multisend support on remote zone;
* **LiquidityModule** - liquidity module enabled on remote zone;
* **WithdrawalWaitgroup** - tally of pending withdrawal transactions;
* **IbcNextValidatorHash** - 
* **ValidatorSelectionAllocation** - proportional zone rewards allocation for validator selection;
* **HoldingsAllocation** - proportional zone rewards allocation for asset holdings;

### ICAAccount

An `ICAAccount` represents an account on an remote zone under the control of
the protocol.

```
type ICAAccount struct {
    Address          string
    Balance          sdk.Coins
    DelegatedBalance sdk.Coin
    PortName         string
    BalanceWaitgroup uint32
}
```

* **Address** - the account address on the remote zone;
* **Balance** - the account balance on the remote zone;
* **DelegatedBalance** - the account delegation balance on the remote zone;
* **PortName** - the port name to access the remote zone;
* **BalanceWaitgroup** - tally of pending balance query transactions sent to the remote zone;

### Validator

`Validator` represents relevant meta data of a validator within a zone.

```
type Validator struct {
	ValoperAddress  string
	CommissionRate  sdk.Dec
	DelegatorShares sdk.Dec
	VotingPower     sdk.Int
	Score           sdk.Dec
}
```

* **ValoperAddress** - validator address;
* **CommissionRate** - validator commission rate;
* **DelegatorShares** - 
* **VotingPower** - validator voting power on the remote zone;
* **Score** - validator Quicksilver protocol score (decentralization * performance);

### ValidatorIntent

```
type ValidatorIntent struct {
	ValoperAddress string
	Weight         sdk.Dec
}
```

* **ValoperAddress** - remote zone validator address;
* **Weight** - weight of intended delegation to this validator;

## Messages

Description of message types that trigger state transitions;

## Transactions

### signal-intent

Signal validator delegation intent by providing a comma seperated string
containing a decimal weight and the bech32 validator address.

`quicksilverd signal-intent [chain_id] 0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz`

### redeem

Redeem tokens.

`quicksilverd redeem [coins] [destination_address]`

### register-zone

Submit a zone registration proposal.

`quicksilverd register-zone [proposal-file]`

The proposal must include an initial deposit and the details must be provided
as a json file, e.g.

```
{
  "title": "Register cosmoshub-4",
  "description": "Onboard the cosmoshub-4 zone to Quicksilver",
  "connection_id": "connection-3",
  "base_denom": "uatom",
  "local_denom": "uqatom",
  "account_prefix": "cosmos",
  "multi_send": true,
  "liquidity_module": false,
  "deposit": "512000000uqck"
}
```

### update-zone

Submit a zone update proposal.

`quicksilverd update-zone [proposal-file]`

The proposal must include a deposit and the details must be provided as a json
file, e.g.

```
{
  "title": "Enable liquidity module for cosmoshub-4",
  "description": "Update cosmoshub-4 to enable liquidity module",
  "chain_id": "cosmoshub-4",
  "changes": [{
      "key": "liquidity_module",
      "value": "true",
  }],
  "deposit": "512000000uqck"
}
```

## Events

Events emitted by module for tracking messages and index transactions;

## Hooks

Description of hook functions that may be used by other modules to execute operations at specific points within this module;

## Queries

```
service Query {
  // RegisteredZoneInfos provides meta data on connected zones.
  rpc RegisteredZoneInfos(QueryRegisteredZonesInfoRequest)
      returns (QueryRegisteredZonesInfoResponse) {
    option (google.api.http).get = "/quicksilver/interchainstaking/v1/zones";
  }
  // DepositAccount provides data on the deposit address for a connected zone.
  rpc DepositAccount(QueryDepositAccountForChainRequest)
      returns (QueryDepositAccountForChainResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/deposit_address";
  }
  // DelegatorIntent provides data on the intent of the delegator for the given
  // zone.
  rpc DelegatorIntent(QueryDelegatorIntentRequest)
      returns (QueryDelegatorIntentResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/delegator_intent/"
        "{delegator_address}";
  }

  // Delegations provides data on the delegations for the given zone.
  rpc Delegations(QueryDelegationsRequest) returns (QueryDelegationsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/delegations";
  }

  // DelegatorDelegations provides data on the delegations from a given
  // delegator for the given zone.
  rpc DelegatorDelegations(QueryDelegatorDelegationsRequest)
      returns (QueryDelegatorDelegationsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/"
        "delegator_delegations/{delegator_address}";
  }

  // ValidatorDelegations provides data on the delegations to a given validator
  // for the given zone.
  rpc ValidatorDelegations(QueryValidatorDelegationsRequest)
      returns (QueryValidatorDelegationsResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/"
        "validator_delegations/{validator_address}";
  }

  // DelegationPlans provides data on the delegations to a given validator for
  // the given zone.
  rpc DelegationPlans(QueryDelegationPlansRequest)
      returns (QueryDelegationPlansResponse) {
    option (google.api.http).get =
        "/quicksilver/interchainstaking/v1/zones/{chain_id}/delegation_plans";
  }
}
```

### zones

Query registered zones.

`quicksilverd zones`

### intent

Query delegation intent for a given chain.

`quicksilverd intent [chain_id] [delegator_addr]`

### deposit-account

Query deposit account address for a given chain.

`quicksilverd deposit-account [chain_id]`

## Keepers

Keepers exposed by module;

## Parameters

Module parameters:

| Key                      | Type    | Example |
|:-------------------------|:--------|:--------|
| delegation_account_count | uint64  | 100     |
| delegation_account_split | uint64  | 10      |
| deposit_interval         | uint64  | 50      |
| delegate_interval        | uint64  | 100     |
| delegations_interval     | uint64  | 200     |
| validatorset_interval    | uint64  | 200     |
| commission_rate          | sdk.Dec | "0.02"  |

Description of parameters:

* `delegation_account_count` - the number of delegation accounts per zone;
* `delegation_account_split` - ;
* `deposit_interval` - ;
* `delegate_interval` - ;
* `delegations_interval` - ;
* `validatorset_interval` - ;
* `commission_rate` - ;

## Begin Block

Iterate through all registered zones and check validator set status. If the
status has changed, requery the validator set and update zone state.

## After Epoch End

