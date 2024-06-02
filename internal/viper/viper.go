package viper

import (
	"os"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var v = viper.NewWithOptions(viper.KeyDelimiter("::"))
var decoderOptions = viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
	mapstructure.StringToTimeDurationHookFunc(),
	mapstructure.StringToSliceHookFunc(","),
	mapstructure.TextUnmarshallerHookFunc(),
))

func getDecoderOpts(opts ...viper.DecoderConfigOption) []viper.DecoderConfigOption {
	newOpts := []viper.DecoderConfigOption{decoderOptions}
	newOpts = append(newOpts, opts...)
	return newOpts
}

func GetViper() *viper.Viper                                 { return v }
func Get(key string) any                                     { return v.Get(key) }
func SetConfigFile(in string)                                { v.SetConfigFile(in) }
func SetEnvPrefix(in string)                                 { v.SetEnvPrefix(in) }
func ConfigFileUsed() string                                 { return v.ConfigFileUsed() }
func AddConfigPath(in string)                                { v.AddConfigPath(in) }
func Sub(key string) *viper.Viper                            { return v.Sub(key) }
func GetString(key string) string                            { return v.GetString(key) }
func GetBool(key string) bool                                { return v.GetBool(key) }
func GetInt(key string) int                                  { return v.GetInt(key) }
func GetInt32(key string) int32                              { return v.GetInt32(key) }
func GetInt64(key string) int64                              { return v.GetInt64(key) }
func GetUint(key string) uint                                { return v.GetUint(key) }
func GetUint16(key string) uint16                            { return v.GetUint16(key) }
func GetUint32(key string) uint32                            { return v.GetUint32(key) }
func GetUint64(key string) uint64                            { return v.GetUint64(key) }
func GetFloat64(key string) float64                          { return v.GetFloat64(key) }
func GetTime(key string) time.Time                           { return v.GetTime(key) }
func GetDuration(key string) time.Duration                   { return v.GetDuration(key) }
func GetIntSlice(key string) []int                           { return v.GetIntSlice(key) }
func GetStringSlice(key string) []string                     { return v.GetStringSlice(key) }
func GetStringMap(key string) map[string]any                 { return v.GetStringMap(key) }
func GetStringMapString(key string) map[string]string        { return v.GetStringMapString(key) }
func GetStringMapStringSlice(key string) map[string][]string { return v.GetStringMapStringSlice(key) }
func GetSizeInBytes(key string) uint                         { return v.GetSizeInBytes(key) }
func UnmarshalKey(key string, rawVal any, opts ...viper.DecoderConfigOption) error {
	return v.UnmarshalKey(key, rawVal, getDecoderOpts()...)
}
func Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error {
	return v.Unmarshal(rawVal, getDecoderOpts()...)
}
func UnmarshalExact(rawVal any, opts ...viper.DecoderConfigOption) error {
	return v.UnmarshalExact(rawVal, getDecoderOpts()...)
}
func BindPFlags(flags *pflag.FlagSet) error                { return v.BindPFlags(flags) }
func BindPFlag(key string, flag *pflag.Flag) error         { return v.BindPFlag(key, flag) }
func BindFlagValues(flags viper.FlagValueSet) error        { return v.BindFlagValues(flags) }
func BindFlagValue(key string, flag viper.FlagValue) error { return v.BindFlagValue(key, flag) }
func BindEnv(input ...string) error                        { return v.BindEnv(input...) }
func MustBindEnv(input ...string)                          { v.MustBindEnv(input...) }
func IsSet(key string) bool                                { return v.IsSet(key) }
func AutomaticEnv()                                        { v.AutomaticEnv() }
func SetEnvKeyReplacer(r *strings.Replacer)                { v.SetEnvKeyReplacer(r) }
func RegisterAlias(alias, key string)                      { v.RegisterAlias(alias, key) }
func InConfig(key string) bool                             { return v.InConfig(key) }
func SetDefault(key string, value any)                     { v.SetDefault(key, value) }
func Set(key string, value any)                            { v.Set(key, value) }
func ReadInConfig() error                                  { return v.ReadInConfig() }
func MergeInConfig() error                                 { return v.MergeInConfig() }
func AllKeys() []string                                    { return v.AllKeys() }
func AllSettings() map[string]any                          { return v.AllSettings() }
func SetConfigName(in string)                              { v.SetConfigName(in) }
func SetConfigType(in string)                              { v.SetConfigType(in) }
func SetConfigPermissions(perm os.FileMode)                { v.SetConfigPermissions(perm) }
func Debug()                                               { v.Debug() }
