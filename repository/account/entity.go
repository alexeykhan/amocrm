// Copyright (c) 2020 Alexey Khan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package account

import (
	"net/http"
)

const (
	uuid             = "uuid"
	version          = "version"
	amojoID          = "amojo_id"
	taskTypes        = "task_types"
	userGroups       = "users_groups"
	amojoRights      = "amojo_rights"
	datetimeSettings = "datetime_settings"
)

type Entity struct {

}

func Relations() []string {
	return []string{
		amojoID,
		uuid,
		amojoRights,
		userGroups,
		taskTypes,
		version,
		datetimeSettings,
	}
}

func FromResponse(resp *http.Response) (*Entity, error) {
	return nil, nil
}
