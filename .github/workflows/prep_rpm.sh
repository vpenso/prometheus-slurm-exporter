mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS,tmp}
echo '%_topdir %(echo $HOME)/rpmbuild' > ~/.rpmmacros
cp README.md ~/rpmbuild/SOURCES
cp LICENSE ~/rpmbuild/SOURCES
cp bin/prometheus-slurm-exporter ~/rpmbuild/SOURCES
cp lib/systemd/prometheus-slurm-exporter.service ~/rpmbuild/SOURCES
cp packages/rpm/*.spec ~/rpmbuild/SPECS