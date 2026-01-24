package algorithms

import (
	"crypto/sha256"
	"time"
)

// Consensus algorithms - Agreement protocols for distributed systems
// Used by: All blockchains for block finality and state agreement

// ValidatorID uniquely identifies a validator
type ValidatorID string

// Validator represents a consensus participant
type Validator struct {
	ID     ValidatorID
	Stake  uint64 // For PoS weight
	PubKey []byte
}

// Block represents a simplified block for consensus
type Block struct {
	Height    uint64
	Hash      [32]byte
	ParentHash [32]byte
	Proposer  ValidatorID
	Timestamp time.Time
	Txs       [][]byte
}

// Vote represents a validator's vote on a block
type Vote struct {
	ValidatorID ValidatorID
	BlockHash   [32]byte
	Round       uint64
	Signature   []byte
}

// ------------------------------------------------------------
// PBFT (Practical Byzantine Fault Tolerance)
// Tolerates f Byzantine faults with 3f+1 total nodes
// Used by: Hyperledger Fabric, some permissioned chains
// ------------------------------------------------------------

// PBFTState represents the state of a PBFT node
type PBFTState int

const (
	PBFTIdle PBFTState = iota
	PBFTPrePrepare
	PBFTPrepare
	PBFTCommit
)

// PBFTMessage types
type PBFTMessageType int

const (
	PrePrepare PBFTMessageType = iota
	Prepare
	Commit
	ViewChange
	NewView
)

// PBFTMessage represents a PBFT protocol message
type PBFTMessage struct {
	Type      PBFTMessageType
	View      uint64
	Sequence  uint64
	Digest    [32]byte
	NodeID    ValidatorID
	Signature []byte
}

// PBFTNode simulates a PBFT consensus node
type PBFTNode struct {
	ID         ValidatorID
	View       uint64
	Sequence   uint64
	State      PBFTState
	Validators []Validator

	// Message logs
	PrePrepareLog map[uint64]*PBFTMessage
	PrepareLog    map[uint64][]*PBFTMessage
	CommitLog     map[uint64][]*PBFTMessage
}

// NewPBFTNode creates a new PBFT node
func NewPBFTNode(id ValidatorID, validators []Validator) *PBFTNode {
	// TODO: Implement
	return nil
}

// IsPrimary returns true if this node is the current primary/leader
func (n *PBFTNode) IsPrimary() bool {
	// TODO: Implement
	// Primary = validators[view % len(validators)]
	return false
}

// OnPrePrepare handles a pre-prepare message from the primary
func (n *PBFTNode) OnPrePrepare(msg *PBFTMessage, block *Block) error {
	// TODO: Implement
	// 1. Verify sender is primary for current view
	// 2. Verify we haven't accepted a different pre-prepare for this sequence
	// 3. Verify block digest matches
	// 4. Broadcast Prepare message to all validators
	return nil
}

// OnPrepare handles a prepare message
func (n *PBFTNode) OnPrepare(msg *PBFTMessage) error {
	// TODO: Implement
	// 1. Add to prepare log
	// 2. If we have 2f+1 matching prepares (including our own), move to Commit
	// 3. Broadcast Commit message
	return nil
}

// OnCommit handles a commit message
func (n *PBFTNode) OnCommit(msg *PBFTMessage) error {
	// TODO: Implement
	// 1. Add to commit log
	// 2. If we have 2f+1 matching commits, block is finalized
	// 3. Execute block and update state
	return nil
}

// HasQuorum checks if we have enough matching messages
// quorumSize = 2f+1 where f = (n-1)/3
func (n *PBFTNode) HasQuorum(count int) bool {
	// TODO: Implement
	return false
}

// StartViewChange initiates a view change (leader election)
func (n *PBFTNode) StartViewChange() {
	// TODO: Implement
	// 1. Increment view
	// 2. Broadcast ViewChange message
	// 3. Wait for 2f+1 ViewChange messages
	// 4. New primary sends NewView message
}

// ------------------------------------------------------------
// Tendermint-style Consensus
// Round-based BFT with propose -> prevote -> precommit
// Used by: Cosmos, Binance Chain, Terra
// ------------------------------------------------------------

// TendermintRoundState represents state within a round
type TendermintRoundState int

const (
	TendermintPropose TendermintRoundState = iota
	TendermintPrevote
	TendermintPrecommit
	TendermintCommit
)

// TendermintNode simulates a Tendermint consensus node
type TendermintNode struct {
	ID         ValidatorID
	Height     uint64
	Round      uint64
	Step       TendermintRoundState
	Validators []Validator

	LockedBlock *Block  // Block we're locked on
	LockedRound int64   // Round we locked (-1 if not locked)
	ValidBlock  *Block  // Valid block seen
	ValidRound  int64   // Round valid block was seen

	Prevotes   map[uint64]map[ValidatorID]*Vote
	Precommits map[uint64]map[ValidatorID]*Vote
}

// NewTendermintNode creates a new Tendermint node
func NewTendermintNode(id ValidatorID, validators []Validator) *TendermintNode {
	// TODO: Implement
	return nil
}

// GetProposer returns the proposer for a given height and round
func (n *TendermintNode) GetProposer(height, round uint64) ValidatorID {
	// TODO: Implement
	// Deterministic proposer selection (often weighted by stake)
	return ""
}

// OnProposal handles a block proposal from the proposer
func (n *TendermintNode) OnProposal(block *Block, round uint64) error {
	// TODO: Implement
	// 1. Verify proposer is correct for this height/round
	// 2. Verify block validity
	// 3. If locked, only accept if block matches or round > lockedRound
	// 4. Broadcast Prevote
	return nil
}

// OnPrevote handles a prevote message
func (n *TendermintNode) OnPrevote(vote *Vote) error {
	// TODO: Implement
	// 1. Add to prevotes
	// 2. If 2/3+ prevotes for a block, broadcast Precommit for that block
	// 3. If 2/3+ prevotes for nil, broadcast Precommit nil
	return nil
}

// OnPrecommit handles a precommit message
func (n *TendermintNode) OnPrecommit(vote *Vote) error {
	// TODO: Implement
	// 1. Add to precommits
	// 2. If 2/3+ precommits for a block, commit the block
	// 3. If 2/3+ precommits for nil or timeout, move to next round
	return nil
}

// HasTwoThirdsMajority checks for 2/3+ votes
func (n *TendermintNode) HasTwoThirdsMajority(votes map[ValidatorID]*Vote, blockHash [32]byte) bool {
	// TODO: Implement
	// For PoS: sum stake of matching votes, check if > 2/3 total stake
	return false
}

// ------------------------------------------------------------
// Simple Proof of Stake Primitives
// Building blocks for PoS consensus
// ------------------------------------------------------------

// ValidatorSet manages the active validator set
type ValidatorSet struct {
	Validators  []Validator
	TotalStake  uint64
}

// NewValidatorSet creates a validator set from a list of validators
func NewValidatorSet(validators []Validator) *ValidatorSet {
	// TODO: Implement
	return nil
}

// GetByID returns a validator by ID
func (vs *ValidatorSet) GetByID(id ValidatorID) *Validator {
	// TODO: Implement
	return nil
}

// SelectProposer selects a proposer weighted by stake
// seed: random seed for deterministic selection (e.g., previous block hash)
func (vs *ValidatorSet) SelectProposer(seed []byte) ValidatorID {
	// TODO: Implement
	// Weighted random selection based on stake
	// Hint: Use seed to create deterministic randomness
	return ""
}

// CalculateVotingPower calculates voting power for a subset of validators
func (vs *ValidatorSet) CalculateVotingPower(voterIDs []ValidatorID) uint64 {
	// TODO: Implement
	return 0
}

// HasTwoThirdsStake checks if voters have > 2/3 of total stake
func (vs *ValidatorSet) HasTwoThirdsStake(voterIDs []ValidatorID) bool {
	// TODO: Implement
	return false
}

// ------------------------------------------------------------
// Finality Gadget (simplified Casper FFG style)
// Provides finality on top of a fork-choice rule
// Used by: Ethereum 2.0 (Casper FFG)
// ------------------------------------------------------------

// Checkpoint represents a finality checkpoint
type Checkpoint struct {
	Epoch      uint64
	BlockHash  [32]byte
	Justified  bool
	Finalized  bool
}

// FinalityGadget tracks justification and finalization
type FinalityGadget struct {
	Checkpoints    map[uint64]*Checkpoint // epoch -> checkpoint
	Validators     *ValidatorSet
	LastJustified  uint64
	LastFinalized  uint64
}

// NewFinalityGadget creates a finality gadget
func NewFinalityGadget(validators *ValidatorSet) *FinalityGadget {
	// TODO: Implement
	return nil
}

// AddCheckpoint adds a checkpoint for an epoch
func (fg *FinalityGadget) AddCheckpoint(epoch uint64, blockHash [32]byte) {
	// TODO: Implement
}

// OnVote processes a finality vote (attestation)
// source: the checkpoint being justified from
// target: the checkpoint being justified
func (fg *FinalityGadget) OnVote(voterID ValidatorID, source, target uint64) error {
	// TODO: Implement
	// 1. Verify source is justified
	// 2. Verify target > source
	// 3. Check for slashable conditions (double vote, surround vote)
	// 4. If 2/3+ votes for source->target, justify target
	// 5. If target is justified and target.epoch == source.epoch + 1, finalize source
	return nil
}

// IsFinalized checks if a block at given epoch is finalized
func (fg *FinalityGadget) IsFinalized(epoch uint64) bool {
	// TODO: Implement
	return false
}

// ------------------------------------------------------------
// Utility Functions
// ------------------------------------------------------------

// HashBlock computes a block hash
func HashBlock(block *Block) [32]byte {
	// TODO: Implement
	// Serialize block fields and hash with SHA256
	return sha256.Sum256(nil)
}

// VerifySignature verifies a validator's signature
func VerifySignature(pubKey []byte, message []byte, signature []byte) bool {
	// TODO: Implement
	// In practice, use ed25519 or BLS signatures
	return false
}

// Consensus concepts to understand:
//
// 1. Safety vs Liveness
//    - Safety: Never finalize conflicting blocks
//    - Liveness: Eventually make progress
//    - CAP theorem trade-offs
//
// 2. Byzantine Fault Tolerance
//    - Tolerate f faults with 3f+1 nodes
//    - Why 3f+1? Need 2f+1 honest to form quorum, f could be silent
//
// 3. Finality
//    - Probabilistic (Bitcoin): More confirmations = more final
//    - Absolute (PBFT, Tendermint): Once committed, never reverts
//    - Economic (Casper): Revert costs slashed stake
//
// 4. Resources:
//    - PBFT paper: https://pmg.csail.mit.edu/papers/osdi99.pdf
//    - Tendermint: https://arxiv.org/abs/1807.04938
//    - Casper FFG: https://arxiv.org/abs/1710.09437
//    - HotStuff: https://arxiv.org/abs/1803.05069
