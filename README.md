# Merkle Weave

A Merkle weave is a verifiable data structure with improved write throughput.

Merkle trees are usually constrained by the tradeoff between ordering and write throughput.

An append-only Merkle tree provides total ordering at the cost of very low write throughput. A sparse Merkle tree improves write throughput by batching updates across lower levels of the tree, but cannot guarantee consistent ordering of updates across the entire tree.

Instead, a Merkle weave shards writes across a number of underlying Merkle trees to improve performance. New data is atomically written to a fixed subset of Merkle trees in a deterministic and verifiable way. As a result, ordering of writes is maintained and write throughput can horizontally scale.

## Performance

A Merkle weave is parameterized by the total number of underlying trees `N`, the number of trees touched by each write `k` and the total writes per second to the Merkle weave `r`.

TODO: probability of collision based on N choose k?

The write throughput capacity of a Merkle weave scales with `N/k`. So why not have large `N` and small (min. 2) `k`?

Each write will touch each tree with probability `k/N`. The average tree head will be `N/2k` writes old, or `N/2kr` seconds old.

We want to maximize throughput `N/k` while minimizing tree head age `N/2kr`.

From the top down, we can assume a write rate of 10 QPS for the overall Merkle weave. For 256 trees (1-byte prefix), and 2 trees per write, the average tree head will be 6.4 seconds old.

| Trees (N) | Trees per write (k) | Throughput capacity | Writes per second (r) | Head age |
|-----------|---------------------|---------------------|-----------------------|----------|
|   256     |           2         |                     |           10          |    6.4s  |

From the bottom up, we can assume a write capacity of 1000 QPS per transaction (Vitess docs somewhere).

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
