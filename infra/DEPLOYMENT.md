# Deployment Guide: laundry-management-system on k3s via GitLab CI

## 1. Создать GitLab Deploy Token

GitLab -> Settings -> Repository -> **Deploy tokens**

- Name: `k8s-registry`
- Scope: `read_registry`
- Нажать **Create deploy token**, сохранить `username` и `token`

Они понадобятся на следующем шаге.

---

## 2. Подготовка Ansible

### Установить зависимости

```bash
ansible-galaxy collection install -r requirements.yml
```

### Заполнить инвентарь

`inventory/hosts.yml`:

```yaml
all:
  children:
    webservers:
      hosts:
        cms_server:
          ansible_host: <IP_СЕРВЕРА>
          ansible_user: ubuntu
          ansible_ssh_private_key_file: ~/.ssh/id_rsa
```

### Заполнить переменные и зашифровать секреты

Вписать в `group_vars/all.yml` токен из шага 1 и путь к образу:

```yaml
gitlab_registry_username: "your_deploy_token_username"
gitlab_registry_token: "your_deploy_token_value"
image_repo: registry.gitlab.com/<your-gitlab-group>/laundry-management-system
```

Зашифровать перед коммитом:

```bash
ansible-vault encrypt group_vars/all.yml
```

Для редактирования:

```bash
ansible-vault edit group_vars/all.yml
```

### Запустить playbook

```bash
ansible-playbook site.yml -i inventory/hosts.yml --ask-vault-pass --ask-become-pass
```

---

## 3. Что делает playbook

Роль `k3s_db_role` выполняет:

1. Устанавливает k3s (одиночный node)
2. Открывает порты через UFW (22, 80, 443, 6443)
3. Создаёт namespace'ы `staging` и `production`
4. Добавляет секрет `gitlab-registry` для pull образов в оба namespace'а
5. Рендерит k8s-манифесты из шаблона и применяет их через `kubectl apply`

---

## 4. Настройка GitLab CI/CD

После того как playbook отработал и k3s установлен, снять kubeconfig с сервера.

### Получить KUBE_CONFIG

На сервере (заменить IP на реальный):

```bash
sudo cat /etc/rancher/k3s/k3s.yaml \
  | sed 's/127.0.0.1/<IP_СЕРВЕРА>/g' \
  | base64 -w0
```

### Добавить переменную в GitLab

GitLab -> Settings -> CI/CD -> Variables -> **Add variable**:

| Key           | Value              | Type     |
|---------------|--------------------|----------|
| `KUBE_CONFIG` | вывод команды выше | Variable |

> Убедись что значение - одна строка без переносов строки.

---

## 5. K8s манифесты

Манифесты генерируются из Jinja2-шаблона `roles/k3s_db_role/templates/k8s_deployment.yaml.j2`.

Параметры задаются в `group_vars/all.yml` или `roles/k3s_db_role/defaults/main.yml`:

| Переменная                | По умолчанию              | Описание                        |
|---------------------------|---------------------------|---------------------------------|
| `app_name`                | `laundry-backend`         | Имя приложения в k8s            |
| `image_repo`              | `registry.gitlab.com/...` | Путь к образу в GitLab Registry |
| `image_tag`               | `latest`                  | Тег образа                      |
| `container_port`          | `8080`                    | Порт контейнера                 |
| `staging_replicas`        | `1`                       | Число реплик в staging          |
| `staging_ingress_host`    | `""`                      | Хост Ingress (пусто = все)      |
| `production_replicas`     | `2`                       | Число реплик в production       |
| `production_ingress_host` | `laundry.example.com`     | Хост Ingress для production     |

### Секрет с переменными окружения приложения

Создаётся вручную на сервере после первого деплоя:

```bash
sudo kubectl create secret generic laundry-backend-env \
  --from-literal=POSTGRES_PASSWORD=secret \
  --namespace=staging

sudo kubectl create secret generic laundry-backend-env \
  --from-literal=POSTGRES_PASSWORD=secret \
  --namespace=production
```

---

## 6. Проверка деплоя

После пуша в `main`:

1. **test** - `go test ./...`
2. **build** - сборка Docker-образа, пуш в GitLab Registry
3. **deploy:staging** - автоматически обновляет deployment в `staging`
4. **deploy:production** - запускается вручную через GitLab UI

Проверить статус подов:

```bash
sudo kubectl get pods -n staging
sudo kubectl get pods -n production
```

Приложение доступно по `http://<IP_СЕРВЕРА>` (через Traefik на порту 80).
