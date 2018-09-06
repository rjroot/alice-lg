package main

import (
	"fmt"
	"strings"
)

/*
Implement BGP Communities Lookup Base

We initialize the dictionary with well known communities and
store the representation as a string with : as delimiter.

From: https://www.iana.org/assignments/bgp-well-known-communities/bgp-well-known-communities.xhtml

    0x00000000-0x0000FFFF   Reserved    [RFC1997]
    0x00010000-0xFFFEFFFF   Reserved for Private Use    [RFC1997]

    0xFFFF0000  GRACEFUL_SHUTDOWN   [RFC8326]
    0xFFFF0001  ACCEPT_OWN          [RFC7611]
    0xFFFF0002  ROUTE_FILTER_TRANSLATED_v4  [draft-l3vpn-legacy-rtc]
    0xFFFF0003  ROUTE_FILTER_v4     [draft-l3vpn-legacy-rtc]
    0xFFFF0004  ROUTE_FILTER_TRANSLATED_v6  [draft-l3vpn-legacy-rtc]
    0xFFFF0005  ROUTE_FILTER_v6     [draft-l3vpn-legacy-rtc]
    0xFFFF0006  LLGR_STALE          [draft-uttaro-idr-bgp-persistence]
    0xFFFF0007  NO_LLGR             [draft-uttaro-idr-bgp-persistence]
    0xFFFF0008  accept-own-nexthop  [draft-agrewal-idr-accept-own-nexthop]

    0xFFFF0009-0xFFFF0299   Unassigned

    0xFFFF029A  BLACKHOLE           [RFC7999]

    0xFFFF029B-0xFFFFFF00   Unassigned

    0xFFFFFF01  NO_EXPORT           [RFC1997]
    0xFFFFFF02  NO_ADVERTISE        [RFC1997]
    0xFFFFFF03  NO_EXPORT_SUBCONFED [RFC1997]
    0xFFFFFF04  NOPEER              [RFC3765]
    0xFFFFFF05-0xFFFFFFFF   Unassigned
*/

type BgpCommunities map[string]string

func MakeWellKnownBgpCommunities() BgpCommunities {
	c := BgpCommunities{
		"65535:0": "graceful shutdown",
		"65535:1": "accept own",
		"65535:2": "route filter translated v4",
		"65535:3": "route filter v4",
		"65535:4": "route filter translated v6",
		"65535:5": "route filter v6",
		"65535:6": "llgr stale",
		"65535:7": "no llgr",
		"65535:8": "accept-own-nexthop",

		"65535:666": "blackhole",

		"65535:1048321": "no export",
		"65535:1048322": "no advertise",
		"65535:1048323": "no export subconfed",
		"65535:1048324": "nopeer",
	}

	return c
}

func (self BgpCommunities) Merge(communities BgpCommunities) BgpCommunities {
	merged := BgpCommunities{}

	// Make copy, don't mutate
	for k, v := range self {
		merged[k] = v
	}

	for k, v := range communities {
		merged[k] = v
	}

	return merged
}

type NgBgpCommunities map[string]interface{}

func NgMakeWellKnownBgpCommunities() NgBgpCommunities {
	c := NgBgpCommunities{
		"65535": NgBgpCommunities{
			"0": "graceful shutdown",
			"1": "accept own",
			"2": "route filter translated v4",
			"3": "route filter v4",
			"4": "route filter translated v6",
			"5": "route filter v6",
			"6": "llgr stale",
			"7": "no llgr",
			"8": "accept-own-nexthop",

			"666": "blackhole",

			"1048321": "no export",
			"1048322": "no advertise",
			"1048323": "no export subconfed",
			"1048324": "nopeer",
		},
	}

	return c
}

func (self NgBgpCommunities) Lookup(community string) (string, error) {
	path := strings.Split(community, ":")
	var lookup interface{} // This is all much too dynamic...
	lookup = self

	for _, key := range path {
		clookup, ok := lookup.(NgBgpCommunities)
		if !ok {
			break
		}

		res, ok := clookup[key]
		if !ok {
			// Try to fall back to wildcard key
			res, ok = clookup["*"]
			if !ok {
				break // we did everything we could.
			}
		}

		lookup = res
	}

	label, ok := lookup.(string)
	if !ok {
		return "", fmt.Errorf("community not found")
	}

	return label, nil
}

func (self NgBgpCommunities) Set(community string, label string) {
	path := strings.Split(community, ":")
	var lookup interface{} // This is all much too dynamic...
	lookup = self

	for _, key := range path {
		clookup, ok := lookup.(NgBgpCommunities)
		if !ok {
			break
		}

		res, ok := clookup[key]
		if !ok {
			// Try to fall back to wildcard key
			res, ok = clookup["*"]
			if !ok {
				break // we did everything we could.
			}
		}

		lookup = res
	}

	label, ok := lookup.(string)
	if !ok {
		return "", fmt.Errorf("community not found")
	}

	return label, nil
}
