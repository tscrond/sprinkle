
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

#### `init`

Initializes a configuration file for managing Proxmox resources.

```bash
sprinkle init
```

This command creates a basic configuration file, typically named `sprinkle.yaml`, in your current working directory. You can modify this file to define your VMs and LXCs.

#### `apply`

Applies the changes defined in the configuration file to your Proxmox setup.

```bash
sprinkle apply
```

This command reads the `sprinkle.yaml` file and ensures that the Proxmox infrastructure is in sync with the configurations specified in the file. It will create, update, or delete VMs and LXCs as necessary.

#### `destroy`

Destroys the resources (VMs or LXCs) defined in your configuration file.

```bash
sprinkle destroy
```

This command deletes the VMs or LXCs that were previously created using `sprinkle apply`. 

### Configuration File (`sprinkle.yaml`)

The configuration file is a YAML file where you define your infrastructure. Here’s an example:

```yaml
vms:
  - id: 100
    name: "web-server"
    cores: 2
    memory: 4096
    disk: 20
    net:
      bridge: "vmbr0"
      ip: "dhcp"
      firewall: true
    storage: "local-lvm"
    template: "debian-11"

lxcs:
  - id: 200
    name: "app-container"
    cores: 1
    memory: 1024
    disk: 10
    net:
      bridge: "vmbr0"
      ip: "192.168.1.100"
      firewall: false
    storage: "local-lvm"
    template: "ubuntu-20.04"
```

In this example, the configuration defines both a VM (`web-server`) and an LXC (`app-container`) with specific settings, including CPU cores, memory, disk size, and network configurations.

### Flags

- `-f`, `--file`: Specify the path to the configuration file. If not specified, `sprinkle` looks for a file named `sprinkle.yaml` in the current directory.
- `-d`, `--dry-run`: Simulate the changes without actually applying them. This is useful for previewing the changes before making them.
- `-v`, `--verbose`: Enable verbose output for more detailed logs.

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

- **VM/LXC Network Configuration:** Define network bridge, IP address type (`dhcp` or static), and firewall settings.
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
