# Deployment

## Быстрый старт

### Зашифровать секреты в файле `group_vars/all.yml`

```bash
ansible-vault encrypt group_vars/all.yml
```

#### Для расшифровки

```bash
ansible-vault decrypt group_vars/all.yml
```

### Установить Ansible-коллекции

```bash
ansible-galaxy collection install -r requirements.yml
```

### Запустить playbook

```bash
ansible-playbook site.yml -i inventory/hosts.yml --ask-vault-pass --ask-become-pass
```
