# awsp – AWS Profile Switcher

CLI para cambiar entre perfiles y regiones de AWS. Modo interactivo con selector de perfil y región.

## Instalación

**Opción 1 – Un solo comando (recomendado, no requiere Go)**

Descarga el binario desde [Releases](https://github.com/danixts/awsp/releases) e instala en `~/.local/bin` + configura shell (awsp + completion):

```bash
curl -sSL https://raw.githubusercontent.com/danixts/awsp/main/install.sh | bash
```

O con wget:

```bash
wget -qO- https://raw.githubusercontent.com/danixts/awsp/main/install.sh | bash
```

Luego reinicia la terminal o ejecuta `source ~/.zshrc` (o tu config).

**Opción 2 – Desde código (requiere Go)**

Solo necesitas Go si compilas desde el repo:

```bash
git clone https://github.com/danixts/awsp.git && cd awsp
go run . install
```

O si ya tienes el binario local: `./aws-profile install`.

**Opción 3 – Make (solo copia)**

```bash
make install
```

Copia el binario a `$(GOPATH)/bin/aws-profile`. Asegúrate de tener ese directorio en tu `PATH`.

## Releases (automático con GitHub Action)

Al hacer push de un **tag** `v*` (ej. `v1.0.0`), el workflow [`.github/workflows/release.yml`](.github/workflows/release.yml) compila los binarios y publica la release con los assets.

```bash
git tag v1.0.0
git push origin v1.0.0
```

Tras unos minutos, la release aparecerá en [Releases](https://github.com/danixts/awsp/releases) con los binarios: `aws-profile-linux-amd64`, `aws-profile-linux-arm64`, `aws-profile-darwin-amd64`, `aws-profile-darwin-arm64`, `aws-profile-windows-amd64.exe`.

Para compilar y subir manualmente: `make release` y luego subir los archivos de `dist/` a una release en GitHub.

## Uso

Binario: **aws-profile**. Atajo: **awsp** (función que abre el menú y aplica las variables).

| Comando | Descripción |
|---------|-------------|
| `awsp` | Atajo: menú interactivo y aplicar variables (recomendado) |
| `aws-profile` | Igual que `awsp` (modo interactivo) |
| `aws-profile switch [profile] [region]` | Cambiar perfil (o interactivo sin args) |
| `aws-profile list` | Lista todos los perfiles (★ = favorito) |
| `aws-profile current` | Muestra el perfil actual |
| `aws-profile favorite add/remove/list` | Gestionar favoritos |
| `aws-profile install` | Instalar binario y configurar shell |

**Flags**

- `-v`, `--validate`: Verifica credenciales con `aws sts get-caller-identity` antes de exportar (modo interactivo y `switch` con perfil).
- `-f`, `--full`: Incluye `AWS_SDK_LOAD_CONFIG=1` en el export (disponible en raíz y en `switch`).

**Cache y favoritos**

- El último perfil y región usados se guardan en `~/.config/awsp/last.json` y se usan como valor por defecto en el selector.
- Los favoritos se guardan en `~/.config/awsp/favorites` y aparecen primero en el selector y con ★ en `aws-profile list`.

## Aplicar las variables en la shell

- **Recomendado:** escribe **awsp** (función corta): se abre el menú y se aplican las variables.
- O imprime exports y aplícalos: `eval $(aws-profile)` o `eval $(aws-profile switch prod us-east-1)`.

La función **awsp** la añade `aws-profile install`; es un atajo para el menú interactivo.

## Autocompletado

Con `aws-profile install` se genera el script en `~/.config/awsp/completion.zsh` (o `.bash`). Tras `source ~/.zshrc`, al escribir **aws-profile** y Tab verás: `switch`, `list`, `current`, `favorite`, `install`, etc.

## Requisitos

- Perfiles en `~/.aws/credentials`
- Opcional: regiones en `~/.aws/config`

## Desarrollo

Requisito: **Go** (solo para compilar desde código; la instalación con `curl` no necesita Go).

```bash
make build    # compilar binario local
make run      # ejecutar con go run
make release  # binarios para Linux/macOS/Windows (dist/)
make clean    # borrar binario y dist/
```
