# awsp – AWS Profile Switcher

CLI para cambiar entre perfiles y regiones de AWS. Modo interactivo con selector de perfil y región.

## Instalación

**Opción 1 – Instalador (recomendado, multiplataforma)**

Desde el repo (con Go instalado):

```bash
go run . install
```

o si ya tienes el binario:

```bash
./aws-profile install
```

El comando:
1. **Si estás en el proyecto** (hay `go.mod`): compila con `go build` y luego copia el binario `aws-profile` a `~/.local/bin`.
2. **Si ejecutas un `aws-profile` ya instalado**: no recompila; solo copia el binario actual. Para una versión nueva, ejecuta `go run . install` desde el proyecto.
3. Pide elegir shell: **Zsh (Oh My Zsh)**, **Bash** o **Windows (PowerShell)**.
4. Añade la función **awsp** (atajo interactivo) y el autocompletado para `aws-profile`.

Reinicia la terminal o ejecuta `source ~/.zshrc` (o tu config).

**Opción 2 – Make (solo copia)**

```bash
make install
```

Copia el binario a `$(GOPATH)/bin/aws-profile`. Asegúrate de tener ese directorio en tu `PATH`.

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

```bash
make build   # compilar
make run     # ejecutar con go run
make clean   # borrar binario
```
