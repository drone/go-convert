// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package travis

import "strings"

func convertImageMaybe(ctx *context, ok bool) (image string) {
	if ok {
		return convertImage(ctx)
	} else {
		return
	}
}

func convertImage(ctx *context) string {
	switch strings.ToLower(ctx.config.Language) {
	case "android":
		// TODO
	case "c":
		return convertImageC(ctx)
	case "clojure":
		return convertImageClojure(ctx)
	case "cpp":
		return convertImageC(ctx)
	case "crystal":
		return convertImageCrystal(ctx)
	case "csharp":
		// TODO
	case "d":
		// TODO
	case "dart":
		// TODO
	case "elixir":
		// TODO
	case "elm":
		// TODO
	case "erlang":
		return convertImageErlang(ctx)
	case "go":
		return convertImageGo(ctx)
	case "groovy":
		// TODO
	case "hack":
		// TODO
	case "haskell":
		// TODO
	case "haxe":
		// TODO
	case "java":
		// TODO
	case "julia":
		// TODO
	case "nix":
		// TODO
	case "node_js":
		return convertImageNode(ctx)
	case "objective-c":
		return "" // no docker image for objective c
	case "perl":
		// TODO
	case "perl6":
		// TODO
	case "php":
		// TODO
	case "python":
		return convertImagePy(ctx)
	case "r":
		// TODO
	case "ruby":
		// TODO
	case "rust":
		return convertImageRust(ctx)
	case "scala":
		// TODO
	case "smalltalk":
		// TODO
	case "minimal", "generic", "shell":
		return "ubuntu"
	}
	return "ubuntu"
}

func convertImageC(ctx *context) string {
	if len(ctx.config.Compiler) == 0 {
		return "gcc"
	}
	if len(ctx.config.Compiler) == 1 {
		switch ctx.config.Compiler[0] {
		case "gcc":
			return "gcc"
		case "clang":
			return "gcc" // TODO official clang image
		}
	}
	return "gcc" // TODO strategy to convert C matrix to image
}

func convertImageCrystal(ctx *context) string {
	if len(ctx.config.Crystal) == 0 {
		return "crystallang/crystal"
	}
	if len(ctx.config.Crystal) == 1 {
		version := ctx.config.Python[0]
		if version == "nightly" {
			version = "latest"
		}
		return "crystallang/crystal:" + version
	}
	return "crystallang/crystal:<+matrix.crystal>"
}

func convertImageClojure(ctx *context) string {
	// TODO support for jdk version
	// TODO support for lein version
	return "clojure"
}

func convertImageErlang(ctx *context) string {
	if len(ctx.config.ErlangOTP) == 0 {
		return "erlang"
	}
	if len(ctx.config.ErlangOTP) == 1 {
		version := ctx.config.ErlangOTP[0]
		version = strings.ReplaceAll(version, ".x", "")
		return "erlang:" + version
	}
	return "golang:<+matrix.otp_release>"
}

func convertImageGo(ctx *context) string {
	if len(ctx.config.Go) == 0 {
		return "golang"
	}
	if len(ctx.config.Go) == 1 {
		return "golang:" + strings.ReplaceAll(ctx.config.Go[0], ".x", "")
	}
	return "golang:<+matrix.go>"
}

func convertImageNode(ctx *context) string {
	if len(ctx.config.Node) == 0 {
		return "node"
	}
	if len(ctx.config.Node) == 1 {
		version := ctx.config.Node[0]
		switch version {
		case "lts/*":
			version = "lts"
		case "node":
			version = "latest"
		}
		return "node:" + version
	}
	return "node:<+matrix.node_js>"
}

func convertImagePy(ctx *context) string {
	if len(ctx.config.Python) == 0 {
		return "python"
	}
	if len(ctx.config.Python) == 1 {
		version := ctx.config.Python[0]
		version = strings.TrimSuffix(version, "-dev")
		if version == "nightly" {
			version = "latest"
		}
		return "python:" + version
	}
	return "python:<+matrix.python>"
}

func convertImageRust(ctx *context) string {
	if len(ctx.config.Rust) == 0 {
		return "rust"
	}
	if len(ctx.config.Rust) == 1 {
		version := ctx.config.Rust[0]
		switch version {
		case "stable":
			version = "1"
		case "beta":
			version = "latest"
		case "nightly":
			version = "latest"
		}
		return "rust:" + version
	}
	return "rust:<+matrix.rust>"
}
