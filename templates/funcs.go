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
	"strings"

	"github.com/dspinhirne/netaddr-go"
)

var templateFuncs = map[string]interface{}{
	"split":   split,
	"ip4":     ip4,
	"ip4mask": ip4mask,
	//"ip6":      ip6,
}

func ip4(ip string, idx int) (string, error) {
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

func ip4mask(ip string) (string, error) {
	ipv4net, err := netaddr.ParseIPv4Net(ip)
	if err != nil {
		return "", err
	}

	return ipv4net.Netmask().Extended(), nil
}

func split(str string, sep string, idx int) string {
	arr := strings.Split(str, sep)
	if idx > len(arr)-1 {
		return ""
	}

	return arr[idx]
}
