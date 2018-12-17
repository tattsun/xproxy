package proxy

import (
	"net"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"
)

type Matcher struct {
	hostGlobs []glob.Glob
	ipNets    []*net.IPNet
}

func NewMatcher(hosts []string, ips []string) (*Matcher, error) {
	hostGlobs := make([]glob.Glob, len(hosts))
	for i, host := range hosts {
		g, err := glob.Compile(host)
		if err != nil {
			return nil, errors.Wrapf(err, `failed to compile host "%s"`, host)
		}
		hostGlobs[i] = g
	}

	ipNets := make([]*net.IPNet, len(ips))
	for i, ip := range ips {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, errors.Wrapf(err, `failed to parse cidr "%s"`, ip)
		}
		ipNets[i] = ipNet
	}

	return &Matcher{
		hostGlobs: hostGlobs,
		ipNets:    ipNets,
	}, nil
}

func (m *Matcher) Host(target string) bool {
	for _, g := range m.hostGlobs {
		if g.Match(target) {
			return true
		}
	}
	return false
}

func (m *Matcher) IP(target net.IP) bool {
	for _, ipNet := range m.ipNets {
		if ipNet.Contains(target) {
			return true
		}
	}
	return false
}
