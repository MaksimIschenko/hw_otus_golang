//go:build !bench
// +build !bench

package stats

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

func TestGetUserEmails(t *testing.T) {
	t.Run("find 'com' domain", func(t *testing.T) {
		result, err := GetUserEmails(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, LightUsers{
			{Email: "mLynch@broWsecat.com"},
			{Email: "RoseSmith@Browsecat.com"},
			{Email: "nulla@Linktype.com"},
		}, result)
	})

	t.Run("find 'gov' domain", func(t *testing.T) {
		result, err := GetUserEmails(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, LightUsers{
			{Email: "aliquid_qui_ea@Browsedrive.gov"},
		}, result)
	})

	t.Run("find 'unknown' domain", func(t *testing.T) {
		result, err := GetUserEmails(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, LightUsers{}, result)
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := GetUserEmails(bytes.NewBufferString(""), "com")
		require.NoError(t, err)
		require.Equal(t, LightUsers{}, result)
	})
}

func TestGetDomainStat(t *testing.T) {
	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}
