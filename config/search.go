package config

type Search struct {
	Url      string `mapstructure:"url" json:"url" yaml:"url"`
	UserName string `mapstructure:"user_name" json:"user_name" yaml:"user_name"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}