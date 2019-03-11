// Copyright © 2019 Pawel Potrykus <pawel.potrykus@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package templates

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dspinhirne/netaddr-go"
)

var templateFuncs = map[string]interface{}{
	"split":           Split,
	"ip4":             IP4,
	"ip4mask":         IP4Mask,
	"ip4cidr":         IP4Cidr,
	"ip4mask_to_cidr": IP4MaskToCidr,
	"ip4cidr_to_mask": IP4CidrToMask,
	//"ip6":      IP6,
	//"ip6mask": IP6Mask,
}

func IP4(ip string, idx int) (string, error) {
	if idx < 0 {
		return "", fmt.Errorf("negative value of argument passed to ip4 func not allowed")
	}

	ipv4net, err := netaddr.ParseIPv4Net(ip)
	if err != nil {
		return "", err
	}

	ipv4 := ipv4net.Network()

	for i := 0; i < idx; i++ {
		ipv4 = ipv4.Next()
	}

	return ipv4.String(), nil
}

func IP4Mask(ip string) (string, error) {
	ipv4net, err := netaddr.ParseIPv4Net(ip)
	if err != nil {
		return "", err
	}

	return ipv4net.Netmask().Extended(), nil
}

func IP4Cidr(ip string) (string, error) {
	ipv4net, err := netaddr.ParseIPv4Net(ip)
	if err != nil {
		return "", err
	}

	cidr := strconv.Itoa(int(ipv4net.Netmask().PrefixLen()))

	return cidr, nil
}

func IP4MaskToCidr(mask string) (string, error) {
	mask32, err := netaddr.ParseMask32(mask)
	if err != nil {
		return "", err
	}

	cidr := strconv.Itoa(int(mask32.PrefixLen()))

	return cidr, nil
}

func IP4CidrToMask(cidr string) (string, error) {
	cidrInt, err := strconv.Atoi(cidr)
	if err != nil {
		return "", err
	}

	mask32, err := netaddr.NewMask32(uint(cidrInt))
	if err != nil {
		return "", err
	}

	return mask32.Extended(), nil
}

func Split(str string, sep string, idx int) string {
	arr := strings.Split(str, sep)
	if idx > len(arr)-1 {
		return ""
	}

	return arr[idx]
}
