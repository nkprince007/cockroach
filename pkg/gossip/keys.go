// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package gossip

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/pkg/errors"
)

// separator is used to separate the non-prefix components of a
// Gossip key, facilitating the automated generation of regular
// expressions for various prefixes.
// It must not be contained in any of the other keys defined here.
const separator = ":"

// Constants for gossip keys.
const (
	// KeyClusterID is the unique UUID for this Cockroach cluster.
	// The value is a string UUID for the cluster.  The cluster ID is
	// gossiped by all nodes that contain a replica of the first range,
	// and it serves as a check for basic gossip connectivity. The
	// Gossip.Connected channel is closed when we see this key.
	KeyClusterID = "cluster-id"

	// KeyStorePrefix is the key prefix for gossiping stores in the network.
	// The suffix is a store ID and the value is roachpb.StoreDescriptor.
	KeyStorePrefix = "store"

	// KeyDeadReplicasPrefix is the key prefix for gossiping dead replicas in the
	// network. The suffix is a store ID and the value is
	// roachpb.StoreDeadReplicas.
	KeyDeadReplicasPrefix = "replica-dead"

	// KeyNodeIDPrefix is the key prefix for gossiping node id
	// addresses. The actual key is suffixed with the decimal
	// representation of the node id and the value is the host:port
	// string address of the node. E.g. node:1 => 127.0.0.1:24001
	KeyNodeIDPrefix = "node"

	// KeyNodeLivenessPrefix is the key prefix for gossiping node liveness info.
	KeyNodeLivenessPrefix = "liveness"

	// KeySentinel is a key for gossip which must not expire or
	// else the node considers itself partitioned and will retry with
	// bootstrap hosts.  The sentinel is gossiped by the node that holds
	// the range lease for the first range.
	KeySentinel = "sentinel"

	// KeyFirstRangeDescriptor is the descriptor for the "first"
	// range. The "first" range contains the meta1 key range, the first
	// level of the bi-level key addressing scheme. The value is a slice
	// of storage.Replica structs.
	KeyFirstRangeDescriptor = "first-range"

	// KeySystemConfig is the gossip key for the system DB span.
	// The value if a config.SystemConfig which holds all key/value
	// pairs in the system DB span.
	KeySystemConfig = "system-db"

	// KeyDistSQLNodeVersionKeyPrefix is key prefix for each node's DistSQL
	// version.
	KeyDistSQLNodeVersionKeyPrefix = "distsql-version"
)

// MakeKey creates a canonical key under which to gossip a piece of
// information. The first argument will typically be one of the key constants
// defined in this package.
func MakeKey(components ...string) string {
	return strings.Join(components, separator)
}

// MakePrefixPattern returns a regular expression pattern that
// matches precisely the Gossip keys created by invocations of
// MakeKey with multiple arguments for which the first argument
// is equal to the given prefix.
func MakePrefixPattern(prefix string) string {
	return regexp.QuoteMeta(prefix+separator) + ".*"
}

// MakeNodeIDKey returns the gossip key for node ID info.
func MakeNodeIDKey(nodeID roachpb.NodeID) string {
	return MakeKey(KeyNodeIDPrefix, nodeID.String())
}

// IsNodeIDKey returns true iff the provided key is a valid node ID key.
func IsNodeIDKey(key string) bool {
	return strings.HasPrefix(key, KeyNodeIDPrefix+separator)
}

// NodeIDFromKey attempts to extract a NodeID from the provided key.
// The key should have been constructed by MakeNodeIDKey.
// Returns an error if the key is not of the correct type or is not parsable.
func NodeIDFromKey(key string) (roachpb.NodeID, error) {
	trimmedKey := strings.TrimPrefix(key, KeyNodeIDPrefix+separator)
	if trimmedKey == key {
		return 0, errors.Errorf("%q is not a NodeID Key", key)
	}
	nodeID, err := strconv.ParseInt(trimmedKey, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed parsing NodeID from key %q", key)
	}
	return roachpb.NodeID(nodeID), nil
}

// MakeNodeLivenessKey returns the gossip key for node liveness info.
func MakeNodeLivenessKey(nodeID roachpb.NodeID) string {
	return MakeKey(KeyNodeLivenessPrefix, nodeID.String())
}

// MakeStoreKey returns the gossip key for the given store.
func MakeStoreKey(storeID roachpb.StoreID) string {
	return MakeKey(KeyStorePrefix, storeID.String())
}

// MakeDeadReplicasKey returns the dead replicas gossip key for the given store.
func MakeDeadReplicasKey(storeID roachpb.StoreID) string {
	return MakeKey(KeyDeadReplicasPrefix, storeID.String())
}

// MakeDistSQLNodeVersionKey returns the gossip key for the given store.
func MakeDistSQLNodeVersionKey(nodeID roachpb.NodeID) string {
	return MakeKey(KeyDistSQLNodeVersionKeyPrefix, nodeID.String())
}
