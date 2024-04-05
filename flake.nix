{
  description = "simple http file server with tui";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
      };
    in {
      defaultPackage = pkgs.buildGoModule rec {
        pname = "cths";
        version = "1";
        src = ./.;
        vendorHash = "sha256-9INrlvLTIm9ZJ97T1K0PLbcPAgVNi+KOs1Ac9aCF8YU=";

        ldflags = ["-s -w"];
        CGO_ENABLED = 0;

        meta = with pkgs.lib; {
          description = "simple http file server with tui";
          license = licenses.wtfpl;
        };
      };
      formatter = pkgs.alejandra;
    });
}
