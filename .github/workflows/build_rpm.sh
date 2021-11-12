cd $HOME/rpmbuild || exit 1
spectool -g -R SPECS/prometheus-slurm-exporter.spec
rpmbuild -ba SPECS/prometheus-slurm-exporter.spec