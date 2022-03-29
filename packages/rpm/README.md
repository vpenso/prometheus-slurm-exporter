## Setup a development environment

In order to build your own RPMs you have to setup your own _development environment_ as explained [here](https://wiki.centos.org/HowTos/SetupRpmBuildEnvironment).

### Quick setup

1. Install the ``rpm-build`` and ``rpmdevtools`` packages.

2. Login as **normal user** and create the following structure in the home directory:
```bash
mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
```
3. Create the RPM macros file with the following content:
```bash
echo "%_topdir $HOME/rpmbuild" > ~/.rpmmacros
```

### Prepare the build
 
4. Clone this repo:
```bash
git clone https://github.com/vpenso/prometheus-slurm-exporter.git
```
5. Get into the source directory and copy the following files under ``~/rpmbuild/SOURCES``:
```bash
cd prometheus-slurm-exporter
cp README.md ~/rpmbuild/SOURCES
cp LICENSE ~/rpmbuild/SOURCES
cp lib/systemd/prometheus-slurm-exporter.service ~/rpmbuild/SOURCES
```
6. Copy the SPEC file in the proper directory:
```bash
cd prometheus-slurm-exporter
cp packages/rpm/*.spec ~/rpmbuild/SPECS
```

### Build the RPM package

8. Build the RPM based on your SPEC file:
```bash
cd $HOME/rpmbuild/SPECS
rpmbuild -ba prometheus-slurm-exporter.spec
```
9. The RPM package will be placed under $HOME/rpmbuild/RPMS/x86_64
