{ pkgs ?
  import (fetchTarball "http://nixos.org/channels/nixos-20.09/nixexprs.tar.xz")
  { } }:

pkgs.buildGoModule {
  name = "einkauf";
  version = "1.0.0";
  src = ./.;
  vendorSha256 = null;
}

