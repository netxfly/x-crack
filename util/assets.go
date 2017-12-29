/*

Copyright (c) 2017 xsec.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package util

import (
	"x-crack/models"
)

func DivideAsset(ipList []models.IpAddr) (assetsGroup map[string][]models.IpAddr) {
	assetsGroup = make(map[string][]models.IpAddr)
	for _, addr := range ipList {
		protocol := addr.Protocol
		if _, ok := assetsGroup[protocol]; ok {
			assetsGroup[protocol] = append(assetsGroup[protocol], addr)
		} else {
			assetsGroup[protocol] = make([]models.IpAddr, 0)
			assetsGroup[protocol] = append(assetsGroup[protocol], addr)
		}
	}
	return assetsGroup
}
