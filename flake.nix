{
  description = "Identity Service";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    let
      readVersion = { versionPath, fallback }:
        if builtins.pathExists versionPath
        then builtins.replaceStrings ["\n"] [""] (builtins.readFile versionPath)
        else fallback;
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        pname = "identity-service";
        version = readVersion {
          versionPath = ./VERSION;
          fallback = "latest";
        };
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            inherit pname version;

            src = ./.;
            vendorHash = "sha256-WJ3s+XgETHMkSFeEhpOdVjVuFno3SsoKFGtV4IjU64Y=";
            subPackages = [ "cmd/identity-service" ];
            buildFlags = [ "-mod=mod" ];
            ldflags = [
              "-s -w"
              "-X main.version=${version}"
            ];
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = pname;
            tag = version;
            contents = [
              self.packages.${system}.default
              pkgs.cacert
            ];
            config = {
              Cmd = [ "${self.packages.${system}.default}/bin/identity-service" ];
              Env = [
                "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              ];
            };
          };
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_25
            goose
            go-swag
            just
          ];

          shellHook = ''
            echo "Development environment ready!"
            export GOPATH=$HOME/go
            export PATH=$GOPATH/bin:$PATH
          '';
        };
      });
}
