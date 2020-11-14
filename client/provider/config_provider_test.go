package provider

import (
	"fmt"
	"testing"

	"github.com/gojek/stevedore/client/internal/mocks/micro/go-micro"
	"github.com/gojek/stevedore/pkg/plugin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestConsulConfigProvider_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfig := mocks.NewMockConfig(ctrl)
	c := ConsulConfigProvider{conf: mockConfig}
	mockValue := mocks.NewMockValue(ctrl)

	type args struct {
		context map[string]string
		data    interface{}
	}
	tests := []struct {
		name   string
		args   args
		preRun func()
		want   map[string]interface{}
		err    error
	}{
		{
			name: "should return error if host is not present",
			args: args{
				context: map[string]string{},
				data:    new(interface{}),
			},
			want: nil,
			err:  fmt.Errorf("host not set"),
		},
		{
			name: "should return error if port is not present",
			args: args{
				context: map[string]string{"host": "consul"},
				data:    new(interface{}),
			},
			want: nil,
			err:  fmt.Errorf("port not set"),
		},
		{
			name: "should return error if prefix is not present",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234"},
				data:    new(interface{}),
			},
			want: nil,
			err:  fmt.Errorf("prefix not set"),
		},
		{
			name: "should return error if strip-prefix is not present",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234", "prefix": "prefix"},
				data:    new(interface{}),
			},
			want: nil,
			err:  fmt.Errorf("strip-prefix not set"),
		},
		{
			name: "should return error if strip-prefix is not bool",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234", "prefix": "prefix", "strip-prefix": "not-bool"},
				data:    new(interface{}),
			},
			want: nil,
			err:  fmt.Errorf(`invalid value for strip-prefix: strconv.ParseBool: parsing "not-bool": invalid syntax`),
		},
		{
			name: "should return error if data is parse-able",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234", "prefix": "prefix", "strip-prefix": "false"},
				data:    new(interface{}),
			},
			want: nil,
			err:  fmt.Errorf(`invalid consul configs: '' expected a map, got 'interface'`),
		},
		{
			name: "should fail when unable to load config from consul",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234", "prefix": "prefix", "strip-prefix": "false"},
				data: map[interface{}]interface{}{
					"path": []string{"path1"},
				},
			},
			preRun: func() {
				mockConfig.EXPECT().Load(gomock.Any()).Return(fmt.Errorf("error while loading"))
			},
			want: nil,
			err:  fmt.Errorf("error loading from consul: error while loading"),
		},
		{
			name: "should return empty map for empty path",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234", "prefix": "prefix", "strip-prefix": "false"},
				data:    make(map[interface{}]interface{}),
			},
			preRun: func() {
				mockConfig.EXPECT().Load(gomock.Any()).Return(nil)
			},
			want: make(map[string]interface{}),
			err:  nil,
		},
		{
			name: "should fail when unable to load config from consul",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234", "prefix": "prefix", "strip-prefix": "false"},
				data: map[interface{}]interface{}{
					"path": []string{"path1"},
				},
			},
			preRun: func() {
				mockConfig.EXPECT().Load(gomock.Any()).Return(nil)
				mockValue := mocks.NewMockValue(ctrl)
				mockConfig.EXPECT().Get("path1").Return(mockValue)
				mockValue.EXPECT().Scan(gomock.Any()).Return(fmt.Errorf("some scan error"))
			},
			want: nil,
			err:  fmt.Errorf("unable to decode configs under path path1: some scan error"),
		},
		{
			name: "should fetch from consul",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234", "prefix": "prefix", "strip-prefix": "false"},
				data: map[interface{}]interface{}{
					"path": []string{"path1"},
				},
			},
			preRun: func() {
				mockConfig.EXPECT().Load(gomock.Any()).Return(nil)
				mockValue := mocks.NewMockValue(ctrl)
				mockConfig.EXPECT().Get("path1").Return(mockValue)
				mockValue.EXPECT().Scan(gomock.Any()).Do(func(mapAddress interface{}) {
					config := *mapAddress.(*map[string]interface{})
					config["hi"] = "bye"
				})
			},
			want: map[string]interface{}{"hi": "bye"},
			err:  nil,
		},
		{
			name: "should fetch and merge from consul multi-path",
			args: args{
				context: map[string]string{"host": "consul", "port": "1234", "prefix": "prefix", "strip-prefix": "false"},
				data: map[interface{}]interface{}{
					"path": []string{"path1", "path2"},
				},
			},
			preRun: func() {
				mockConfig.EXPECT().Load(gomock.Any()).Return(nil)
				mockConfig.EXPECT().Get("path1").Return(mockValue)
				mockValue.EXPECT().Scan(gomock.Any()).Do(func(mapAddress interface{}) {
					config := *mapAddress.(*map[string]interface{})
					config["key1"] = "value1"
					config["key2"] = "value2"
				})
				mockValue2 := mocks.NewMockValue(ctrl)
				mockConfig.EXPECT().Get("path2").Return(mockValue2)
				mockValue2.EXPECT().Scan(gomock.Any()).Do(func(mapAddress interface{}) {
					config := *mapAddress.(*map[string]interface{})
					config["key3"] = "value3"
					config["key2"] = "value4"
				})
			},
			want: map[string]interface{}{
				"key1": "value1",
				"key2": "value4",
				"key3": "value3",
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preRun != nil {
				test.preRun()
			}
			got, err := c.Fetch(test.args.context, test.args.data)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestConsulConfigProvider_Flags(t *testing.T) {
	t.Run("should returns consul flags", func(t *testing.T) {
		c := ConsulConfigProvider{}
		expectedFlags := []plugin.Flag{
			{Name: "host", Default: "http://127.0.0.1", Usage: "host for consul"},
			{Name: "port", Default: "8500", Usage: "port for consul"},
			{Name: "prefix", Default: "/", Usage: "prefix for contacting consul"},
			{Name: "strip-prefix", Default: "true", Usage: "strip-prefix indicates whether to remove the prefix from config entries, or leave it in place."},
		}

		actual, err := c.Flags()

		assert.NoError(t, err)
		assert.Equal(t, expectedFlags, actual)
	})

}
