# -*- mode: ruby -*-
# vi: set ft=ruby :


# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://atlas.hashicorp.com/search.
  config.vm.box = "ubuntu/xenial64"

  config.ssh.forward_agent = true
  config.ssh.forward_x11 = true

  # Staging / Development / Testing port
  config.vm.network "forwarded_port", guest: 3000, host: 3000
  config.vm.network "forwarded_port", guest: 8005, host: 8005
  config.vm.network "forwarded_port", guest: 8007, host: 8007
  config.vm.network "forwarded_port", guest: 7777, host: 7777
  config.vm.network "forwarded_port", guest: 2345, host: 2345
  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  # config.vm.network "private_network", ip: "192.168.33.10"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  config.vm.synced_folder "./", "/home/vagrant/go/src/github.com/phzfi/RIC/"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  config.vm.provider "virtualbox" do |vb|
      # Headless box named RIC
      vb.gui = false
        vb.name = "RIC"

      # 2 CPU:s with 2GB of RAM
	     vb.memory = 2048
	     vb.cpus = 2
  end

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.

  config.vm.provision "shell", path: "scripts/provision/vagrant_setup.sh"
  config.vm.provision "shell", path: "scripts/provision/vagrant_setup_imagemagick.sh"
  config.vm.provision "shell", path: "scripts/provision/vagrant_setup_go.sh"
end

