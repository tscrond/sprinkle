
# Sprinkle - A Proxmox VM/LXC Management Tool

**Sprinkle** is a command-line utility designed to simplify the management of virtual machines (VMs) and Linux Containers (LXCs) on a Proxmox server. It offers a declarative, Terraform-like experience, allowing users to define infrastructure in configuration files and apply changes seamlessly.

## Features

- **Declarative Management:** Define VMs and LXCs in configuration files, which serve as a blueprint for your infrastructure.
- **Terraform-Like Workflow:** Similar to Terraform, you define the desired state and apply it, and Sprinkle ensures that your Proxmox setup matches that state.
- **Simplified Proxmox Interaction:** Communicate with Proxmox directly from the command line without needing to navigate through the web interface.
- **Support for VMs and LXCs:** Manage both virtual machines and containers.
- **Configuration Flexibility:** Easily manage VM and LXC configurations, including storage, network settings, and more.

## Installation

To install `sprinkle`, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/tscrond/sprinkle.git
   ```
2. Navigate into the project directory:
   ```bash
   cd sprinkle
   ```
3. Build the binary (assuming you have Go installed):
   ```bash
   go build -o sprinkle
   ```
4. Place the binary in a directory accessible in your PATH for easy access:
   ```bash
   sudo mv sprinkle /usr/local/bin/
   ```

Alternatively, you can download pre-compiled binaries from the [Releases](https://github.com/tscrond/sprinkle/releases) section of the repository.

## Usage

The primary command for Sprinkle is `sprinkle`, which interacts with configuration files and applies the desired infrastructure changes.

### Command Syntax

```bash
sprinkle <command> [flags]
```

### Available Commands

#### `init` - TODO

Initializes a configuration file for managing Proxmox resources.

```bash
sprinkle init
```

This command creates a basic configuration file, typically named `sprinkle-config.yaml`, in your current working directory. You can modify this file to define your VMs and LXCs.

#### `apply`

Applies the changes defined in the configuration file to your Proxmox setup.

```bash
sprinkle apply
```

This command reads the `sprinkle-config.yaml` file and ensures that the Proxmox infrastructure is in sync with the configurations specified in the file. It will create, update, or delete VMs and LXCs as necessary.

#### `destroy` - TODO

Destroys the resources (VMs or LXCs) defined in your configuration file.

```bash
sprinkle destroy
```

This command deletes the VMs or LXCs that were previously created using `sprinkle apply`. 

### Configuration File (`sprinkle-config.yaml`)

The configuration file is a YAML file where you define your infrastructure. Here’s an example:

```yaml
hosts:
  genesis:
    api-url: "192.168.1.102:8006"
    target-node: "genesis"
    lxc:
      default:
        start-on-boot: false
        storage-backend: local-lvm
        template-backend: local
        default-gateway: "192.168.1.1"
        network-bridge: "vmbr0"
        network-interface: "eth0"
      machines:
        - name: "firstlxc123"
          vmid: 112
          os-template: "debian-11-standard_11.7-1_amd64.tar.zst"
          ssh-public-keys:
            - key: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIC48nEoa2rRazXTxZ4anL+6CL2bGXTo6w6XcDpmcd3pE tomasz.skrond@boar.network"
            - path: "/Users/tskr/.ssh/genesis.pub"
          ip-address: "192.168.1.30/24"
          cpus: 4
          memory: 2048
          disk-size: 100
          swap-size: 20
          tags: "asdf;fdsa;fds"
    vm:
      default:
        start-on-boot: false
        storage-backend: local-lvm
        template-backend: local
        default-gateway: "192.168.1.1"
        network-bridge: "vmbr0"
        network-interface: "eth0"
      machines:
        - name: "firstvm"
          vmid: 114
          iso: "ubuntu-24.04.1-live-server-amd64.iso"
          ssh-public-keys:
            - key: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5aaaafIC48nEoa2rRazXTxZ4anL+6CL2bGXTo6w6XcDpmcd3pE key@tomo"
            - path: "/Users/user/.ssh/genesis.pub"
          cpus: 4
          memory: 2048
          disk-size: 100
          swap-size: 0
          tags: "asdf;fdsa;oooooooo"


  hyperbook:
    api-url: "192.168.1.100:8006"
    target-node: "hyperbook"
    lxc:
      default:
        start-on-boot: false
        storage-backend: local-lvm
        template-backend: local
        default-gateway: "192.168.1.1"
        network-bridge: "vmbr0"
        network-interface: "eth0"
      machines:
      vm:
        default:
          start-on-boot: false
          storage-backend: local-lvm
          template-backend: local
          default-gateway: "192.168.1.1"
          network-bridge: "vmbr0"
          network-interface: "eth0"
        machines:
          - name: "firstvm"
            vmid: 11443
            iso: "ubuntu-24.04.1-live-server-amd64.iso"
            ssh-public-keys:
              - key: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5aaaafIC48nEoa2rRazXTxZ4anL+6CL2bGXTo6w6XcDpmcd3pE key@tomo"
              - path: "/Users/user/.ssh/genesis.pub"
            cpus: 4
            memory: 2048
            disk-size: 100
            swap-size: 0
            tags: "ffasdfqwer1111234;fdsa;fds"
```

In this example, the configuration defines both a VM (`web-server`) and an LXC (`app-container`) with specific settings, including CPU cores, memory, disk size, and network configurations.

### Flags

- `-f`, `--file`: Specify the path to the configuration file. If not specified, `sprinkle` looks for a file named `sprinkle-config.yaml` in the current directory.
- `-d`, `--dry-run` - TODO: Simulate the changes without actually applying them. This is useful for previewing the changes before making them.
- `-v`, `--verbose` - TODO: Enable verbose output for more detailed logs.

### Example Usage

1. **Initialize a new configuration:**

   ```bash
   sprinkle init
   ```

2. **Apply the configuration to create or update resources:**

   ```bash
   sprinkle apply
   ```

3. **Destroy the resources defined in the configuration:**

   ```bash
   sprinkle destroy
   ```

4. **Preview changes without applying them (dry run):**

   ```bash
   sprinkle apply --dry-run
   ```

## Configuration Options

- **VM/LXC Network Configuration:** Define network bridge, IP addresses and SSH keys.
- **Resources:** Specify CPU cores, memory, disk size, and storage backend.
- **Templates:** Choose a template for VMs and LXCs, which should be pre-existing on your Proxmox server.

## Requirements

- **Proxmox Server**: Ensure your Proxmox server is set up and accessible.
- **Proxmox API**: Sprinkle interacts with the Proxmox API, so make sure the API is enabled and accessible on your Proxmox instance.
- **Go**: For building from source, you’ll need Go installed.

## Contributing

Contributions are welcome! If you find bugs or have suggestions for improvement, please open an issue or create a pull request. 

- Fork the repository
- Create a feature branch (`git checkout -b feature-branch`)
- Commit your changes (`git commit -am 'Add new feature'`)
- Push to the branch (`git push origin feature-branch`)
- Open a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
