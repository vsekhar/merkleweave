# Merkle Weave

A Merkle Weave is a verifiable data structure with improved write throughput.

Merkle trees are usually constrained by the tradeoff between ordering and write throughput.

An append-only Merkle tree provides total ordering at the cost of very low write throughput. A sparse Merkle tree improves write throughput by batching updates across lower levels of the tree, but cannot guarantee consistent ordering of updates across the entire tree.


