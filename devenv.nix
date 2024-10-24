{
  pkgs,
  ...
}:

{
  dotenv.enable = true;

  packages = with pkgs; [
    go
    gopls
    golangci-lint
    go-tools
    git

    xorg.libX11
    xorg.libXrandr
    xorg.libXinerama
    xorg.libXcursor
    xorg.libXi
    xorg.xorgproto
  ];

  languages.go = {
    enable = true;
  };

  scripts.start.exec = ''
    go run ./cmd/cbtxt/main.go ./cmd/cbtxt/traverse.go ./cmd/cbtxt/gitignore.go
  '';

  pre-commit.hooks.gofmt.enable = true;
  pre-commit.hooks.golangci-lint.enable = true;
}
