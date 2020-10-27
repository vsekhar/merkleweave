# Merkle Weave

Merkle weave is a verifiable provable append-only blind log with high write throughput.

* Verifiable: anyone can audit the Merkle weave to verify that it has been operated correctly
* Provable: the log can produce a compact proof that an entry is included and the size of this proof grows slower than the size of the log the entry is in
* Append-only: once a value is written to the log, removing it would be apparent to auditors
* Blind: auditors do not need to know anything about client data to verify the log
* High write throughput: writes are not serialized

## Building up to the Merkle weave

A verifiable log allows users to commit log entries and auditors to verify that log entries have not be removed or changed by the operator.

### Hash chain

The simplest way to implement a verifiable log is as a _hash chain_:

```none
e_0 = Hash(d0)
e_1 = Hash(e_0, d1)
e_2 = Hash(e_1, d2)
...
```

You can audit a hash chain by reading it from start to finish. Given this auditability, the last hash in the chain `e_n` is enough to uniquely summarize a particular log with a particular sequence of entries.

```none
hc_n = e_n
```

A _proof_ is data that allows one to locally demonstrate that `e_k` is in the log identified by `hc_n`. In the case of a hash chain, the proof is simply the list of log entries from `e_k+1 = Hash(e_k, d_k+1)` to `e_n = Hash(e_n-1, d_n)`.

```none
proof(e_k, hc_n) = [e_k+1, ..., e_n]
```

For a suitably strong hash functions, it would be infeasible to produce this proof if `hc_n` had not been constructed from a chain of hashes that included `e_k`.

A client in possession of such a proof can hash `e_k` and the intervening entries contained in the proof of `e_k` to confirm that `e_n` is produced. Unfortunately the size of this proof (and the work required to produce and verify it) scales linearly (`O(N)`) with the size of the log. For a log with a billion entries, each proof requires roughly half a billion values and validating it requires hashing all those values. Thus proving inclusion in a hash chain is generally considered infeasible.

### Merkle Mountain Range

To improve on this, we instead place the entries inside a [Merkle Mountain Range](https://github.com/mimblewimble/grin/blob/master/doc/mmr.md) (MMR). With entries numbered from zero, the first 30 entries form a binary tree from the bottom up:

```none
              27
     11                26
   6   10       18           25
 2   5   9   14    17     21   24     30
0 1 3 4 7 8 12 13 15 16 19 20 22 23 28 29
```

Computing each entry works much like the hash chain, except in addition to hashing the previous entry, each new entry that is not a leaf also includes the hash of an earlier "left child" entry in the tree. The first 8 nodes would be computed as follows (indented to indicate parent/child relationships):

```none
e_0 = Hash(d_0)       --\
e_1 = Hash(e_0, d_1)  --|
                        e_2 = Hash(e_1, e_0, d_2) --\
                                                    |
e_3 = Hash(e_2, d_3) --\                            |
e_4 = Hash(e_3, d_4) --|                            |
                       e_5 = Hash(e_4, e_3, d_5)  --|
                                                    e_6 = Hash(e_5, e_2, d_6)

e_7 = Hash(e_6, d_7)
...
```

Each non-leaf node summarizes the nodes below it. The _peaks_ of an MMR of any given size are all nodes that have not yet been summarized by other higher-level nodes. For the 8-node MMR above, the peaks are `e_6` and `e_7`. Since an MMR grows deterministically, knowing the MMR's size is enough to tell you the indexes of its peaks. The MMR's size and the hashes of its peaks uniquely summarizes the MMR and all its entries. For the above 8 nodes:

```
mmr_8 = {size: 8, peaks: [e_6, e_7]}
```

The number of peaks scale `O(log_2(N))` with the size of the MMR, so the summary of an MMR is not as compact as that of a hash chain (which is just its last entry).

However the MMR providesmore compact proofs. For example, the proof that `e_3` is in `mmr_8` consists:

```none
proof(e_3, mmr_8) = [e_4, e_2]
```

For a suitably strong hash function, it would be infeasible to construct this proof if `mmr_8` had not been constructed from a tree of hashes that included `e_k`.

Verifying this proof involves hashing `e_3` with `e_4` (from the proof) to get `e_5`, then hashing `e_2` (from the proof) with `e_5` to get `e_6`, which is included in `mmr_8`.

For the hash chain with the same entries:

```none
proof(e_3, hc_8) = [e_4, e_5, e_6, e_7]
```

The size of proofs for `mmr_8` scale in `O(log_2(N))`, instead of `O(N)` for the hash chain. A proof that an entry exists in an MMR with a billion entires can be stored in less than 2 KB (30 peaks, 64 bytes hashes). Producing a summary requires `O(log_2(N))` reads however these are highly cacheable (once a node becomes a peak it remains so for some number of entries and all nodes are immutable).

Writing to an MMR requires looking up and hashing an additional entry: the left child. This `O(1)` penalty is incurred on half of all writes (non-leaf nodes).

Finally, writing to an MMR still requires strict sequencing, looking up and hashing the immediately prior entry. So write throughput remains within a constant factor of that of the hash chain.

### Sharding

A simple solution to improve write throughput would be to shard writes to a set of MMRs rather than only a single one. This effectively puts writes in separate causal chains. For a given entry `e`, we can only establish whether `e` came before or after a randomly-chosen subset of entries in the sharded MMR. This ordering is of limited usefulness.

To provide a more useful global order among entries in the sharded MMR we can set a timestamp on each entry. Notice, however, that generating a globally consistent and monotonically increasing timestamp in a distributed system is as difficult as strictly sequencing an MMR: the source of such a timestamp would become a single point of contention and limit write throughput as much as operating a single MMR does.

Even if we could generate such a timestamp, the independence of the MMRs means the operator can exploit timestamp skew between shards to potentially violate causality.

> **Attacking causality in a Sharded MMR**: consider an operator logging `e_4` on behalf of a client. Based on the hash prefix of `e_4`, the operator logs it to `mmr_6`.
>
> With knowledge of the timestamp of `e_4`, the operator (or an accomplice) can then generate `e_5`. Based on the hash prefix of `e_5`, the operator determines it is to be written to `mmr_2`. If the time stamp of the last entry in `mmr_2` is before the timestamp of `e_4`, then the operator can falsely timestamp `e_5` earlier than `e_4` and log it to `mmr_2`.
>
> Both `mmr_2` and `mmr_6` would independently verify by their hashes and their timestamp ordering, however causality would be violated _across_ the two MMRs because `e_5` in `mmr_2` incorrectly claims to have occurred before `e_4` in `mmr_6`.
>
> Including the timestamp in the hashes makes this worse because valid timestamp values are dense. The operator could generate many candidate values of `e_5 = Hash(d, s, ts)` by trying several values of `ts < e_4.ts` until they hit on one that directed `e_5` to an MMR with a last timestamp earlier than `ts` and the attack proceeds as above.
>
> Letting the client choose the sequence for their entry (rather than by hash prefix) doesn't help since the operator would just choose a sequence known to be behind.
>
> Logically this is solved by logging to all shards so that no shard is behind. Alternatively, we could log to `N/2+1` shards so that each log entry has at least one shard in common with the immediately previous entry. But either approach nullifies the throughput benefit of sharding in the first place.
>
> Anything less than a majority has the same effect to different proportions: a high probability of collision reduces the chance for operator malfeasance, but also reduces the throughput of the weave.
>
> Can we solve this by relying on the unspoofability of records at the application layer? If Alice and Bob timestamp a document and Bob is colluding with the operator, can Bob do Alice harm without being able to spoof Alice? Yes. Bob is the patent office. Alice submits a notarized invention application, Bob replies with a notarized receipt. Bob (or an accomplice) can then create an application copying Alice's invention and collude with the operator to have the copy timestamped earlier than Alice's. Bob can later reject Alice's application by pointing to the notarized copy created by Bob. The notarization of Alice's application and of Bob's receipt is thus meaningless.
>
> Caveat: this is a misuse of the infrastructure. Bob would have to create an application chain, independent of the notarization infrastructure, to have the patent office commit to a specific application ordering. Alternatively, Alice can notarize her application herself, and wait until the operator has committed to a HWM (see below) that is after than her application timestamp, storing the corresponding summary. Only then does she share her application with Bob. If Bob and the operator attempt to sneak in a backdated application into the log, Alice can prove operator malfeasance because the operator will be unable to prove inclusion of her summary in any future summary.

### Verifiable timestamping

The issue comes down to verifying that a set of independent (for throughput) timestampers are nevertheless providing timestamps (via commit waits) that are externally consistent. For the timestampers to have high throughput, they must operate independently. However independent operation means timestampers skew relative to each other. A malicious operator can use this skew to undetectably backdate an entry by finding (or waiting to happen upon) a suitably behind timestamper.

Countering this requires determining a global high water mark for timestamps. I.e. the infrastructure must commit to a monotonically increasing global minimum timestamp. Clients can wait until this global HWM has passed the timestamp of their entry before making their entry known to the operator. This ensures the operator cannot later produce an entry with a prior timestamp with knowledge of the client's entry.

This is equivalent to obtaining a summary of the entire system: the latest timestamp and hashes of all MMRs in the Merkle weave. The minimum timestamp among MMRs constitutes the HWM of the system. A client can wait until after the HWM of the system has passed the timestamp assigned to their transaction before acting on it. That is, the client can wait until they have collected a summary of each MMR that has as its latest timestamp a value greater than the timestamp of its entry. This is equivalent to the "merge window" of a more typical verifiable log.

MMRs can themselves generate entries every second to ensure they are never more than a second behind. They are like clock ticks. Clocks don't tick only when we look at them, they also tick when we don't.

Alternatively, we can write sentinel commitments to laggard MMRs only when a summary is requested. These sentinels do not vouch for any external data so they can have a data_hash of PREFIX000... Given that their data_hash is fixed, that hash also need not be stored. An empty data_hash is replaced with PREFIX000... As overall traffic increases, the need to do this lessens. So maybe we don't need to worry about saving storage by omitting data_hash in storage. We will need to decide between waiting for a sequence to update (if it's last value is relatively recent and therefore throughput is likely high) or writing a sentinel (if its last value is not recent and therefore throughput is likely low).

Paranoid clients can ask the operator to return a summary of the Merkle weave whose HWM is greater than the entry's timestamp. Building this summary is expensive, however with knowledge of the desired timestamp, some work can be cached/saved.

> **Partial summaries**: A lower cost alternative is to request a summary where a client-determined random set of prefixes is advanced past a certain timestamp. The more prefixes included in the request, the higher the level of certainty provided by the operator.
>
> Note that letting the operator choose prefixes doesn't work as they can just avoid the prefix chain on which they intend to front-run the client.

Once a summary of the Merkle weave is obtained, its inclusion in any future summary can also be proved. Summarized logging exists to earn client trust, but clients should eventually have enough trust in the operator that typical usage can involve regular writes.

How does this interact with uncertainty and commit waits across operators? Do they all gang together with this operation and produce summaries with suitably-high HWM? This can be done in parallel.

> **Cross logging**: previously considered logging each entry to some deterministic set of MMR shards. This doesn't help since it only partially achieves what we achieve with summaries above and reduces write throughput in the process. In the limit, writing to a majority of MMR shards guarantees the operator cannot backdate an entry (since it would have to conflict at at least one of the shards) but brings throughput back down to that of a single MMR.

## Open issues

* How to summarize?
  * Can't currently go from one summary to another so just have to trust the latest publication of the Merkle weave digest.
  * Summary must contain a timestamp to be reproducible?
  * Version number?
  * Vector clock...
  * Summary is just an array of summaries of all Merkle trees
    * There is no value in summarizing any further (even in hashing this array)
  * Going from one summary to another produces summary-to-summary proofs of all Merkle trees
    * Very expensive, but rare
* Storage:
  * Format: PREFIX:INDEX --> data_sha3256, node_sha3256, timestamp
  * Raw: 64 + 64 + 8 = 136 bytes
  * Protobuf: (1 byte per tag * 3 fields) + (1 byte length + 64 bytes) + (1 byte length + 64 bytes) + 1 byte length + (timestamp: (1 byte per tag * 2 fields) + 8 bytes + 4 bytes) = 148 bytes (8.8% overhead)
* Indexing
  * If we return the notarization hash to the user, they may want to look it up later by hash
  * To do so, we need to separately store an index of hashes to MMR index values
  * Alternatively, the notarization can consist of an index value and a hash
    * E.g. 148:ab8f3c....
  * This lets us skip the index and look it up directly in core storage via PREFIX:INDEX
* Caching
  * Values at PREFIX:index-->(data_hash, node_hash, timestamp)
    * Useful for assembling proofs, peaks are stable and likely to remain in cache based on LRU
  * Last value at PREFIX-->(index, node_hash, timestamp)
    * Check on each write, hint for next index number, node_hash for inclusion
    * Update with each write (the success of the write tells us it was the last value)
    * Can reconstruct by binary probing indexes in a prefix to find the latest, then read timestamp and cache
    * Useful to quickly build summaries for a given timestamp staleness, verifying that the timestamp of each entry from the cache is after the desired timestamp
