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

## Concepts

Key concepts, mechanisms and core logic for the module;

## State

### Registered Zone

A `Registered Zone` represents a Cosmos based blockchain (or zone) that
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

* **ConnectionId** - 
* **ChainId** - 
* **DepositAddress** - 
* **WithdrawalAddress** - 
* **PerformanceAddress** - 
* **DelegationAddresses** - 
* **AccountPrefix** - 
* **LocalDenom** - 
* **BaseDenom** - 
* **RedemptionRate** - 
* **LastRedemptionRate** - 
* **Validators** - 
* **AggregateIntent** - 
* **MultiSend** - 
* **LiquidityModule** - 
* **WithdrawalWaitgroup** - 
* **IbcNextValidatorHash** - 
* **ValidatorSelectionAllocation** - 
* **HoldingsAllocation** - 

### ICAAccount

```
type ICAAccount struct {
    Address          string
    Balance          sdk.Coins
    DelegatedBalance sdk.Coin
    PortName         string
    BalanceWaitgroup uint32
}
```

* **Address** - 
* **Balance** - 
* **DelegatedBalance** - 
* **PortName** - 
* **BalanceWaitgroup** - 

### Validator

```
type Validator struct {
	ValoperAddress  string
	CommissionRate  sdk.Dec
	DelegatorShares sdk.Dec
	VotingPower     sdk.Int
	Score           sdk.Dec
}
```

* **ValoperAddress** - 
* **CommissionRate** - 
* **DelegatorShares** - 
* **VotingPower** - 
* **Score** - 

### ValidatorIntent

```
type ValidatorIntent struct {
	ValoperAddress string
	Weight         sdk.Dec
}
```

* **ValoperAddress** - 
* **Weight** - 

## Messages

Description of message types that trigger state transitions;

## Transactions

Description of transactions that collect messages in specific contexts to trigger state transitions;

## Events

Events emitted by module for tracking messages and index transactions;

## Hooks

Description of hook functions that may be used by other modules to execute operations at specific points within this module;

## Queries

Description of available information request queries;

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
