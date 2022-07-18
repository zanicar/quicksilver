# Interchain Staking

## Abstract

Module, `x/interchainstaking`, defines and implements the core Quicksilver
protocol. _(wip)_

## Contents

1. [Concepts](#Concepts)
1. [State](#State)
1. [Messages](#Messages)
1. [Transactions](#Transactions)
1. [Proposals](#Proposals)
1. [Events](#Events)
1. [Hooks](#Hooks)
1. [Queries](#Queries)
1. [Keepers](#Keepers)
1. [Parameters](#Parameters)
1. [Begin Block](#Begin-Block)
1. [After Epoch End](#After-Epoch-End)
1. [IBC](#IBC)

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
* **PerformanceAddress** - remote zone performance address (each validator gets
  an exact equal delegation from this account to measure performance);
* **DelegationAddresses** - remote zone delegation addresses to represent
  granular voting power;
* **AccountPrefix** - remote zone account address prefix;
* **LocalDenom** - protocol denomination (qAsset), e.g. uqatom;
* **BaseDenom** - remote zone denomination (uStake), e.g. uatom;
* **RedemptionRate** - redemption rate between protocol qAsset and native
  remote asset;
* **LastRedemptionRate** - redemption rate as at previous epoch boundary
  (used to prevent epoch boundary gaming);
* **Validators** - list of validators on the remote zone;
* **AggregateIntent** - the aggregated delegation intent of the protocol for
  this remote zone. The map key is the corresponding validator address
  contained in the `ValidatorIntent`;
* **MultiSend** - multisend support on remote zone;
* **LiquidityModule** - liquidity module enabled on remote zone;
* **WithdrawalWaitgroup** - tally of pending withdrawal transactions;
* **IbcNextValidatorHash** - 
* **ValidatorSelectionAllocation** - proportional zone rewards allocation for
  validator selection;
* **HoldingsAllocation** - proportional zone rewards allocation for asset
  holdings;

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
* **BalanceWaitgroup** - the tally of pending balance query transactions sent
  to the remote zone;

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

* **ValoperAddress** - the validator address;
* **CommissionRate** - the validator commission rate;
* **DelegatorShares** - 
* **VotingPower** - the validator voting power on the remote zone;
* **Score** - the validator Quicksilver protocol overall score;

### ValidatorIntent

`ValidatorIntent` represents the weighted delegation intent to a particular
validator.

```
type ValidatorIntent struct {
	ValoperAddress string
	Weight         sdk.Dec
}
```

* **ValoperAddress** - the remote zone validator address;
* **Weight** - the weight of intended delegation to this validator;

### DelegatorIntent

`DelegatorIntent` represents the delegation intent for every individual
`RegisteredZone.DelegationAddresses`. The overall delegations must match the
`RegisteredZone.AggregateIntent` for each validator in the zone as closely as
possible. The protocol spreads and balances delegations across delegation
accounts for efficiency purposes.

```
type DelegatorIntent struct {
	Delegator string
	Intents   []*ValidatorIntent
}
```

* **Delegator** - the delegation account address on the remote zone;
* **Intents** - the delegation intents to individual validators on the remote
  zone;

### Delegation

`Delegation` represents the actual delegations made by
`RegisteredZone.DelegationAddresses` to validators on the remote zone;

```
type Delegation struct {
	DelegationAddress string
	ValidatorAddress  string
	Amount            github_com_cosmos_cosmos_sdk_types.Coin
	Height            int64
	RedelegationEnd   int64
}
```

* **DelegationAddress** - the delegator address on the remote zone;
* **ValidatorAddress** - the validator address on the remote zone;
* **Amount** - the amount delegated;
* **Height** - the block height at which the delegation occured;
* **RedelegationEnd** - ;

## Messages

### MsgRequestRedemption

Redeems the indicated qAsset coin amount from the protocol, converting the
qAsset back to the native asset at the appropriate redemption rate.

```
type MsgRequestRedemption struct {
	Coin               string
	DestinationAddress string
	FromAddress        string
}
```

* **Coin** - qAsset as standard cosmos sdk cli coin string, {amount}{denomination};
* **DestinationAddress** - standard cosmos sdk bech32 address string;
* **FromAddress** - standard cosmos sdk bech32 address string;

**Transaction**: [`redeem`](#redeem)

### MsgSignalIntent

Signal validator delegation intent for a given zone by weight.

```
type MsgSignalIntent struct {
	ChainId     string
	Intents     []*ValidatorIntent
	FromAddress string
}
```

* **ChainId** - zone identifier string;
* **Intents** - list of validator intents according to weight;
* **FromAddress** - standard cosmos sdk bech32 address string;

**Transaction**: [`signal-intent`](#signal-intent)

## Transactions

### signal-intent

Signal validator delegation intent by providing a comma seperated string
containing a decimal weight and the bech32 validator address.

`quicksilverd signal-intent [chain_id] [delegation_intent]`

Example:

`quicksilverd signal-intent cosmoshub-4 0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz`

### redeem

Redeem qAssets for native tokens.

`quicksilverd redeem [coins] [destination_address]`

Example:

`quicksilverd redeem 2500000uatom cosmos1pgfzn0zhxjjgte7hprwtnqyhrn534lqk437x2w`

## Proposals

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

`quicksilverd query interchainstaking zones`

Example response:

```
pagination:
  next_key: null
  total: "2"
zones:
- account_prefix: cosmos
  aggregate_intent: {}
  base_denom: uatom
  chain_id: lstest-1
  connection_id: connection-0
  delegation_addresses:
  - address: cosmos12hww50r7q7xyhspt72c9c8n3uyknqhv208sxuq9mcqdqjv0mcreq62maa2
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "25083333"
      denom: uatom
    port_name: icacontroller-lstest-1.delegate.9
  ...
  deposit_address:
    address: cosmos146xjrj2tass9fvtcw30dtl9v8f4t26z7cjxlxw0paxyyxmx2hqcq73vk6p
    balance:
    - amount: "25000000"
      denom: cosmosvaloper16pxh2v4hr28h2gkntgfk8qgh47pfmjfhvcamkc3
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uatom
    port_name: icacontroller-lstest-1.deposit
  holdings_allocation: []
  ibc_next_validators_hash: Qn4t+8M6bod6ewSYwnPScdWwwbSE7mc47GlMpuo15d0=
  last_redemption_rate: "1.000000000000000000"
  liquidity_module: true
  local_denom: uqatom
  multi_send: true
  performance_address:
    address: cosmos1yp64sfc5d4g4xtemptachyd2jaraxz8c5vptp7swgnv86l3ll3yqzz72wk
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uatom
    port_name: icacontroller-lstest-1.performance
  redemption_rate: "1.000000000000000000"
  validator_selection_allocation: []
  validators:
  - commission_rate: "0.030000000000000000"
    delegator_shares: "4000093333.000000000000000000"
    score: "0.000000000000000000"
    valoper_address: cosmosvaloper12evgzwsc2av7nfc5x7p74g9ppmfwm30xug6pwv
    voting_power: "4000093333"
  ...
  withdrawal_address:
    address: cosmos1w7x78xu4ms3qwspryl8jjy57l3esns8ayh6mj9g3544wmgnfnzrs86lr9p
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uatom
    port_name: icacontroller-lstest-1.withdrawal
  withdrawal_waitgroup: 12
- account_prefix: osmo
  aggregate_intent: {}
  base_denom: uosmo
  chain_id: lstest-2
  connection_id: connection-1
  delegation_addresses: []
  deposit_address:
    address: osmo14s68pery7n8s9cm6lzwxv4s0ppucctv28fcmtqg852965hfgpuvsmq5edm
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uosmo
    port_name: icacontroller-lstest-2.deposit
  holdings_allocation: []
  ibc_next_validators_hash: TyRByZjTIrfQ81mMlvoRyg1crPx4Kk9Lur+Kkty06h8=
  last_redemption_rate: "1.000000000000000000"
  liquidity_module: true
  local_denom: uqosmo
  multi_send: true
  performance_address:
    address: osmo1v799s4plwuyux8xunmzxcw6y2g8t5u373ravsy2zxg7k0x8g7pdsaa8ve9
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uosmo
    port_name: icacontroller-lstest-2.performance
  redemption_rate: "1.000000000000000000"
  validator_selection_allocation: []
  validators: []
  withdrawal_address:
    address: osmo1xcs9r3ssmjndgr09jww29cs60ygck8g3udyl6savg3nkercfhl2qtp3lwv
    balance: []
    balance_waitgroup: 0
    delegated_balance:
      amount: "0"
      denom: uosmo
    port_name: icacontroller-lstest-2.withdrawal
  withdrawal_waitgroup: 0
```

### intent

Query delegation intent for a given chain.

`quicksilverd query interchainstaking intent [chain_id] [delegator_addr]`

### deposit-account

Query deposit account address for a given chain.

`quicksilverd query interchainstaking deposit-account [chain_id]`

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

The following is performed at the end of every epoch for each registered zone:

* Aggregate Intents:
  1. Iterate through all stored instances of `DelegatorIntent` for each zone
     and obtain the **delegator account balance**;
  2. Compute the **base balance** using the account balance and `RedpemtionRate`;
  3. Ordinalize the delegator's validator intents by `Weight`;
  4. Set the zone `AggregateIntent` and update zone state;
* Query delegator delegations for each zone and update delegation records:
  1. Query delegator delegations `cosmos.staking.v1beta1.Query/DelegatorDelegations`;
  2. For each response (per delegator `DelegationsCallback`), verify every
     delegation record (via IBC `DelegationCallback`) and update delegation
     record accordingly (add, update or remove);
  3. Update validator set;
  4. Update zone;
* Withdraw delegation rewards for each zone and distribute:
  1. Query delegator rewards `cosmos.distribution.v1beta1.Query/DelegationTotalRewards`;
  2. For each response (per delegator `RewardsCallback`), send withdrawal
     messages for each of its validator delegations and add tally to
     `WithdrawalWaitgroup`;
  3. For each IBC acknowledgement decrement the `WithdrawalWaitgroup`. Once
     all responses are collected (`WithdrawalWaitgroup == 0`) query the balance
     of `WithdrawalAddress` (`cosmos.bank.v1beta1.Query/AllBalances`), then
     distribute rewards (`DistributeRewardsFromWithdrawAccount`).

     This approach ensures the exact rewards amount is known at the time of
     distribution.

## IBC

### Messages, Acknowledgements & Handlers

#### MsgWithdrawDelegatorReward

Triggered at the end of every epoch if delegator accounts have accrued rewards.
Collects rewards to zone withdrawal account `WithdrawalAddress` and distributes
rewards once all delegator rewards withdrawals have been acknowledged.

* **Endpoint:** `/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward`
* **Handler:** `HandleWithdrawRewards`

#### MsgRedeemTokensforShares

Triggered during execution of `Delegate` for delegation allocations that are
not in the zone `BaseDenom`. During callback the relevant delegation record is
updated.

* **Endpoint:** `/cosmos.staking.v1beta1.MsgRedeemTokensforShares`
* **Handler:** `HandleRedeemTokens`

#### MsgTokenizeShares

Triggered by `RequestRedemption` when a user redeems qAssets. Withdrawal
records are set or updated accordingly.  
See [MsgRequestRedemption](#MsgRequestRedemption).

* **Endpoint:** `/cosmos.staking.v1beta1.MsgTokenizeShares`
* **Handler:** `HandleTokenizedShares`

#### MsgDelegate

Triggered by `Delegate` whenever delagtions are made by the protocol to zone
validators. `HandleDelegate` distinguishes `DelegationAddresses` and updates
delegation records for these delegation accounts.

* **Endpoint:** `/cosmos.staking.v1beta1.MsgDelegate`
* **Handler:** `HandleDelegate`

#### MsgBeginRedelegate

Not implemented.

* **Endpoint:** `/cosmos.staking.v1beta1.MsgBeginRedelegate`
* **Handler:** `HandleBeginRedelegate`

#### MsgSend

Triggered by `TransferToDelegate` during `HandleReceiptTransaction`.  
See [Deposit Interval](#Deposit-Interval).

* **Endpoint:** `/cosmos.bank.v1beta1.MsgSend`
* **Handler:** `HandleCompleteSend`

`HandleCompleteSend` executes one of the following options based on the
`FromAddress` and `ToAddress` of the msg:

1. **Delegate rewards accoring to global intents.**  
   (If `FromAddress` is the zone's `WithdrawalAddress`);
2. **Withdraw native assets for user.**  
   (If `FromAddress` is one of zone's `DelegationAddresses`);
3. **Delegate amount according to delegation plan.**  
   (If `FromAddress` is `DepositAddress` and `ToAddress` is one of zone's `DelegationAddresses`);

#### MsgMultiSend

Not sent?

* **Endpoint:** `/cosmos.bank.v1beta1.MsgMultiSend`
* **Handler:** `HandleCompleteMultiSend`

#### MsgSetWithdrawAddress

Triggered during zone initialization for every `DelegationAddresses` and
for the `PerformanceAddress`. The purpose of using a dedicated withdrawal
account allows for accurate rewards withdrawal accounting, that would otherwise
be impossible as the rewards amount will only be known at the time the msg is
triggered, and not at the time it was executed by the remote zone (due to network
latency and different zone block times, etc).

* **Endpoint:** `/cosmos.distribution.v1beta1.MsgSetWithdrawAddress`
* **Handler:** `HandleUpdatedWithdrawAddress`

#### MsgTransfer

Triggered by `DistributeRewardsFromWithdrawAccount` to distribute rewards
across the zone delegation accounts and collect fees to the module fee account.
The `RedemptionRate` is updated accordingly.  
See [WithdrawalAddress Balances](#WithdrawalAddress-Balances).

* **Endpoint:** `/ibc.applications.transfer.v1.MsgTransfer`
* **Handler:** `HandleMsgTransfer`

### Queries, Requests & Callbacks

This module registeres the following queries, requests and callbacks.

#### DepositAddress Balances

For every registered zone a periodic `AllBalances` query is run against the
`DepositAddress`. The query is proven by utilizing provable KV queries that
update the individual account balances `AccountBalanceCallback`, trigger the
`depositInterval` and finally update the zone state.

* **Query:** `cosmos.bank.v1beta1.Query/AllBalances`
* **Callback:** `AllBalancesCallback`

#### Delegator Delegations

Query delegator delegations for each zone and update delegation records.  
See [After Epoch End](#After-Epoch-End).

* **Query:** `cosmos.staking.v1beta1.Query/DelegatorDelegations`
* **Callback:** `DelegationsCallback`

#### Delegate Total Rewards

Withdraw delegation rewards for each zone and distribute.  
See [After Epoch End](#After-Epoch-End).

* **Query:** `cosmos.distribution.v1beta1.Query/DelegationTotalRewards`
* **Callback:** `RewardsCallback`

#### WithdrawalAddress Balances

Triggered by `HandleWithdrawRewards`.  
See [MsgWithdrawDelegatorReward](#MsgWithdrawDelegatorReward).

* **Query:** `cosmos.bank.v1beta1.Query/AllBalances`
* **Callback:** `DistributeRewardsFromWithdrawAccount`

#### Deposit Interval

Monitors transaction events of the zone `DepositAddress` on the remote chain
for receipt transactions that are then handled by `HandleReceiptTransaction`.
On valid receipts the delegation intent is updated (`UpdateIntent`) and new
qAssets minted and transferred to the sender (`MintQAsset`). A delegation
plan is computed (`DeterminePlanForDelegation`) and then executed
(`TransferToDelegate`). Successfully executed receipts are recorded to state.

* **Query:** `cosmos.tx.v1beta1.Service/GetTxsEvent`
* **Callback:** `DepositIntervalCallback`

#### Performance Balance Query

Triggered at zone registration when the zone performance account
`PerformanceAddress` is created. It monitors the performance account balance
until sufficient funds are available to execute the performance delegations.  
See [x/participationrewards/spec](../../participationrewards/spec/README.md).

* **Query:** `cosmos.bank.v1beta1.Query/AllBalances`
* **Callback:** `PerfBalanceCallback`

#### Validator Set Query

An essential query to ensure that the registred zone state accurately reflects
the validator set of the remote zone for bonded, unbonded and unbonding
validators.

* **Query:** `cosmos.staking.v1beta1.Query/Validators`
* **Callback:** `ValsetCallback`
