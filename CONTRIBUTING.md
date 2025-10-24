# Contributing Guide / Руководство по внесению вклада

## English

### How to Add a Package

#### Using the Web Interface (Recommended)

1. **Visit the web interface:**
   Go to [registry website URL]

2. **Sign in with GitHub:**
   Click "Sign in with GitHub" and authorize the application

3. **Upload your package:**
   - Navigate to the "Upload" page
   - Drag and drop your .lua file(s)
   - The system will automatically:
     - Extract metadata from `script_name()`, `script_version()`, `script_author()`
     - Analyze and detect all dependencies
     - Detect security features (FFI, network access, file I/O)
   - Review the detected metadata and edit if needed
   - Add tags (comma-separated)
   - Add source URL (optional)

4. **Create Pull Request:**
   - Click "Upload & Create PR"
   - The system will automatically create a fork, commit files, and open a PR
   - Wait for maintainer review

#### Using the CLI

1. **Download the CLI:**
   Download pre-compiled binary from [Releases](https://github.com/Deps-Tech/deps-registry/releases/latest):
   - Windows: `tools-cli-windows-amd64.exe`
   - Linux: `tools-cli-linux-amd64` or `tools-cli-linux-arm64`
   - macOS: `tools-cli-darwin-amd64` or `tools-cli-darwin-arm64`

2. **Clone the repository:**
```bash
git clone https://github.com/Deps-Tech/deps-registry.git
cd deps-registry
```

3. **Add your script:**
```bash
./tools-cli add script \
  --source ~/path/to/YourScript.lua \
  --tags "helper,rp,automation"
```

The CLI will automatically:
- Extract metadata from `script_name()`, `script_version()`, `script_author()`
- Analyze and detect all dependencies
- Detect security features (FFI, network access, file I/O)
- Generate manifest with SHA256 hashes
- Create proper directory structure

4. **Review the generated files:**
```bash
git diff
```

5. **Commit and push:**
```bash
git add scripts/your-script/
git commit -m "feat(scripts): add your-script v1.0"
git push origin main
```

6. **Create a Pull Request** on GitHub

#### Manual Method

1. Create directory: `scripts/your-script/1.0/`
2. Copy your Lua file(s)
3. Create `dep.json` manifest
4. Run validation: `./tools-cli validate`

### Package Guidelines

- **Naming:** Use lowercase with hyphens (e.g., `fake-documents`)
- **Versioning:** Follow semantic versioning (e.g., `1.0.0`)
- **Dependencies:** Declare all `require()` dependencies
- **Security:** Mark if uses FFI, network, or file access
- **Testing:** Test your script before submitting

### Review Process

1. Automated checks run on your PR
2. Maintainer reviews code for security
3. If approved, PR is merged
4. Package auto-deploys to CDN within 2 minutes

---

## Русский

### Как добавить пакет

#### Используя веб-интерфейс (Рекомендуется)

1. **Откройте веб-интерфейс:**
   Перейдите на [URL сайта реестра]

2. **Войдите через GitHub:**
   Нажмите "Войти через GitHub" и авторизуйте приложение

3. **Загрузите ваш пакет:**
   - Перейдите на страницу "Загрузить"
   - Перетащите ваш .lua файл(ы)
   - Система автоматически:
     - Извлечёт метаданные из `script_name()`, `script_version()`, `script_author()`
     - Проанализирует и найдёт все зависимости
     - Обнаружит функции безопасности (FFI, сетевой доступ, файловый I/O)
   - Проверьте обнаруженные метаданные и отредактируйте при необходимости
   - Добавьте теги (через запятую)
   - Добавьте ссылку на источник (опционально)

4. **Создайте Pull Request:**
   - Нажмите "Загрузить и создать PR"
   - Система автоматически создаст форк, закоммитит файлы и откроет PR
   - Дождитесь проверки мейнтейнера

#### Используя CLI

1. **Скачайте CLI:**
   Скачайте готовый бинарник из [Releases](https://github.com/Deps-Tech/deps-registry/releases/latest):
   - Windows: `tools-cli-windows-amd64.exe`
   - Linux: `tools-cli-linux-amd64` или `tools-cli-linux-arm64`
   - macOS: `tools-cli-darwin-amd64` или `tools-cli-darwin-arm64`

2. **Клонируйте репозиторий:**
```bash
git clone https://github.com/Deps-Tech/deps-registry.git
cd deps-registry
```

3. **Добавьте свой скрипт:**
```bash
./tools-cli add script \
  --source ~/путь/к/ВашСкрипт.lua \
  --tags "helper,rp,automation"
```

CLI автоматически:
- Извлечёт метаданные из `script_name()`, `script_version()`, `script_author()`
- Проанализирует и найдёт все зависимости
- Обнаружит функции безопасности (FFI, сетевой доступ, файловый I/O)
- Сгенерирует манифест с SHA256 хешами
- Создаст правильную структуру директорий

4. **Проверьте сгенерированные файлы:**
```bash
git diff
```

5. **Закоммитьте и запушьте:**
```bash
git add scripts/ваш-скрипт/
git commit -m "feat(scripts): add ваш-скрипт v1.0"
git push origin main
```

6. **Создайте Pull Request** на GitHub

#### Ручной метод

1. Создайте директорию: `scripts/ваш-скрипт/1.0/`
2. Скопируйте ваш Lua файл(ы)
3. Создайте манифест `dep.json`
4. Запустите валидацию: `./tools-cli validate`

### Рекомендации по пакетам

- **Именование:** Используйте lowercase с дефисами (например, `fake-documents`)
- **Версионирование:** Следуйте семантическому версионированию (например, `1.0.0`)
- **Зависимости:** Объявляйте все `require()` зависимости
- **Безопасность:** Отмечайте использование FFI, сети или файлового доступа
- **Тестирование:** Протестируйте скрипт перед отправкой

### Процесс ревью

1. Автоматические проверки запускаются на вашем PR
2. Мейнтейнер проверяет код на безопасность
3. При одобрении PR сливается
4. Пакет автоматически деплоится на CDN в течение 2 минут
