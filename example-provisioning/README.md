# Provisioning Example

This directory contains a complete inventory and playbook for `minipaas-node` and `minipaas-swarm`.

## Requirements

- Debian or Ubuntu hosts with SSH access
- Python 3 and Ansible on your workstation
- A key-based SSH user with sudo access

Install the collection and role dependencies:

```bash
ansible-galaxy install -r ../minipaas-swarm/requirements.yml
ansible-galaxy role install geerlingguy.docker
```

## Configure

Edit `inventory.ini` with your hosts.

The manager group opens the Swarm management ports and HTTP/HTTPS. The worker group opens the gossip and overlay ports. Add any published service ports to `node_firewall_tcp_ports` in the matching group file.

Telegram alerts are optional. Set both values when you want monitoring and log alerts:

```bash
export TELEGRAM_TOKEN="123:ABC"
export TELEGRAM_CHAT="999111222"
```

## Run

```bash
cd example-provisioning
ansible-playbook playbook.yml --check --tags node
ansible-playbook playbook.yml
```

`minipaas-node` hardens SSH by default. Run the node layer with `--check` first and confirm that key-based access works before the first real run. To disable SSH hardening, set `node_ssh_harden: false` in `group_vars/managers.yml`, `group_vars/workers.yml`, or `inventory.ini`.

Update OS packages only when you ask for it:

```bash
ansible-playbook update.yml
```

## Selective runs

```bash
ansible-playbook playbook.yml --tags node
ansible-playbook playbook.yml --tags swarm
ansible-playbook playbook.yml --tags docker
```
