{
  description = "Easily convert epoch timestamps to human readable formats and vice versa.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
          ];
        };

        packages.default = pkgs.buildGoModule {
          pname = "epoch";
          version = "undefined";
          src = ./.;
          vendorHash = null;
          subPackages = [ "./cmd/epoch" ];
        };
      });
}
