## RPM packaging

This directory holds the RPM spec used to build `prometheus-slurm-exporter`
packages for RHEL 9 / Rocky Linux 9 and compatible distributions.

## Automated build (maintainers)

The **Release Prometheus Slurm Exporter RPM** workflow
(`.github/workflows/build-and-release-rpm.yaml`) builds an RPM inside a
Rocky Linux 9 container and attaches it to an existing GitHub release. It is
triggered manually from the Actions tab ("Run workflow") with a `release-tag`
input naming an already-published tag (e.g. `0.23`). The workflow checks out
that tag, so the tag must exist before it is run.

## Building an RPM manually

1. Install the build tools:

   ```bash
   dnf install -y rpm-build rpmdevtools make git go jq systemd
   ```

2. Create the rpmbuild tree:

   ```bash
   mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
   echo '%_topdir %(echo $HOME)/rpmbuild' > ~/.rpmmacros
   ```

3. From a checkout of this repo, stage the sources the spec expects:

   ```bash
   cp README.md LICENSE ~/rpmbuild/SOURCES
   cp lib/systemd/prometheus-slurm-exporter.service ~/rpmbuild/SOURCES
   cp packages/rpm-ci/*.spec ~/rpmbuild/SPECS
   ```

4. Build, substituting a published tag for `VERSION`:

   ```bash
   cd ~/rpmbuild
   spectool -g -R --define '_version VERSION' --define '_release 1' \
       SPECS/prometheus-slurm-exporter.spec
   rpmbuild -bb --define '_version VERSION' --define '_release 1' \
       SPECS/prometheus-slurm-exporter.spec
   ```

The resulting RPM is placed under `~/rpmbuild/RPMS/<arch>/`. Note that the
packaged unit runs the exporter as the dedicated `slurm_exporter` system user
(created automatically at install time by the spec) and binds it to
`0.0.0.0:9341` rather than the program default of `:8080`.
