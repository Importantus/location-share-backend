{
  description = "Location-Share Backend";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  outputs = { nixpkgs, self, ... }:
    let
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      nixosModules.default = { pkgs, lib, config, ... }:
        let
          cfg = config.services.location-share;
          username = "locationshare";
          app = throw "make app derivation";

          # has to be equal to username to ensure ownership for this db to user
          dbName = username;
        in
        {
          options.services.location-share = with lib.types;
            {
              enable = lib.mkEnableOption "Location-Share Backend";
              port = lib.mkOption {
                type = int;
                default = 8000;
                description = "port to liten at";
              };
              clientOrigin = lib.mkOption {
                type = uniq str;
                default = "http://localhost:8000";
                description = "url the client connects to. possibly relevant if behind a reverse proxy";
              };
              jwtSecret = lib.mkOption {
                type = uniq str;
                default = "verySecret";
                description = "JWT Secret";
              };
              registrationSecret = lib.mkOption {
                type = nullOr (uniq str);
                default = null;
                description = "This secret is needed to register a new user. If set to empty string or null, anyone can register";
              };
              googleApplicationCredentials = lib.mkOption {
                type = attrs;
                description = "google application credentials.";
              };
              dbPassword = lib.mkOption {
                type = uniq str;
                default = "password123";
                description = "Password for Postgres db. is ignored by the postgresqls server anyway";
              };
            };

          config = lib.mkIf cfg.enable {
            services.postgresql = {
              enable = true;
              ensureDatabases = [ dbName ];
              ensureUsers = [{
                name = username;
                ensureDBOwnership = true;
              }];
              authentication = ''
                #type database  DBuser  auth-method
                local all       all     trust
              '';
            };

            systemd.services.locationshare = {
              after = [ "postgresql.service" ];
              wants = [ "postgresql.service" ];
              wantedBy = [ "multi-user.target" ];
              script = ''

              # Set up working directory
              WORK_DIR=$(mktemp -d)
              cd "$WORK_DIR"

              cat ${pkgs.writeText "app.env" ''
                POSTGRES_HOST=/var/run/postgresql
                POSTGRES_USER=${username}
                POSTGRES_PASSWORD=${cfg.dbPassword}
                POSTGRES_DB=${dbName}
                POSTGRES_PORT=${builtins.toString config.services.postgresql.settings.port}

                PORT=${builtins.toString cfg.port}
                CLIENT_ORIGIN=${cfg.clientOrigin}

                JWT_SECRET=${cfg.jwtSecret}

                # This secret is needed to register a new user. If ste to empty string, anyone can register
                REGISTRATION_SECRET=${if cfg.registrationSecret == null then "" else cfg.registrationSecret}

                GOOGLE_APPLICATION_CREDENTIALS=${cfg.googleApplicationCredentials
                  |> builtins.toJSON
                  |> pkgs.writeText "credentials.json"
                }
              ''} > ./app.env
              ${
                self.packages.${config.nixpkgs.localSystem.system}.default
              }/bin/migrate
              ${
                self.packages.${config.nixpkgs.localSystem.system}.default
              }/bin/location-share-backend
            '';
              serviceConfig = {
                User = username;
                Restart = "on-failure";
              };
            };

            users.users.${username} = {
              isSystemUser = true;
              group = username;
            };
            users.groups.${username} = { };
          };
        };
      packages = forAllSystems ({ pkgs }: {
        default = pkgs.buildGo123Module rec {
          pname = "locationshare";
          version = "1.0.0";
          src = ./.;
          vendorHash = "sha256-e0J/bSYWQNqdstQBrKf17sPvTOQtFPB5PPAWuyKgqUI=";
        };
      });
    };
}
