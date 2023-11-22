## GitHub Actions rpmbuild environment

Generally this directory is used by the Release Prometheus Slurm Exporter RPM [workflow](https://github.com/stackhpc/prometheus-slurm-exporter/actions/workflows/build-and-release-rpm.yaml) to
build and release prometheus-slurm-exporter RPMs.

To trigger a build of an RPM for a new version of prometheus-slurm-exporter, run the following commands to sync the upstream repository:

```bash
git clone git@github.com:stackhpc/prometheus-slurm-exporter.git
cd prometheus-slurm-exporter
git remote add upstream https://github.com/vpenso/prometheus-slurm-exporter
git fetch upstream
git fetch upstream --tags
git checkout master
git merge upstream/master
git push origin master
git push origin master --tags
```

Then, the Release Prometheus Slurm Exporter RPM [workflow](https://github.com/stackhpc/prometheus-slurm-exporter/actions/workflows/build-and-release-rpm.yaml) can then be triggered with
with a newly-synced tag.


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
cp packages/rpm-ci/*.spec ~/rpmbuild/SPECS
```

### Build the RPM package

8. Build the RPM based on your SPEC file:
```bash
cd $HOME/rpmbuild/SPECS
spectool -g -R --define '_version {prometheus_slurm_exporter_release_version}' prometheus-slurm-exporter.spec
rpmbuild -ba --define '_version {prometheus_slurm_exporter_release_version}' prometheus-slurm-exporter.spec
```
9. The RPM package will be placed under $HOME/rpmbuild/RPMS/x86_64
