# Contributing Guide / Руководство по внесению вклада

## English

### How to Add a Package

#### Using the CLI (Recommended)

1. **Clone the repository:**
```bash
git clone https://github.com/Deps-Tech/deps-registry.git
cd deps-registry/tools
```

2. **Build the CLI:**
```bash
go build -o tools-cli ./cmd/tools-cli
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
cd ..
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

#### Используя CLI (Рекомендуется)

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/Deps-Tech/deps-registry.git
cd deps-registry/tools
```

2. **Соберите CLI:**
```bash
go build -o tools-cli ./cmd/tools-cli
```

3. **Добавьте свой скрипт:**
```bash
./tools-cli add script \
  --source ~/путь/к/ВашСкрипт.lua \
  --tags "helper,rp,automation"
```

CLI автоматически:
- Извлечет метаданные из `script_name()`, `script_version()`, `script_author()`
- Проанализирует и найдет все зависимости
- Обнаружит функции безопасности (FFI, сетевой доступ, файловый I/O)
- Сгенерирует манифест с SHA256 хешами
- Создаст правильную структуру директорий

4. **Проверьте сгенерированные файлы:**
```bash
cd ..
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

