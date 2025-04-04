{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.pkg-config
  ];

  buildInputs = [
    pkgs.go
    pkgs.raylib
    pkgs.wayland
    pkgs.wayland-protocols
    pkgs.libxkbcommon
    pkgs.xorg.libX11
    pkgs.xorg.libXrandr
    pkgs.xorg.libXinerama
    pkgs.xorg.libXcursor
    pkgs.xorg.libXxf86vm
    pkgs.xorg.libXext
    pkgs.xorg.libXi  # Required for XInput2.h
    pkgs.mesa
  ];
}

