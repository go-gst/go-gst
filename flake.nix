{
  description = "Golang+gstreamer dev environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
    ...
  }: let
    systems = ["x86_64-linux" "aarch64-darwin"];
    forAllSystems = nixpkgs.lib.genAttrs systems;
  in {
    formatter = forAllSystems (system: nixpkgs.legacyPackages.${system}.alejandra);
    devShells = forAllSystems (system: let
      pkgs = import nixpkgs {
        inherit system;
      };
    in {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go
          pkg-config

          gst_all_1.gst-editing-services
          gst_all_1.gst-libav
          gst_all_1.gst-plugins-bad
          gst_all_1.gst-plugins-base
          gst_all_1.gst-plugins-good
          gst_all_1.gst-plugins-rs
          gst_all_1.gst-plugins-ugly
          gst_all_1.gst-rtsp-server
          gst_all_1.gst-vaapi
          gst_all_1.gstreamer

          #gotk4 deps:
          atk
          gtk3
          gtk4
          glib
          graphene
          gdk-pixbuf
          gobject-introspection
          librsvg
        ];

        GO111MODULE = "on";

        # needed for running delve
        # https://github.com/go-delve/delve/issues/3085
        # https://nixos.wiki/wiki/C#Hardening_flags
        hardeningDisable = ["all"];

        # print the go version and gstreamer version on shell startup
        shellHook = ''
          ${pkgs.go}/bin/go version
          ${pkgs.gst_all_1.gstreamer}/bin/gst-launch-1.0 --version
        '';
      };
    });
  };
}
