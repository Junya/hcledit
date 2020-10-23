package cmd

import (
	"testing"
)

func TestAttributeGet(t *testing.T) {
	src := `terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "minamijoyo-hcledit"
    key    = "services/hoge/dev/terraform.tfstate"
  }
}
`

	cases := []struct {
		name string
		args []string
		ok   bool
		want string
	}{
		{
			name: "simple",
			args: []string{"terraform.backend.s3.key"},
			ok:   true,
			want: "\"services/hoge/dev/terraform.tfstate\"\n",
		},
		{
			name: "no match",
			args: []string{"hoge"},
			ok:   true,
			want: "",
		},
		{
			name: "no args",
			args: []string{},
			ok:   false,
			want: "",
		},
		{
			name: "too many args",
			args: []string{"hoge", "fuga"},
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(newAttributeGetCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}

func TestAttributeSet(t *testing.T) {
	src := `terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "minamijoyo-hcledit"
    key    = "services/hoge/dev/terraform.tfstate"
  }
}
module "hoge" {
  source = "./hoge"
  env    = "dev"
}
`

	cases := []struct {
		name string
		args []string
		ok   bool
		want string
	}{
		{
			name: "string literal",
			args: []string{"terraform.backend.s3.key", `"services/fuga/dev/terraform.tfstate"`},
			ok:   true,
			want: `terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "minamijoyo-hcledit"
    key    = "services/fuga/dev/terraform.tfstate"
  }
}
module "hoge" {
  source = "./hoge"
  env    = "dev"
}
`,
		},
		{
			name: "string literal to variable reference",
			args: []string{"module.hoge.env", "var.env"},
			ok:   true,
			want: `terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "minamijoyo-hcledit"
    key    = "services/hoge/dev/terraform.tfstate"
  }
}
module "hoge" {
  source = "./hoge"
  env    = var.env
}
`,
		},
		{
			name: "no match",
			args: []string{"hoge", "fuga"},
			ok:   true,
			want: src,
		},
		{
			name: "no args",
			args: []string{},
			ok:   false,
			want: "",
		},
		{
			name: "1 arg",
			args: []string{"hoge"},
			ok:   false,
			want: "",
		},
		{
			name: "too many args",
			args: []string{"hoge", "fuga", "piyo"},
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(newAttributeSetCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}

func TestAttributeRm(t *testing.T) {
	src := `locals {
  service = "hoge"
  env     = "dev"
  region  = "ap-northeast-1"
}`

	cases := []struct {
		name string
		args []string
		ok   bool
		want string
	}{
		{
			name: "remove an unused local variable",
			args: []string{"locals.region"},
			ok:   true,
			want: `locals {
  service = "hoge"
  env     = "dev"
}`,
		},
		{
			name: "no match",
			args: []string{"hoge"},
			ok:   true,
			want: src,
		},
		{
			name: "no args",
			args: []string{},
			ok:   false,
			want: "",
		},
		{
			name: "too many args",
			args: []string{"hoge", "fuga"},
			ok:   false,
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newMockCmd(newAttributeRmCmd(), src)
			assertMockCmd(t, cmd, tc.args, tc.ok, tc.want)
		})
	}
}
