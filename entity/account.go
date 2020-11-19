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

package entity

const (
	uuid             = "uuid"
	version          = "version"
	amojoID          = "amojo_id"
	taskTypes        = "task_types"
	userGroups       = "users_groups"
	amojoRights      = "amojo_rights"
	datetimeSettings = "datetime_settings"
)

type Account struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	Subdomain               string `json:"subdomain"`
	CreatedAt               int    `json:"created_at"`
	CreatedBy               int    `json:"created_by"`
	UpdatedAt               int    `json:"updated_at"`
	UpdatedBy               int    `json:"updated_by"`
	CurrentUserID           int    `json:"current_user_id"`
	Country                 string `json:"country"`
	Currency                string `json:"currency"`
	CustomersMode           string `json:"customers_mode"`
	IsUnsortedOn            bool   `json:"is_unsorted_on"`
	MobileFeatureVersion    int    `json:"mobile_feature_version"`
	IsLossReasonEnabled     bool   `json:"is_loss_reason_enabled"`
	IsHelpbotEnabled        bool   `json:"is_helpbot_enabled"`
	IsTechnicalAccount      bool   `json:"is_technical_account"`
	ContactNameDisplayOrder int    `json:"contact_name_display_order"`
	Links                   struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
}

func (a *Account) Relations() []string {
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

