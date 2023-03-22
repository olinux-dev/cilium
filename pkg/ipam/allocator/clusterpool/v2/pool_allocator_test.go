// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package v2

import (
	"fmt"
	"math/big"
	"net/netip"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cilium/cilium/pkg/ipam/allocator/clusterpool/cidralloc"
	ipamTypes "github.com/cilium/cilium/pkg/ipam/types"
	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
)

func newAllocators(maskSize int, cidrs ...string) []cidralloc.CIDRAllocator {
	isIPv6 := false
	for _, c := range cidrs {
		isIPv6 = isIPv6 || strings.Contains(c, ":")
	}
	cidrsets, err := cidralloc.NewCIDRSets(isIPv6, cidrs, maskSize)
	if err != nil {
		panic(err)
	}
	return cidrsets
}

func newAllocatorsWithOccupied(maskSize int, cidrs []string, occupied []string) []cidralloc.CIDRAllocator {
	cidrsets := newAllocators(maskSize, cidrs...)
	for _, c := range occupied {
		err := occupyCIDR(cidrsets, netip.MustParsePrefix(c))
		if err != nil {
			panic(err)
		}
	}
	return cidrsets
}

func TestPoolAllocator_AllocateToNode(t *testing.T) {
	type fields struct {
		pools map[string]cidrPool
		nodes map[string]poolToCIDRs
		ready bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    *v2.CiliumNode
		want    *v2.CiliumNode
		wantErr string
	}{
		{
			name: "empty pool",
			fields: fields{
				pools: map[string]cidrPool{
					"default": {
						v4: nil,
						v6: nil,
					},
				},
				ready: true,
			},
			args: &v2.CiliumNode{
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
										IPv6Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
			wantErr: errPoolEmpty.Error(),
		},
		{
			name: "empty ipv6 pool, but ipv4 works",
			fields: fields{
				pools: map[string]cidrPool{
					"default": {
						v4: newAllocators(24, "192.168.0.0/16"),
						v6: nil,
					},
				},
				nodes: map[string]poolToCIDRs{},
				ready: true,
			},
			args: &v2.CiliumNode{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node1",
				},
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
										IPv6Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
			want: &v2.CiliumNode{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node1",
				},
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Allocated: []ipamTypes.IPAMPoolAllocation{
								{
									Pool:  "default",
									CIDRs: []string{"192.168.0.0/24"},
								},
							},
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
										IPv6Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
			wantErr: fmt.Errorf(`failed to allocate ipv6 address for node "node1" from pool "default": %w`, errPoolEmpty).Error(),
		},
		{
			name: "basic default pool",
			fields: fields{
				pools: map[string]cidrPool{
					"default": {
						v4: newAllocators(24, "192.168.0.0/16"),
						v6: newAllocators(96, "f00d::/80"),
					},
				},
				nodes: map[string]poolToCIDRs{},
				ready: true,
			},
			args: &v2.CiliumNode{
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
										IPv6Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
			want: &v2.CiliumNode{
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Allocated: []ipamTypes.IPAMPoolAllocation{
								{
									Pool:  "default",
									CIDRs: []string{"192.168.0.0/24", "f00d::/96"},
								},
							},
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
										IPv6Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "occupy ipv4, allocate ipv6",
			fields: fields{
				pools: map[string]cidrPool{
					"default": {
						v4: newAllocators(24, "192.168.0.0/16"),
						v6: newAllocators(96, "f00d::/80"),
					},
				},
				nodes: map[string]poolToCIDRs{},
				ready: true,
			},
			args: &v2.CiliumNode{
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Allocated: []ipamTypes.IPAMPoolAllocation{
								{
									Pool:  "default",
									CIDRs: []string{"192.168.0.0/24"},
								},
							},
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
										IPv6Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
			want: &v2.CiliumNode{
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Allocated: []ipamTypes.IPAMPoolAllocation{
								{
									Pool:  "default",
									CIDRs: []string{"192.168.0.0/24", "f00d::/96"},
								},
							},
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
										IPv6Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "cannot occupy already allocated CIDR",
			fields: fields{
				pools: map[string]cidrPool{
					"default": {
						v4: newAllocatorsWithOccupied(24, []string{"192.168.0.0/16"}, []string{"192.168.1.0/24"}),
					},
				},
				nodes: map[string]poolToCIDRs{},
				ready: true,
			},
			args: &v2.CiliumNode{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-node",
				},
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Allocated: []ipamTypes.IPAMPoolAllocation{
								{
									Pool:  "default",
									CIDRs: []string{"192.168.1.0/24"},
								},
							},
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
			want: &v2.CiliumNode{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-node",
				},
				Spec: v2.NodeSpec{
					IPAM: ipamTypes.IPAMSpec{
						Pools: ipamTypes.IPAMPoolSpec{
							Allocated: []ipamTypes.IPAMPoolAllocation{
								{
									Pool:  "default",
									CIDRs: []string{"192.168.0.0/24"}, // TODO why does allocate a different cidr for it?
								},
							},
							Requested: []ipamTypes.IPAMPoolRequest{
								{
									Pool: "default",
									Needed: ipamTypes.IPAMPoolDemand{
										IPv4Addrs: 10,
									},
								},
							},
						},
					},
				},
			},
			wantErr: "cidr 192.168.1.0/24 has already been allocated",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PoolAllocator{
				pools: tt.fields.pools,
				nodes: tt.fields.nodes,
				ready: tt.fields.ready,
			}
			if err := p.AllocateToNode(tt.args); err != nil && (tt.wantErr == "" || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Errorf("AllocateToNode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != nil {
				p.PopulateNodeSpec(tt.args)
				if diff := cmp.Diff(tt.args, tt.want); diff != "" {
					t.Errorf("AllocateToNode() diff = %s", diff)
				}
			}
		})
	}
}

func Test_addrsInPrefix(t *testing.T) {
	mustParseBigInt := func(s string) *big.Int {
		r := new(big.Int)
		r.SetString(s, 0)
		return r
	}

	tests := []struct {
		name string
		args netip.Prefix
		want *big.Int
	}{
		{
			name: "ipv4",
			args: netip.MustParsePrefix("10.0.0.0/24"),
			want: big.NewInt(254),
		},
		{
			name: "ipv6",
			args: netip.MustParsePrefix("f00d::/48"),
			want: mustParseBigInt("1208925819614629174706174"),
		},
		{
			name: "zero",
			args: netip.Prefix{},
			want: big.NewInt(0),
		},
		{
			name: "two",
			args: netip.MustParsePrefix("10.0.0.0/30"),
			want: big.NewInt(2),
		},
		{
			name: "underflow /31",
			args: netip.MustParsePrefix("10.0.0.0/31"),
			want: big.NewInt(0),
		},
		{
			name: "underflow /32",
			args: netip.MustParsePrefix("10.0.0.0/32"),
			want: big.NewInt(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addrsInPrefix(tt.args); got.Cmp(tt.want) != 0 {
				t.Errorf("addrsInPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
