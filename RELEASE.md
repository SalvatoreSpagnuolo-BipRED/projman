# Release Guide

Guida rapida per creare una nuova release di Projman.

## 📋 Prerequisiti

- [GitHub CLI](https://cli.github.com/) installata e autenticata: `gh auth login`
- Accesso push al repository

## 🚀 Creare una Release

### Opzione 1: Da Locale (Consigliato)

```bash
make release VERSION=v1.2.3
```

### Opzione 2: Da GitHub Web

1. Vai su [Actions → Release](https://github.com/SalvatoreSpagnuolo-BipRED/projman/actions/workflows/release.yml)
2. Click su **Run workflow**
3. Inserisci la versione (es: `v1.2.3`)
4. Click su **Run workflow**

## ⚙️ Cosa Succede Automaticamente

1. ✅ Badge versione in README aggiornato
2. ✅ Commit delle modifiche su `main`
3. ✅ Tag Git creato (es: `v1.2.3`)
4. ✅ Build multi-piattaforma (Windows, macOS, Linux)
5. ✅ Changelog generato dai commit
6. ✅ GitHub Release pubblicata con binari

## 📝 Convenzioni

- **Formato versione**: `vX.Y.Z` (es: `v1.2.3`)
- **Commit messages**: Usa [Conventional Commits](https://www.conventionalcommits.org/)
  - `feat:` - Nuove funzionalità
  - `fix:` - Bug fix
  - `chore:` - Manutenzione
  - `docs:` - Documentazione

## 🔍 Monitorare la Release

```bash
# Lista workflow run
gh run list --workflow=release.yml

# Dettagli ultimo run
gh run view
```

## 🧪 Test Locale (Senza Pubblicare)

```bash
make snapshot
```

Crea build locali in `dist/` senza pubblicare release.
