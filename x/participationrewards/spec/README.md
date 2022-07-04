# Participation Rewards

## Abstract

Module, x/participatiorewards, defines and implements the mechanisms to track,
allocate and distribute protocol participation rewards to users.

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
1. [End Block](#End-Block)
1. [After Epoch End](#After-Epoch-End)

## Concepts

The purpose of the participation rewards module is to reward users for protocol participation.

Specifically, we want to reward users for:

1. Staking and locking of QCK on the Quicksilver chain;
2. Positive validator selection, validators are ranked equally on performance and decentralization;
3. Holdings of off-chain assets (qAssets);

The total inflation allocation for participation rewards is divided
proportionally for each of the above according to the module [parameters](#Parameters).

### 1. Lockup Rewards

The staking and lockup rewards allocation is moved to the fee collector account
to be distributed by the staking module on the next begin blocker. Thus, the
**user rewards allocation** will be proportional to their stake of the overall
staked pool.

### 2. Validator Selection Rewards

Validators are ranked on two aspects with equal weighting, namely
decentralization and performance.

The **decentralilzation scores** are based on the normalized voting power of the
validators within a given zone, favouring smaller validators.

The **performance scores** are based on the validator rewards earned by a
special performance account that delegates an exact amount to each validator.
The total rewards earned by the performance account is divided by the number of
active validators to obtain the expected rewards. The performance score for
each validator is then simply the percentage of actual rewards compared to the
expected rewards (capped at 100%).

The overall **validator scores** are simply the multiple of their
decentralization score and their performance score.

Individual **users scores** are based on their validator selection intent
signalled at the previous epoch boundary. The user intent weights are
multiplied by the corresponding validator scores for the given zone and an
overall user score is calculated for the given zone along with an
**overall zone score**.

Each zone receives a **proportional rewards allocation** based on the total
value locked (TVL) for the zone relative to the TVL of all zones across the
protocol.

The overall zone score and the proportional rewards allocation determines the
amount of **tokens per point (TPP)** to be allocated for the given zone. Thus,
the **user rewards allocation** is the user's score multiplied by the TPP.

### 3. Holdings Rewards

Each zone receives a **proportional rewards allocation** based on the total
value locked (TVL) for the zone relative to the TVL of all zones across the
protocol.

Thus, the **user rewards allocation** is proportional to their holdings of
qAssets across all zones, capped at 2% per account.

## State

Participation Rewards maintains a `Score` for every `Validator` within a `Zone`.
`Score` is initially set to zero and is updated at the end of every epoch to
reflect the **overall score** for the validator (decntralization_score * performance_score).

## Messages

Description of message types that trigger state transitions;

## Transactions

Description of transactions that collect messages in specific contexts to trigger state transitions;

## Events

Events emitted by module for tracking messages and index transactions;

## Hooks

Description of hook functions that may be used by other modules to execute operations at specific points within this module;

## Queries

Participation Rewards module provides the below queries to check the module's state:

```
service Query {
  // Params returns the total set of participation rewards parameters.
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get =
        "/quicksilver/participationrewards/v1beta1/params";
  }
}
```

## Keepers

Keepers exposed by module;

## Parameters

Module parameters:

| Key                                                     | Type         | Example |
|:--                                                   ---|:--        ---|:--   ---|
| distribution_proportions.validator_selection_allocation | string (dec) | "0.34"  |
| distribution_proportions.holdings_allocation            | string (dec) | "0.33"  |
| distribution_proportions.lockup_allocation              | string (dec) | "0.33"  |

Description of parameters:

* `validator_selection_allocation` - the percentage of inflation rewards allocated to validator selection rewards;
* `holdings_allocation` - the percentage of inflation rewards allocated to qAssets hoildings rewards;
* `lockup_allocation` - the percentage of inflation rewards allocated to staking and locking of QCK;

## Begin Block

Description of logic executed with optional methods or external hooks;

## End Block

Description of logic executed with optional methods or external hooks;

## After Epoch End

The following is performed at the end of every epoch:

* Obtains the rewards allocations according to the module balances and
  distribution proportions parameters;
* Allocate zone rewards according to the proportional zone Total Value Locked
  (TVL) for both **Validator Selection** and **qAsset Holdings**;
* Calculate validator selection scores and allocations for every zone:
  1. Obtain performance account delegation rewards (`performanceScores`);
  2. Calculate decentralization scores (`distributionScores`);
  3. Calculate overall validator scores;
  4. Calculate user validator selection rewards;
  5. Distribute validator selection rewards;
* Calculate qAsset holdings:
  1. Obtain qAssets held by account (locally and off-chain via claims / Poof of
     Posession);
  2. Calculate user proportion (cap at 2%);
  3. Normalize and distribute allocation;
* Allocate lockup rewards by sending portion to `feeCollector` for distribution
  by Staking Module;
